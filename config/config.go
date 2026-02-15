package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Server      ServerConfig
	Auth        AuthConfig
}

type ServerConfig struct {
	Port string
}

type AuthConfig struct {
	ValidTokens []string
}

var AppConfig *Config

// LoadConfig carga la configuración desde variables de entorno
func LoadConfig() error {
	// Intentar cargar el archivo .env (opcional en producción)
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	AppConfig = &Config{
		DatabaseURL: getEnv("DATABASE_URL", ""),
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Auth: AuthConfig{
			ValidTokens: parseTokens(getEnv("AUTH_TOKENS", "")),
		},
	}

	// Validar configuración crítica
	if AppConfig.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL no está configurado")
	}

	if len(AppConfig.Auth.ValidTokens) == 0 {
		return fmt.Errorf("AUTH_TOKENS no está configurado")
	}

	return nil
}

// GetDSN retorna el string de conexión para PostgreSQL
func (c *Config) GetDSN() string {
	return c.DatabaseURL
}

// IsValidToken verifica si un token es válido
func (c *Config) IsValidToken(token string) bool {
	for _, validToken := range c.Auth.ValidTokens {
		if validToken == token {
			return true
		}
	}
	return false
}

// parseTokens divide los tokens separados por coma
func parseTokens(tokensStr string) []string {
	if tokensStr == "" {
		return []string{}
	}
	
	tokens := strings.Split(tokensStr, ",")
	var cleanTokens []string
	
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token != "" {
			cleanTokens = append(cleanTokens, token)
		}
	}
	
	return cleanTokens
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}