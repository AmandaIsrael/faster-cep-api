package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/AmandaIsrael/faster-cep-api/configs"
	"github.com/AmandaIsrael/faster-cep-api/internal/dto"
	"github.com/AmandaIsrael/faster-cep-api/internal/entity"
)

type ICEPGateway interface {
	GetBrasilAPICEP(ctx context.Context, cep string) (*dto.CEP, error)
	GetViaCEP(ctx context.Context, cep string) (*dto.CEP, error)
}

type CEPGateway struct {
	config *configs.Config
}

func NewCEPGateway(config *configs.Config) *CEPGateway {
	return &CEPGateway{
		config: config,
	}
}

func (c *CEPGateway) GetBrasilAPICEP(ctx context.Context, cep string) (*dto.CEP, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(c.config.BrasilAPIURL, cep), nil)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao criar requisição BrasilAPI: %v\n", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao executar requisição BrasilAPI: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CEPGATEWAY] BrasilAPI retornou status %d\n", resp.StatusCode)
		return nil, fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao ler resposta BrasilAPI: %v\n", err)
		return nil, err
	}

	var apiResp entity.BrasilAPICEP
	err = json.Unmarshal(result, &apiResp)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao decodificar JSON BrasilAPI: %v\n", err)
		return nil, err
	}

	log.Println("[CEPGATEWAY] Dados obtidos da BrasilAPI com sucesso")
	return &dto.CEP{
		Cep:     apiResp.Cep,
		Estado:  apiResp.State,
		Cidade:  apiResp.City,
		Bairro:  apiResp.Neighborhood,
		Rua:     apiResp.Street,
		Servico: apiResp.Service,
	}, nil
}

func (c *CEPGateway) GetViaCEP(ctx context.Context, cep string) (*dto.CEP, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf(c.config.ViaCEPURL, cep), nil)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao criar requisição ViaCEP: %v\n", err)
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao executar requisição ViaCEP: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[CEPGATEWAY] ViaCEP retornou status %d\n", resp.StatusCode)
		return nil, fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao ler resposta ViaCEP: %v\n", err)
		return nil, err
	}

	var apiResp entity.ViaCEP
	err = json.Unmarshal(result, &apiResp)
	if err != nil {
		log.Printf("[CEPGATEWAY] Erro ao decodificar JSON ViaCEP: %v\n", err)
		return nil, err
	}

	log.Println("[CEPGATEWAY] Dados obtidos da ViaCEP com sucesso")
	return &dto.CEP{
		Cep:         apiResp.Cep,
		Logradouro:  apiResp.Logradouro,
		Complemento: apiResp.Complemento,
		Unidade:     apiResp.Unidade,
		Bairro:      apiResp.Bairro,
		Localidade:  apiResp.Localidade,
		Uf:          apiResp.Uf,
		Estado:      apiResp.Uf,
		Regiao:      apiResp.Regiao,
		Ibge:        apiResp.Ibge,
		Gia:         apiResp.Gia,
		Ddd:         apiResp.Ddd,
		Siafi:       apiResp.Siafi,
	}, nil
}
