package handlers

import (
	"context"

	"github.com/AmandaIsrael/faster-cep-api/internal/dto"
	"github.com/stretchr/testify/mock"
)

type MockCEPGateway struct {
	mock.Mock
}

func (m *MockCEPGateway) GetBrasilAPICEP(ctx context.Context, cep string) (*dto.CEP, error) {
	// Mock para BrasilAPI
	args := m.Called(ctx, cep)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CEP), args.Error(1)
}

func (m *MockCEPGateway) GetViaCEP(ctx context.Context, cep string) (*dto.CEP, error) {
	// Mock para ViaCEP API
	args := m.Called(ctx, cep)
	if result := args.Get(0); result == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.CEP), args.Error(1)
}
