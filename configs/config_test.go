package configs

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigWithDefaults(t *testing.T) {
	os.Clearenv()

	config := Load()

	assert.Equal(t, "https://brasilapi.com.br/api/cep/v1/%s", config.BrasilAPIURL)
	assert.Equal(t, "http://viacep.com.br/ws/%s/json/", config.ViaCEPURL)
	assert.Equal(t, time.Second, config.Timeout)
	assert.Equal(t, "8080", config.Port)
}

func TestLoadConfigWithEnvVars(t *testing.T) {
	os.Setenv("BRASILAPI_URL", "https://custom-brasil-api.com/%s")
	os.Setenv("VIACEP_URL", "https://custom-viacep.com/%s")
	os.Setenv("TIMEOUT", "5s")
	os.Setenv("PORT", "3000")

	defer func() {
		os.Clearenv()
	}()

	config := Load()

	assert.Equal(t, "https://custom-brasil-api.com/%s", config.BrasilAPIURL)
	assert.Equal(t, "https://custom-viacep.com/%s", config.ViaCEPURL)
	assert.Equal(t, time.Second*5, config.Timeout)
	assert.Equal(t, "3000", config.Port)
}

func TestGetEnvWithDefaultValue(t *testing.T) {
	os.Clearenv()

	result := getEnv("NON_EXISTENT_KEY", "default_value")

	assert.Equal(t, "default_value", result)
}

func TestGetEnvWithExistingValue(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := getEnv("TEST_KEY", "default_value")

	assert.Equal(t, "test_value", result)
}

func TestGetDurationWithDefaultValue(t *testing.T) {
	os.Clearenv()

	result := getDuration("NON_EXISTENT_TIMEOUT", time.Minute)

	assert.Equal(t, time.Minute, result)
}

func TestGetDurationWithValidValue(t *testing.T) {
	os.Setenv("TEST_TIMEOUT", "10s")
	defer os.Unsetenv("TEST_TIMEOUT")

	result := getDuration("TEST_TIMEOUT", time.Minute)

	assert.Equal(t, time.Second*10, result)
}

func TestGetDurationWithInvalidValue(t *testing.T) {
	os.Setenv("TEST_TIMEOUT", "invalid_duration")
	defer os.Unsetenv("TEST_TIMEOUT")

	result := getDuration("TEST_TIMEOUT", time.Minute)

	assert.Equal(t, time.Minute, result)
}

func TestGetDurationWithEmptyValue(t *testing.T) {
	os.Setenv("TEST_TIMEOUT", "")
	defer os.Unsetenv("TEST_TIMEOUT")

	result := getDuration("TEST_TIMEOUT", time.Minute)

	assert.Equal(t, time.Minute, result)
}
