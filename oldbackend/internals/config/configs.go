package config

import (
	"fmt"
	"os"
)

type ReleasMode int

const (
	Release ReleasMode = 0
	Develop ReleasMode = 1
)

type JWTConfig struct {
	SecretKey []byte
}

type Config struct {
	// BotToken    string
	Port        int
	ReleaseMode ReleasMode
	DatabseUrl  string
	// BaseUrl     string
	Jwt JWTConfig
}

func ParseConfig() (*Config, error) {
	var conf Config

	release_mode := os.Getenv("RELEASE_MODE")
	if release_mode == "" {
		return nil, fmt.Errorf("env var release mode is empty")
	}

	if release_mode == "release" {
		conf.ReleaseMode = Release
	} else {
		conf.ReleaseMode = Develop
	}

	db_url := os.Getenv("DATABASE_URL")
	if db_url == "" {
		return nil, fmt.Errorf("database url is empy")
	}
	conf.DatabseUrl = db_url

	// bot_token := os.Getenv("BOT_TOKEN")
	// if bot_token == "" {
	// 	return nil, fmt.Errorf("bot token is empty")
	// }
	// conf.BotToken = bot_token
	// conf.BotToken = "1346646247:7t2Gzq223R9F6DfA9dK1bcoPbH3xZ3C5xEA"

	conf.Jwt = JWTConfig{
		SecretKey: []byte("palfd34nm3n834n74riq3v4k23avq65b5q7n7n45nw7nwsp8w5b434v5bqb56q35v5b6n7opawrr"),
	}
	return &conf, nil
}
