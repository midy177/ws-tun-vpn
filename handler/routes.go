package handler

import (
	"context"
	"fmt"
	"net/http"
)

// RegisterHandlers registers all the handlers for the server.
func RegisterHandlers(ctx context.Context) {
	Handlle := NewHandler(ctx)
	http.Handle("/connect", recoverHandler(Handlle.AcceptConnectHandler()))
}

// recoverHandler wraps the given http.Handler with a panic recovery mechanism.
func recoverHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from panic:", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
