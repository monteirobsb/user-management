package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// AppConfig holds the application configuration
type AppConfig struct {
	DatabaseHost       string
	DatabasePort       string
	PostgresUser       string
	PostgresPassword   string
	PostgresDB         string
	DatabaseSSLMode    string
	JWTSecretKey       string
	JWTExpirationMinutes time.Duration
	APIPort            string
}

// Config is the global application configuration
var Config *AppConfig

func init() {
	Config = LoadConfig()
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *AppConfig {
	// Load .env file from the project root
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("INFO: No .env file found, relying on environment variables")
	}

	appConfig := &AppConfig{}

	appConfig.DatabaseHost = getEnv("DATABASE_HOST", "")
	if appConfig.DatabaseHost == "" {
		log.Fatal("CRITICAL: DATABASE_HOST is not set")
	}

	appConfig.DatabasePort = getEnv("DATABASE_PORT", "")
	if appConfig.DatabasePort == "" {
		log.Fatal("CRITICAL: DATABASE_PORT is not set")
	}

	appConfig.PostgresUser = getEnv("POSTGRES_USER", "")
	if appConfig.PostgresUser == "" {
		log.Fatal("CRITICAL: POSTGRES_USER is not set")
	}

	appConfig.PostgresPassword = getEnv("POSTGRES_PASSWORD", "")
	if appConfig.PostgresPassword == "" {
		log.Fatal("CRITICAL: POSTGRES_PASSWORD is not set")
	}

	appConfig.PostgresDB = getEnv("POSTGRES_DB", "")
	if appConfig.PostgresDB == "" {
		log.Fatal("CRITICAL: POSTGRES_DB is not set")
	}

	appConfig.DatabaseSSLMode = getEnv("DATABASE_SSLMODE", "disable")
	if appConfig.DatabaseSSLMode == "disable" {
		log.Println("INFO: DATABASE_SSLMODE not set, using default 'disable'")
	} else {
		log.Printf("INFO: DATABASE_SSLMODE set to '%s'", appConfig.DatabaseSSLMode)
	}

	appConfig.JWTSecretKey = getEnv("JWT_SECRET_KEY", "")
	if appConfig.JWTSecretKey == "" {
		log.Fatal("CRITICAL: JWT_SECRET_KEY is not set")
	}

	jwtExpirationMinutesStr := getEnv("JWT_EXPIRATION_MINUTES", "1440")
	jwtExpirationMinutes, err := strconv.Atoi(jwtExpirationMinutesStr)
	if err != nil {
		log.Printf("WARNING: Invalid JWT_EXPIRATION_MINUTES value '%s', using default 1440 minutes", jwtExpirationMinutesStr)
		appConfig.JWTExpirationMinutes = 1440 * time.Minute
	} else {
		log.Printf("INFO: JWT_EXPIRATION_MINUTES set to %d minutes", jwtExpirationMinutes)
		appConfig.JWTExpirationMinutes = time.Duration(jwtExpirationMinutes) * time.Minute
	}

	appConfig.APIPort = getEnv("API_PORT", "8080")
	if appConfig.APIPort == "8080" && getEnv("API_PORT", "") == "" { // Only log default if it wasn't explicitly set to 8080
		log.Println("INFO: API_PORT not set, using default '8080'")
	} else {
		log.Printf("INFO: API_PORT set to '%s'", appConfig.APIPort)
	}

	return appConfig
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if defaultValue != "" {
			// Log only for optional variables that have a default and are not set
			// Critical variables without defaults will be handled by the caller (LoadConfig)
			switch key {
			case "DATABASE_SSLMODE", "JWT_EXPIRATION_MINUTES", "API_PORT":
				log.Printf("INFO: %s not set, using default '%s'", key, defaultValue)
			}
		}
		return defaultValue
	}
	return value
}
