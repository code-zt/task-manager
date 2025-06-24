package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabasePort                              string
	DatabaseName                              string
	JWTSecretKey                              string
	DatabaseHost                              string
	UseHttps                                  bool
	AppPort                                   string
	AccessTokenLifetime, RefreshTokenLifetime int
}

func LoadConfig() (*Config, error) {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	err := godotenv.Load(envFile)
	if err != nil {
		return nil, err
	}

	accessTokenLifetime := os.Getenv("ACCESS_TOKEN_LIFETIME")
	refreshTokenLifetime := os.Getenv("REFRESH_TOKEN_LIFETIME")
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	appPort := os.Getenv("APP_PORT")
	useHttps := os.Getenv("USE_HTTPS")

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

	useHttpsBool := parseBool(useHttps)
	accessTokenLifetimeInt, err := strconv.Atoi(accessTokenLifetime)
	if err != nil {
		return nil, err
	}
	refreshTokenLifetimeInt, err := strconv.Atoi(refreshTokenLifetime)
	if err != nil {
		return nil, err
	}
	return &Config{
		DatabasePort:         dbPort,
		JWTSecretKey:         jwtSecretKey,
		DatabaseHost:         dbHost,
		DatabaseName:         dbName,
		UseHttps:             useHttpsBool,
		AppPort:              appPort,
		AccessTokenLifetime:  accessTokenLifetimeInt,
		RefreshTokenLifetime: refreshTokenLifetimeInt,
	}, nil

}

func parseBool(s string) bool {
	if s == "true" {
		return true
	} else {
		return false
	}

}
