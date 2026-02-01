package main

import (
	"net/http"

	"github.com/AmandaIsrael/faster-cep-api/configs"
	"github.com/AmandaIsrael/faster-cep-api/internal/infra/gateway"
	"github.com/AmandaIsrael/faster-cep-api/internal/infra/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	config := configs.Load()
	server := setupServer(config)
	http.ListenAndServe(":"+config.Port, server)
}

func setupServer(config *configs.Config) http.Handler {
	cepGateway := gateway.NewCEPGateway(config)
	cepHandler := handlers.NewCepHandler(cepGateway, config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/{cep}", cepHandler.GetCEP)

	return r
}
