package config

import (
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Host         string `toml:"host"`
	Login        string `toml:"login"`
	Password     string `toml:"app_password"`
	LogLevel     string `toml:"log_level"`
	MDType       string `toml:"md_type"`
	MailBodyPath string `toml:"body_template_path"`
	BrokerType   string `toml:"broker_type"`
	BrokerURL    string `toml:"broker_url"`
}

func Load() *Config {
	cfgPath := "DOCKER_CONFIG_PATH"
	switch runtime.GOOS {
	case "windows":
		cfgPath = "CONFIG_PATH"
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("unable to load .env file, ended with error: %s", err)
		}
	}

	configPath := os.Getenv(cfgPath)
	if configPath == "" {
		log.Fatal("unable to parse enviromental parameter")
	}

	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatalf("unable to load file, ended with error: %s", err)
	}

	var cfg Config
	_, err = toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalf("unable to umarshal toml config file, ended with error: %s", err)
	}
	return &cfg
}
