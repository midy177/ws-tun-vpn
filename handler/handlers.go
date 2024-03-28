package handler

import (
	"context"
	"github.com/gobwas/ws"
	"log"
	"net/http"
	"strconv"
	"ws-tun-vpn/logic"
)

// Handler is a struct that contains the logic of the server.
type Handler struct {
	logic *logic.ServerLogic
}

// NewHandler creates a new Handler instance.
func NewHandler(ctx context.Context) *Handler {
	serverLogic, err := logic.NewServerLogic(ctx)
	if err != nil {
		log.Fatalf("[server] failed to create serverLogic %v", err)
		return nil
	}
	// start the server tunnel packet route to client
	go serverLogic.ServerTunPacketRouteToClient()
	return &Handler{
		logic: serverLogic,
	}
}

// AcceptConnectHandler is a handler for accepting websocket connections.
func (h *Handler) AcceptConnectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCode := r.Header.Get("Authorization")
		if !h.logic.Authenticate(authCode) {
			response := "unauthorized"
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(response))
			return
		}
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("[server] failed to upgrade http: %v", err)
			return
		}
		defer conn.Close()
		// Create a new connection and distribute a CIDR to it.
		addr, err := h.logic.DistributeCIDR(conn)
		// recycle the cidr when the connection is closed
		defer h.logic.RecycleCIDR(addr)
		if err != nil {
			log.Printf("[server] failed to distribute cidr: %v", err)
			return
		}
		err = h.logic.DistributeRote(conn)
		if err != nil {
			log.Printf("[server] failed to distribute route: %v", err)
			return
		}
		// handle the connection
		err = h.logic.HandleConnection(conn)
		if err != nil {
			log.Printf("[server] failed to handle connection: %v", err)
			return
		}
	}
}
