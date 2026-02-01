package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AmandaIsrael/faster-cep-api/configs"
	"github.com/AmandaIsrael/faster-cep-api/internal/dto"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupHandler(mockGateway *MockCEPGateway) *CepHandler {
	config := &configs.Config{
		Timeout: time.Second * 5,
	}
	return NewCepHandler(mockGateway, config)
}

func createRequest(method, url, cep string) *http.Request {
	req := httptest.NewRequest(method, url, nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("cep", cep)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	return req
}

func TestCepHandlerGetCEPSuccessBrasilAPI(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	expectedCEP := &dto.CEP{
		Cep:     "01310-100",
		Estado:  "São Paulo",
		Cidade:  "São Paulo",
		Bairro:  "Bela Vista",
		Rua:     "Avenida Paulista",
		Servico: "correios",
	}

	mockGateway.On("GetBrasilAPICEP", mock.Anything, "01310100").Return(expectedCEP, nil)
	mockGateway.On("GetViaCEP", mock.Anything, "01310100").Return(&dto.CEP{}, errors.New("timeout"))

	req := createRequest("GET", "/cep/01310100", "01310100")
	recorder := httptest.NewRecorder()

	handler.GetCEP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var response dto.CEP
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCEP.Cep, response.Cep)
	assert.Equal(t, expectedCEP.Cidade, response.Cidade)
	assert.Equal(t, expectedCEP.Bairro, response.Bairro)

	mockGateway.AssertExpectations(t)
}

func TestCepHandlerGetCEPSuccessViaCEP(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	expectedCEP := &dto.CEP{
		Cep:        "01310-100",
		Logradouro: "Avenida Paulista",
		Bairro:     "Bela Vista",
		Localidade: "São Paulo",
		Uf:         "SP",
		Ibge:       "3550308",
		Ddd:        "11",
	}

	mockGateway.On("GetViaCEP", mock.Anything, "01310100").Return(expectedCEP, nil)
	mockGateway.On("GetBrasilAPICEP", mock.Anything, "01310100").Return(&dto.CEP{}, errors.New("timeout"))

	req := createRequest("GET", "/cep/01310100", "01310100")
	recorder := httptest.NewRecorder()

	handler.GetCEP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response dto.CEP
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expectedCEP.Cep, response.Cep)
	assert.Equal(t, expectedCEP.Logradouro, response.Logradouro)
	assert.Equal(t, expectedCEP.Localidade, response.Localidade)

	mockGateway.AssertExpectations(t)
}

func TestCepHandlerGetCEPEmptyCEP(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	req := createRequest("GET", "/cep/", "")
	recorder := httptest.NewRecorder()

	handler.GetCEP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "CEP é obrigatório")

	mockGateway.AssertNotCalled(t, "GetBrasilAPICEP")
	mockGateway.AssertNotCalled(t, "GetViaCEP")
}

func TestCepHandlerGetCEPInvalidCEP(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	invalidCEPs := []string{"123", "12345678a", "123456789", "abcdefgh"}

	for _, invalidCEP := range invalidCEPs {
		t.Run(fmt.Sprintf("CEP_Invalid_%s", invalidCEP), func(t *testing.T) {
			req := createRequest("GET", "/cep/"+invalidCEP, invalidCEP)
			recorder := httptest.NewRecorder()

			handler.GetCEP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "CEP deve conter exatamente 8 dígitos numéricos")
		})
	}

	mockGateway.AssertNotCalled(t, "GetBrasilAPICEP")
	mockGateway.AssertNotCalled(t, "GetViaCEP")
}

func TestCepHandlerGetCEPBothAPIsError(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	mockGateway.On("GetBrasilAPICEP", mock.Anything, "01310100").Return(nil, errors.New("API error"))
	mockGateway.On("GetViaCEP", mock.Anything, "01310100").Return(nil, errors.New("API error"))

	req := createRequest("GET", "/cep/01310100", "01310100")
	recorder := httptest.NewRecorder()

	handler.GetCEP(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Erro ao obter CEP de todas as APIs")

	mockGateway.AssertExpectations(t)
}

func TestCepHandlerGetCEPTimeout(t *testing.T) {
	mockGateway := new(MockCEPGateway)

	config := &configs.Config{
		Timeout: time.Millisecond * 1,
	}

	handler := NewCepHandler(mockGateway, config)

	mockGateway.On("GetBrasilAPICEP", mock.Anything, "01310100").Run(func(args mock.Arguments) {
		time.Sleep(time.Millisecond * 10)
	}).Return(nil, context.DeadlineExceeded)

	mockGateway.On("GetViaCEP", mock.Anything, "01310100").Run(func(args mock.Arguments) {
		time.Sleep(time.Millisecond * 10)
	}).Return(nil, context.DeadlineExceeded)

	req := createRequest("GET", "/cep/01310100", "01310100")
	recorder := httptest.NewRecorder()

	handler.GetCEP(recorder, req)

	assert.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Tempo de espera esgotado para obter o CEP")
}

func TestCepHandlerGetCEPOneAPIErrorOneTimeout(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	mockGateway.On("GetBrasilAPICEP", mock.Anything, "01310100").Return(nil, errors.New("API error"))
	mockGateway.On("GetViaCEP", mock.Anything, "01310100").Run(func(args mock.Arguments) {
		time.Sleep(time.Second * 10)
	}).Return(nil, errors.New("timeout"))

	req := createRequest("GET", "/cep/01310100", "01310100")
	recorder := httptest.NewRecorder()

	handler.GetCEP(recorder, req)

	assert.Equal(t, http.StatusGatewayTimeout, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Tempo de espera esgotado para obter o CEP")

	mockGateway.AssertExpectations(t)
}

func TestCepHandlerNewCepHandler(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	config := &configs.Config{Timeout: time.Second}

	handler := NewCepHandler(mockGateway, config)

	assert.NotNil(t, handler)
	assert.Equal(t, mockGateway, handler.ICEPGateway)
	assert.Equal(t, config, handler.config)
}

func TestCepHandlerGetCEPValidCEPFormats(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	expectedCEP := &dto.CEP{
		Cep:    "01310100",
		Cidade: "São Paulo",
	}

	validCEPs := []string{"01310100", "12345678", "00000000", "99999999"}

	for _, cep := range validCEPs {
		t.Run(fmt.Sprintf("CEP_Valid_%s", cep), func(t *testing.T) {
			mockGateway.On("GetBrasilAPICEP", mock.Anything, cep).Return(expectedCEP, nil).Once()
			mockGateway.On("GetViaCEP", mock.Anything, cep).Return(nil, errors.New("error")).Once()

			req := createRequest("GET", "/cep/"+cep, cep)
			recorder := httptest.NewRecorder()

			handler.GetCEP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
		})
	}

	mockGateway.AssertExpectations(t)
}

func TestCepHandlerGetCEPCepWithHyphenShouldFail(t *testing.T) {
	mockGateway := new(MockCEPGateway)
	handler := setupHandler(mockGateway)

	cepsWithHyphen := []string{"01310-100", "12345-678", "00000-000"}

	for _, cep := range cepsWithHyphen {
		t.Run(fmt.Sprintf("CEP_WithHyphen_%s", cep), func(t *testing.T) {
			req := createRequest("GET", "/cep/"+cep, cep)
			recorder := httptest.NewRecorder()

			handler.GetCEP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
			assert.Contains(t, recorder.Body.String(), "CEP deve conter exatamente 8 dígitos numéricos")
		})
	}

	mockGateway.AssertNotCalled(t, "GetBrasilAPICEP")
	mockGateway.AssertNotCalled(t, "GetViaCEP")
}
