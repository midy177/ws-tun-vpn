package handler

import (
	"context"
	"github.com/gobwas/ws"
	"log"
	"net/http"
)

func CreateAcmeAccountHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			log.Printf("[server] failed to upgrade http %v", err)
			return
		}
	}
}
