package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ride-sharing/services/api-gateway/grpc_clients"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/env"
	"ride-sharing/shared/messaging"
	"ride-sharing/shared/tracing"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
)

var tracer = tracing.GetTracer("api-gateway")

func handleTripStart(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripStart")
	defer span.End()

	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqBody startTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to parse JSON data")
		return
	}

	defer r.Body.Close()

	// Why we need to create a new client for each connection:
	// because if a service is down, we don't want to block the whole application
	// so we create a new client for each connection
	tripService, err := grpc_clients.NewTripServiceClient()
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to close the client to avoid resource leaks!
	defer tripService.Close()

	trip, err := tripService.Client.CreateTrip(ctx, reqBody.toProto())
	if err != nil {
		log.Printf("DEBUG: gRPC CreateTrip failed: %v", err)
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to start trip: %v", err))
		return
	}

	response := contracts.APIResponse{Data: trip}

	writeJSON(w, http.StatusCreated, response)
}

func handleTripPreview(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleTripPreview")
	defer span.End()

	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var reqBody previewTripRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		writeJSONError(w, http.StatusBadRequest, "failed to parse JSON data")
		return
	}

	defer r.Body.Close()

	// validation
	if reqBody.UserID == "" {
		writeJSONError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	// Call the trip service REST API
	tripServiceURL := env.GetString("TRIP_SERVICE_URL", "http://localhost:8080")
	apiURL := tripServiceURL + "/api/preview"

	// Prepare the request payload
	payload := map[string]interface{}{
		"user_id": reqBody.UserID,
		"pickup": map[string]float64{
			"latitude":  reqBody.Pickup.Latitude,
			"longitude": reqBody.Pickup.Longitude,
		},
		"destination": map[string]float64{
			"latitude":  reqBody.Destination.Latitude,
			"longitude": reqBody.Destination.Longitude,
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR: Failed to marshal payload: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to prepare request")
		return
	}

	// Make the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Printf("ERROR: Failed to create request: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to call trip service: %v", err)
		writeJSONError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to preview trip: %v", err))
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("ERROR: Trip service returned status %d: %s", resp.StatusCode, string(bodyBytes))
		writeJSONError(w, http.StatusInternalServerError, "Trip service error")
		return
	}

	// Read and forward the response
	var tripPreview interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tripPreview); err != nil {
		log.Printf("ERROR: Failed to decode trip service response: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "Failed to decode response")
		return
	}

	response := contracts.APIResponse{Data: tripPreview}
	writeJSON(w, http.StatusCreated, response)
}

func handleStripeWebhook(w http.ResponseWriter, r *http.Request, rb *messaging.RabbitMQ) {
	ctx, span := tracer.Start(r.Context(), "handleStripeWebhook")
	defer span.End()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	webhookKey := env.GetString("STRIPE_WEBHOOK_KEY", "")
	if webhookKey == "" {
		log.Printf("Webhook key is required")
		return
	}

	event, err := webhook.ConstructEventWithOptions(
		body,
		r.Header.Get("Stripe-Signature"),
		webhookKey,
		webhook.ConstructEventOptions{
			IgnoreAPIVersionMismatch: true,
		},
	)
	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	log.Printf("Received Stripe event: %v", event)

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			http.Error(w, "Invalid payload", http.StatusBadRequest)
			return
		}

		payload := messaging.PaymentStatusUpdateData{
			TripID:   session.Metadata["trip_id"],
			UserID:   session.Metadata["user_id"],
			DriverID: session.Metadata["driver_id"],
		}

		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Error marshalling payload: %v", err)
			http.Error(w, "Failed to marshal payload", http.StatusInternalServerError)
			return
		}

		message := contracts.AmqpMessage{
			OwnerID: session.Metadata["user_id"],
			Data:    payloadBytes,
		}

		if err := rb.PublishMessage(
			ctx,
			contracts.PaymentEventSuccess,
			message,
		); err != nil {
			log.Printf("Error publishing payment event: %v", err)
			http.Error(w, "Failed to publish payment event", http.StatusInternalServerError)
			return
		}
	}
}

func writeJSONError(w http.ResponseWriter, code int, message string) {
	response := contracts.APIResponse{
		Error: &contracts.APIError{
			Message: message,
			Code:    fmt.Sprintf("%d", code),
		},
	}
	writeJSON(w, code, response)
}
