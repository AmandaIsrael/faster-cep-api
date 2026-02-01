package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AmandaIsrael/faster-cep-api/configs"
	"github.com/AmandaIsrael/faster-cep-api/internal/dto"
	"github.com/AmandaIsrael/faster-cep-api/internal/infra/gateway"
	"github.com/AmandaIsrael/faster-cep-api/pkg"
	"github.com/go-chi/chi/v5"
)

type CepHandler struct {
	ICEPGateway gateway.ICEPGateway
	config      *configs.Config
}

func NewCepHandler(cepGateway gateway.ICEPGateway, config *configs.Config) *CepHandler {
	return &CepHandler{
		ICEPGateway: cepGateway,
		config:      config,
	}
}

func (h *CepHandler) GetCEP(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		http.Error(w, "CEP é obrigatório", http.StatusBadRequest)
		return
	}

	if !pkg.IsValidCEP(cep) {
		http.Error(w, "CEP deve conter exatamente 8 dígitos numéricos", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.config.Timeout)
	defer cancel()

	chanResult := make(chan *dto.APIResponse, 1)
	chanError := make(chan error, 2)

	go func() {
		resp, err := h.ICEPGateway.GetBrasilAPICEP(ctx, cep)
		if err != nil {
			chanError <- err
			return
		}
		chanResult <- &dto.APIResponse{Data: resp, Api: "BrasilAPI"}
	}()

	go func() {
		resp, err := h.ICEPGateway.GetViaCEP(ctx, cep)
		if err != nil {
			chanError <- err
			return
		}
		chanResult <- &dto.APIResponse{Data: resp, Api: "ViaCEP"}
	}()

	select {
	case res := <-chanResult:
		h.logCEPResult(res.Data, res.Api)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res.Data)
	case <-ctx.Done():
		if len(chanError) == 2 {
			http.Error(w, "Erro ao obter CEP de todas as APIs", http.StatusInternalServerError)
		} else {
			http.Error(w, "Tempo de espera esgotado para obter o CEP", http.StatusGatewayTimeout)
		}
	}
}

func (h *CepHandler) logCEPResult(cep *dto.CEP, apiName string) {
	fmt.Printf("=== RESULTADO DA CONSULTA CEP ===\n")
	fmt.Printf("API Utilizada: %s\n", apiName)

	fields := map[string]string{
		"CEP":         cep.Cep,
		"Logradouro":  cep.Logradouro,
		"Complemento": cep.Complemento,
		"Unidade":     cep.Unidade,
		"Bairro":      cep.Bairro,
		"Rua":         cep.Rua,
		"Localidade":  cep.Localidade,
		"UF":          cep.Uf,
		"Cidade":      cep.Cidade,
		"Estado":      cep.Estado,
		"Região":      cep.Regiao,
		"IBGE":        cep.Ibge,
		"GIA":         cep.Gia,
		"DDD":         cep.Ddd,
		"SIAFI":       cep.Siafi,
		"Serviço":     cep.Servico,
	}

	for label, value := range fields {
		if value != "" {
			fmt.Printf("%s: %s\n", label, value)
		}
	}
}
