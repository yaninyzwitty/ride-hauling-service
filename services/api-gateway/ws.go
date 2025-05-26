package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/yaninyzwitty/ride-hauling-app/services/api-gateway/grpc_clients/driver_client"
	pb "github.com/yaninyzwitty/ride-hauling-app/shared/proto/driver"
)

func handleDriversWebSocket(w http.ResponseWriter, r *http.Request, driverClient *driver_client.DriverServiceClient) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("websocket upgrade failed", "error", err)
		return
	}

	defer func(conn *websocket.Conn) {
		conn.Close()
	}(conn)

	slog.Info("WebSocket connection established", "remoteaddr", r.RemoteAddr)

	// add channels for listening to client messages ie user location updates
	clientLocations := make(chan *pb.Location)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					slog.Info("client disconnected", "remote_addr", r.RemoteAddr)
				} else {
					slog.Error("error reading message from websocket", "error", err)
				}
				close(clientLocations)
				return
			}

			var location pb.Location
			if err := json.Unmarshal(msg, &location); err != nil {
				slog.Error("Invalid location data", "error", err)
				continue
			}

			clientLocations <- &location
		}
	}()

	// start streaming driver updates
	stream, err := driverClient.FindNearbyDrivers(r.Context())
	if err != nil {
		slog.Error("failed to start grpc stream", "error", err)
		return
	}
	defer stream.CloseSend()

	// handle sending client location updates to driver service
	go func() {
		for loc := range clientLocations {
			err := stream.Send(&pb.FindNearbyDriversRequest{
				Location: loc,
			})
			if err != nil {
				slog.Error("failed to send location update to driver service", "error", err)
				return
			}

			// Optional: Throttle the rate of location updates
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// receive driver updates from driver service and forward them to client
	for {
		res, err := stream.Recv()
		if err != nil {
			slog.Error("error receiving driver updates", "error", err)
			return
		}

		driverJson, err := json.Marshal(res.NearbyDrivers)
		if err != nil {
			slog.Error("error marshalling driver updates", "error", err)
			return
		}

		// send driver updates using websocket
		if err := conn.WriteMessage(websocket.TextMessage, driverJson); err != nil {
			slog.Error("error sending driver updates", "error", err)
			return
		}
	}
}
