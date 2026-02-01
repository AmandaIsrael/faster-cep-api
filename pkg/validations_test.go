package pkg

import "testing"

func TestValidCEP(t *testing.T) {
	validCEPs := []string{
		"12345678",
		"00000000",
		"98765432",
	}

	for _, cep := range validCEPs {
		if !IsValidCEP(cep) {
			t.Errorf("Expected CEP %s to be valid", cep)
		}
	}
}

func TestInvalidCEP(t *testing.T) {
	invalidCEPs := []string{
		"1234",
		"abcdefgh",
		"1234567a",
		"123456789",
		"12-345678",
		"",
	}

	for _, cep := range invalidCEPs {
		if IsValidCEP(cep) {
			t.Errorf("Expected CEP %s to be invalid", cep)
		}
	}
}