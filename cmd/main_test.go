package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaIsrael/faster-cep-api/configs"
	"github.com/stretchr/testify/assert"
)

func TestSetupServer(t *testing.T) {
	config := &configs.Config{
		BrasilAPIURL: "https://brasilapi.com.br/api/cep/v1/%s",
		ViaCEPURL:    "http://viacep.com.br/ws/%s/json/",
		Timeout:      time.Second * 5,
		Port:         "8000",
	}

	server := setupServer(config)

	assert.NotNil(t, server)
	assert.Implements(t, (*http.Handler)(nil), server)
}

func TestSetupServerRoutes(t *testing.T) {
	config := &configs.Config{
		BrasilAPIURL: "https://brasilapi.com.br/api/cep/v1/%s",
		ViaCEPURL:    "http://viacep.com.br/ws/%s/json/",
		Timeout:      time.Second * 5,
		Port:         "8000",
	}

	server := setupServer(config)

	req := httptest.NewRequest("GET", "/01310100", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, req)
	assert.NotEqual(t, http.StatusNotFound, recorder.Code)
}
