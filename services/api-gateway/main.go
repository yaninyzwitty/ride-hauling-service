package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/yaninyzwitty/ride-hauling-app/services/api-gateway/grpc_clients/driver_client"
	"github.com/yaninyzwitty/ride-hauling-app/services/api-gateway/grpc_clients/trip_client"
)

var HttpAddr = ":8081"

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

func main() {
	// initialize the gRPC clients
	driverClient, err := driver_client.NewDriverServiceClient()
	if err != nil {
		slog.Error("failed to create driver service client", "error", err)
		return
	}

	defer driverClient.Close()

	tripClient, err := trip_client.NewTripServiceClient()
	if err != nil {
		slog.Error("failed to create trip service client", "error", err)
		return
	}

	defer tripClient.Close()

	// handle get services, for testing ie k8

	// http.Handle("/services", EnableCors(HandleGetServices))
	// Live stream driver locations
	http.HandleFunc("/ws/drivers", func(w http.ResponseWriter, r *http.Request) {
		handleDriversWebSocket(w, r, driverClient)
	})

	http.HandleFunc("/trip", EnableCors(func(w http.ResponseWriter, r *http.Request) {
		handleCreateTrip(w, r, tripClient)
	}))

	log.Printf("Starting api gateway on port %s", HttpAddr)
	log.Fatal(http.ListenAndServe(HttpAddr, nil))

}
