package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePort                              string
	DatabaseName                              string
	JWTSecretKey                              string
	DatabaseHost                              string
	AppPort                                   string
	AccessTokenLifetime, RefreshTokenLifetime int
	UseHttps                                  bool
	ContextTimeout                            time.Duration
}

func LoadConfig() *Config {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		log.Errorf("Error loading .env file: %v", err)
	}

	accessTokenLifetime := os.Getenv("ACCESS_TOKEN_LIFETIME")
	refreshTokenLifetime := os.Getenv("REFRESH_TOKEN_LIFETIME")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")
	useHttps := os.Getenv("USE_HTTPS")
	contextTimeout := os.Getenv("CONTEXT_TIMEOUT")

	if jwtSecretKey == "" {
		jwtSecretKey = "secret_key"
	}
	if dbPort == "" {
		dbPort = "27017"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbName == "" {
		dbName = "Database"
	}
	if useHttps == "" {
		useHttps = "false"
	}
	if appPort == "" {
		appPort = "8080"
	}
	if contextTimeout == "" {
		contextTimeout = "10"
	}

	useHttpsBool := parseBool(useHttps)

	accessTokenLifetimeInt := parseInt(accessTokenLifetime)
	refreshTokenLifetimeInt := parseInt(refreshTokenLifetime)

	contextTimeoutInt := parseInt(contextTimeout)
	return &Config{
		DatabasePort:         dbPort,
		JWTSecretKey:         jwtSecretKey,
		DatabaseHost:         dbHost,
		DatabaseName:         dbName,
		UseHttps:             useHttpsBool,
		AppPort:              appPort,
		AccessTokenLifetime:  accessTokenLifetimeInt,
		RefreshTokenLifetime: refreshTokenLifetimeInt,
		ContextTimeout:       time.Duration(contextTimeoutInt) * time.Second,
	}

}

func parseBool(s string) bool {
	if s == "true" {
		return true
	} else {
		return false
	}

}
func parseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Errorf("Error parsing: %v", err)
		return 0
	}
	return i
}
