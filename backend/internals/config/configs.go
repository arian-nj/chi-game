package config

import (
	"fmt"
	"os"
)

type ReleaseMode int

const (
	Release ReleaseMode = 0
	Develop ReleaseMode = 1
)

type JWTConfig struct {
	SecretKey []byte
}

type Config struct {
	Port        int
	ReleaseMode ReleaseMode
	DatabaseURL string
	Jwt         JWTConfig
}

func ParseConfig() (*Config, error) {
	var conf Config

	releaseMode := os.Getenv("RELEASE_MODE")
	if releaseMode == "" {
		return nil, fmt.Errorf("env var RELEASE_MODE is empty")
	}

	if releaseMode == "release" {
		conf.ReleaseMode = Release
	} else {
		conf.ReleaseMode = Develop
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("env var DATABASE_URL is empty")
	}
	conf.DatabaseURL = dbURL

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("env var JWT_SECRET is empty")
	}
	conf.Jwt = JWTConfig{SecretKey: []byte(jwtSecret)}

	return &conf, nil
}
