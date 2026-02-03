package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"ride-sharing/services/trip-service/internal/domain"
	tripTypes "ride-sharing/services/trip-service/pkg/types"
	"ride-sharing/shared/types"
)

type TripService interface {
	GetRoute(ctx context.Context, pickup, destination *types.Coordinate, useOSRMApi bool) (*tripTypes.OsrmApiResponse, error)
	EstimatePackagesPriceWithRoute(route *tripTypes.OsrmApiResponse) []*domain.RideFareModel
	GenerateTripFares(ctx context.Context, rideFares []*domain.RideFareModel, userID string, route *tripTypes.OsrmApiResponse) ([]*domain.RideFareModel, error)
}

type Handler struct {
	service TripService
}

func NewHandler(service TripService) *Handler {
	return &Handler{service: service}
}

type PreviewRequest struct {
	UserID string `json:"user_id"`
	Pickup struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"pickup"`
	Destination struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"destination"`
	RideType string `json:"ride_type"`
}

type PreviewResponse struct {
	Route     *tripTypes.OsrmApiResponse `json:"route"`
	RideFares []*domain.RideFareModel    `json:"ride_fares"`
}

func (h *Handler) HandlePreview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	pickup := &types.Coordinate{
		Latitude:  req.Pickup.Latitude,
		Longitude: req.Pickup.Longitude,
	}

	destination := &types.Coordinate{
		Latitude:  req.Destination.Latitude,
		Longitude: req.Destination.Longitude,
	}

	ctx := r.Context()

	// Get route (using mock for now, change to true for real OSRM)
	route, err := h.service.GetRoute(ctx, pickup, destination, false)
	if err != nil {
		log.Printf("Error getting route: %v", err)
		http.Error(w, "Failed to get route", http.StatusInternalServerError)
		return
	}

	// Estimate fares
	estimatedFares := h.service.EstimatePackagesPriceWithRoute(route)

	// Generate and save fares
	fares, err := h.service.GenerateTripFares(ctx, estimatedFares, req.UserID, route)
	if err != nil {
		log.Printf("Error generating trip fares: %v", err)
		http.Error(w, "Failed to generate trip fares", http.StatusInternalServerError)
		return
	}

	response := PreviewResponse{
		Route:     route,
		RideFares: fares,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
