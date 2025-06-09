package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	driverClientPb "github.com/yaninyzwitty/ride-hauling-app/services/api-gateway/grpc_clients/driver_client"
	driverPb "github.com/yaninyzwitty/ride-hauling-app/shared/proto/driver"
)

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request, client *driverClientPb.DriverServiceClient) {
	var ctx = r.Context()

	// upgrade http connection to websocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade http connection to websocket connection", "error", err)
		return
	}
	defer func(conn *websocket.Conn) {
		_ = conn.Close()
	}(conn)

	slog.Info("websocket connection established")
	// make a channel to listen for client updates
	clientLocations := make(chan *driverPb.Driver)

	go func() {
		for {
			// here we read message from the websocket
			_, msg, err := conn.ReadMessage()
			if err != nil {
				slog.Error("error reading websocket msg:", "error", err)
				close(clientLocations)
				return
			}

			var location driverPb.Driver
			if err := json.Unmarshal(msg, &location); err != nil {
				slog.Info("Invalid location data: ", "error", err)
				continue
			}

			clientLocations <- &location
		}

	}()

	// here we start streaming driver updates
	stream, err := client.FindNearbyDrivers(ctx)

	if err != nil {
		slog.Error("failed to start grpc stream: ", "error", err)
		os.Exit(1)
	}
	defer stream.CloseSend()

	// Handle sending client location updates to the driver service

	go func() {
		for location := range clientLocations {
			if err := stream.Send(&driverPb.FindNearbyDriversRequest{
				Location: &driverPb.Location{
					Latitude:  location.Location.Latitude,
					Longitude: location.Location.Longitude,
				},
			}); err != nil {
				slog.Error("failed to send driver location update: ", "error", err)
				return
			}
		}
	}()

	// here we handle receiving driver updates
	for {
		response, err := stream.Recv()
		if err != nil {
			slog.Error("error receiving driver updates: ", "error", err)
			os.Exit(1)
		}

		driverJson, err := json.Marshal(response.NearbyDrivers)
		if err != nil {
			slog.Error("Error marshaling driver updates: ", "error", err)
			continue // we are in a loop
		}
		// then send driver updates to websocket clients
		if err := conn.WriteMessage(websocket.TextMessage, driverJson); err != nil {
			slog.Error("Error sending driver updates to websocket: ", "error", err)
			return
		}
	}
}
