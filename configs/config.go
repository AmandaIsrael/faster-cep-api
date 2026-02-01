package configs

import (
	"os"
	"time"
)

type Config struct {
	BrasilAPIURL string
	ViaCEPURL    string
	Timeout      time.Duration
	Port         string
}

func Load() *Config {
	return &Config{
		BrasilAPIURL: getEnv("BRASILAPI_URL", "https://brasilapi.com.br/api/cep/v1/%s"),
		ViaCEPURL:    getEnv("VIACEP_URL", "http://viacep.com.br/ws/%s/json/"),
		Timeout:      getDuration("TIMEOUT", time.Second),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}
	return duration
}
