package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	httphandler "ride-sharing/services/trip-service/internal/http"
	"ride-sharing/services/trip-service/internal/infrastructure/events"
	"ride-sharing/services/trip-service/internal/infrastructure/grpc"
	"ride-sharing/services/trip-service/internal/infrastructure/repository"
	"ride-sharing/services/trip-service/internal/service"
	"ride-sharing/shared/db"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"
	"strings"
	"syscall"

	grpcserver "google.golang.org/grpc"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var GrpcAddr = env.GetString("GRPC_ADDR", ":9093")

func main() {
	// Initialize Tracing
	tracerCfg := tracing.Config{
		ServiceName:    "trip-service",
		Environment:    env.GetString("ENVIRONMENT", "development"),
		JaegerEndpoint: env.GetString("JAEGER_ENDPOINT", "http://jaeger:14268/api/traces"),
	}

	sh, err := tracing.InitTracer(tracerCfg)
	if err != nil {
		log.Fatalf("Failed to initialize the tracer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer sh(ctx)

	// Initialize MongoDB
	mongoClient, err := db.NewMongoClient(ctx, db.NewMongoDefaultConfig())
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB, err: %v", err)
	}
	defer mongoClient.Disconnect(ctx)

	mongoDb := db.GetDatabase(mongoClient, db.NewMongoDefaultConfig())

	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")

	mongoDBRepo := repository.NewMongoRepository(mongoDb)
	svc := service.NewService(mongoDBRepo)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	// RabbitMQ connection
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()

	log.Println("Starting RabbitMQ connection")

	publisher := events.NewTripEventPublisher(rabbitmq)

	// Start driver consumer
	driverConsumer := events.NewDriverConsumer(rabbitmq, svc)
	go driverConsumer.Listen()

	// Initialize the gRPC server
	grpcServer := grpcserver.NewServer(tracing.WithTracingInterceptors()...)
	grpc.NewGRPCHandler(grpcServer, svc, publisher)

	// Start payment consumer
	paymentConsumer := events.NewPaymentConsumer(rabbitmq, svc)
	go paymentConsumer.Listen()

	log.Printf("Starting gRPC server Trip service on port %s", GrpcAddr)

	// Combine gRPC and HTTP Health Check on the same port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Trip Service is Healthy"))
	})

	// Add REST endpoint for trip preview
	httpHandler := httphandler.NewHandler(svc)
	mux.HandleFunc("/api/preview", httpHandler.HandlePreview)

	h2Handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		log.Printf("[MULTIPLEXER] Request: %s %s | Proto: HTTP/%d.%d | Content-Type: %s | User-Agent: %s",
			r.Method, r.URL.Path, r.ProtoMajor, r.ProtoMinor, contentType, r.Header.Get("User-Agent"))

		// Check if this is a gRPC request by Content-Type header
		// Note: We don't strictly require HTTP/2 because Render's load balancer may downgrade
		if strings.HasPrefix(contentType, "application/grpc") {
			log.Printf("[MULTIPLEXER] Routing to gRPC handler")
			grpcServer.ServeHTTP(w, r)
		} else {
			log.Printf("[MULTIPLEXER] Routing to HTTP handler (health check)")
			mux.ServeHTTP(w, r)
		}
	}), &http2.Server{})

	server := &http.Server{
		Addr:    ":" + port,
		Handler: h2Handler,
	}

	go func() {
		log.Printf("Starting Multiplexed Server (gRPC + HTTP) on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to serve: %v", err)
			cancel()
		}
	}()

	// wait for the shutdown signal
	<-ctx.Done()
	log.Println("Shutting down the server...")
	grpcServer.GracefulStop()
}
