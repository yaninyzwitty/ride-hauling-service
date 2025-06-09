package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/yaninyzwitty/ride-hauling-app/services/api-gateway/grpc_clients/trip_client"
	"github.com/yaninyzwitty/ride-hauling-app/shared/proto/trip"
	"github.com/yaninyzwitty/ride-hauling-app/shared/types"
)

type Service struct {
	URL  string `json:"url"`
	Name string `json:"name"`
}

// test k8
func HandleGetServices(w http.ResponseWriter, r *http.Request) {
	// Logging the pod name to visually see the replica scaling load balancing
	podName := os.Getenv("HOSTNAME") // k8s sets this automatically
	slog.Info("Handling request",
		"pod", podName,
		"from", r.RemoteAddr,
		"path", r.URL.Path,
	)
	var services []Service
	services = append(services, Service{
		URL:  "http://localhost:8082",
		Name: "driver-service",
	})

	err := writeJSON(w, http.StatusOK, services)
	if err != nil {
		return
	}
}

func handleCreateTrip(w http.ResponseWriter, r *http.Request, tripClient *trip_client.TripServiceClient) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var ctx = r.Context()

	var reqBody struct {
		Pickup      types.Location `json:"pickup"`
		Destination types.Location `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	resp, err := tripClient.CreateTrip(ctx, &trip.CreateTripRequest{
		StartLocation: &trip.Coordinate{
			Latitude:  float32(reqBody.Pickup.Latitude),
			Longitude: float32(reqBody.Pickup.Longitude),
		},
		EndLocation: &trip.Coordinate{
			Latitude:  float32(reqBody.Destination.Latitude),
			Longitude: float32(reqBody.Destination.Longitude),
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, resp)

}
