package server

import (
	"net/http"

	"github.com/thesouldev/goboxd/server/handler"
)

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handler.HealthzHandler)
	mux.HandleFunc("POST /run", handler.RunHandler)
	mux.HandleFunc("POST /runs", handler.RunHandler)
}
