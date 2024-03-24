package handler

import (
	"context"
	"github.com/gobwas/ws"
	"log"
	"net/http"
	"strconv"
	"ws-tun-vpn/logic"
)

type Handler struct {
	logic *logic.ServerLogic
}

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
		authCode := r.Header.Get("authorization")
		if !h.logic.Authenticate(authCode) {
			response := "unauthorized"
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(response))
			return
		}
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		// Create a new connection and distribute a CIDR to it.
		addr, err := h.logic.DistributeCIDR(conn)
		if err != nil {
			return
		}
		// recycle the cidr when the connection is closed
		defer h.logic.RecycleCIDR(addr)
		// handle the connection
		err = h.logic.HandleConnection(conn)
		if err != nil {
			log.Printf("[server] failed to upgrade http %v", err)
			return
		}
	}
}
