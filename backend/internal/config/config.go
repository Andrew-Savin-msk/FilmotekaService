package config

import (
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       string `toml:"port"`
	SessionKey string `toml:"session_key"`
	DbType     string `toml:"db_type"`
	DbPath     string `toml:"db_path"`
	LogLevel   string `toml:"log_level"`
}

func Load() *Config {
	var envParamName = "DOCKER_CONFIG_PATH"
	switch runtime.GOOS {
	case "windows":
		envParamName = "CONFIG_PATH"
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Unable to load .env with error: %s", err)
		}
	}

	configPath := os.Getenv(envParamName)
	if configPath == "" {
		log.Fatal("Enviromental variable doen't exists!")
	}

	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatalf("Error with file stats: %s", err)
	}

	var cfg Config
	_, err = toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalf("Trouble with loading data from config file: %s", err)
	}
	return &cfg
}
