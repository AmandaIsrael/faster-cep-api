package gateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaIsrael/faster-cep-api/configs"
	"github.com/AmandaIsrael/faster-cep-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewCEPGateway(t *testing.T) {
	config := &configs.Config{
		BrasilAPIURL: "https://brasilapi.com.br/api/cep/v1/%s",
		ViaCEPURL:    "http://viacep.com.br/ws/%s/json/",
		Timeout:      time.Second,
	}

	gateway := NewCEPGateway(config)

	assert.NotNil(t, gateway)
	assert.Equal(t, config, gateway.config)
}

func TestCEPGatewayGetBrasilAPICEPSuccess(t *testing.T) {
	mockResponse := entity.BrasilAPICEP{
		Cep:          "01310-100",
		State:        "SP",
		City:         "São Paulo",
		Neighborhood: "Bela Vista",
		Street:       "Avenida Paulista",
		Service:      "correios",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := &configs.Config{
		BrasilAPIURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetBrasilAPICEP(ctx, "01310100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "01310-100", result.Cep)
	assert.Equal(t, "SP", result.Estado)
	assert.Equal(t, "São Paulo", result.Cidade)
	assert.Equal(t, "Bela Vista", result.Bairro)
	assert.Equal(t, "Avenida Paulista", result.Rua)
	assert.Equal(t, "correios", result.Servico)
}

func TestCEPGatewayGetBrasilAPICEPHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := &configs.Config{
		BrasilAPIURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetBrasilAPICEP(ctx, "00000000")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API retornou status 404")
}

func TestCEPGatewayGetBrasilAPICEPInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	config := &configs.Config{
		BrasilAPIURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetBrasilAPICEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCEPGatewayGetBrasilAPICEPContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &configs.Config{
		BrasilAPIURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result, err := gateway.GetBrasilAPICEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCEPGatewayGetViaCEPSuccess(t *testing.T) {
	mockResponse := entity.ViaCEP{
		Cep:         "01310-100",
		Logradouro:  "Avenida Paulista",
		Complemento: "",
		Unidade:     "",
		Bairro:      "Bela Vista",
		Localidade:  "São Paulo",
		Uf:          "SP",
		Estado:      "SP",
		Regiao:      "Sudeste",
		Ibge:        "3550308",
		Gia:         "1004",
		Ddd:         "11",
		Siafi:       "7107",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	config := &configs.Config{
		ViaCEPURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetViaCEP(ctx, "01310100")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "01310-100", result.Cep)
	assert.Equal(t, "Avenida Paulista", result.Logradouro)
	assert.Equal(t, "Bela Vista", result.Bairro)
	assert.Equal(t, "São Paulo", result.Localidade)
	assert.Equal(t, "SP", result.Uf)
	assert.Equal(t, "SP", result.Estado)
	assert.Equal(t, "Sudeste", result.Regiao)
	assert.Equal(t, "3550308", result.Ibge)
	assert.Equal(t, "1004", result.Gia)
	assert.Equal(t, "11", result.Ddd)
	assert.Equal(t, "7107", result.Siafi)
}

func TestCEPGatewayGetViaCEPHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	config := &configs.Config{
		ViaCEPURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetViaCEP(ctx, "00000000")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "API retornou status 500")
}

func TestCEPGatewayGetViaCEPInvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{ invalid json }"))
	}))
	defer server.Close()

	config := &configs.Config{
		ViaCEPURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetViaCEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCEPGatewayGetViaCEPContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	config := &configs.Config{
		ViaCEPURL: server.URL + "/%s",
	}
	gateway := NewCEPGateway(config)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	result, err := gateway.GetViaCEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCEPGatewayGetBrasilAPICEPInvalidURL(t *testing.T) {
	config := &configs.Config{
		BrasilAPIURL: "://invalid-url",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetBrasilAPICEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCEPGatewayGetViaCEPInvalidURL(t *testing.T) {
	config := &configs.Config{
		ViaCEPURL: "://invalid-url",
	}
	gateway := NewCEPGateway(config)

	ctx := context.Background()
	result, err := gateway.GetViaCEP(ctx, "01310100")

	assert.Error(t, err)
	assert.Nil(t, result)
}
