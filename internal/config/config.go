package config

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePort         string
	DatabaseName         string
	JWTSecretKey         string
	DatabaseHost         string
	AppPort              string
	EncryptCookieKey     string `json:"encrypt_cookie_key" env:"ENCRYPT_COOKIE_KEY"`
	AccessTokenLifetime  int
	RefreshTokenLifetime int
	UseHttps             bool
	ContextTimeout       time.Duration
}

func LoadConfig() *Config {

	loadEnv()
	return &Config{
		DatabasePort:         getEnv("DB_PORT", "27017"),
		DatabaseName:         getEnv("DB_NAME", "Database"),
		JWTSecretKey:         getEnv("JWT_SECRET_KEY", "secret_key"),
		DatabaseHost:         getEnv("DB_HOST", "localhost"),
		AppPort:              getEnv("APP_PORT", "8080"),
		EncryptCookieKey:     getValidAESKey("ENCRYPT_COOKIE_KEY"),
		AccessTokenLifetime:  parseInt(getEnv("ACCESS_TOKEN_LIFETIME", "15")),
		RefreshTokenLifetime: parseInt(getEnv("REFRESH_TOKEN_LIFETIME", "43200")),
		UseHttps:             parseBool(getEnv("USE_HTTPS", "false")),
		ContextTimeout:       time.Duration(parseInt(getEnv("CONTEXT_TIMEOUT", "10"))) * time.Second,
	}
}

func loadEnv() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		log.Warnf("Error loading .env file: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseBool(s string) bool {
	return s == "true"
}

func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Errorf("Error parsing integer: %v", err)
		return 0
	}
	return i
}

func getValidAESKey(envVar string) string {
	key := getEnv(envVar, "")

	var keyBytes []byte
	if decoded, err := base64.StdEncoding.DecodeString(key); err == nil {
		keyBytes = decoded
	} else {
		keyBytes = []byte(key)
	}

	switch len(keyBytes) {
	case 16, 24, 32:

		return base64.StdEncoding.EncodeToString(keyBytes)
	default:

		newKey := make([]byte, 32)
		if _, err := rand.Read(newKey); err != nil {
			log.Fatalf("Failed to generate random key: %v", err)
		}
		return base64.StdEncoding.EncodeToString(newKey)
	}
}
