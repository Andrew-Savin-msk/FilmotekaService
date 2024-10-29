package config

import (
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Srv  Server `toml:"server"`
	Bc   Broker `toml:"broker_client"`
	Send Sender `toml:"mail_sender"`
}

type Server struct {
	LogLevel string `toml:"log_level"`
}

type Broker struct {
	BrokerType string `toml:"broker_type"`
	User       string `toml:"user"`
	Password   string `toml:"password"`
	Host       string `toml:"host"`
}

type Sender struct {
	Host         string `toml:"host"`
	Login        string `toml:"login"`
	Password     string `toml:"app_password"`
	MDType       string `toml:"md_type"`
	MailBodyPath string `toml:"body_template_path"`
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
		log.Fatal("unable to parse environmental parameter")
	}

	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatalf("unable to load file, ended with error: %s", err)
	}

	var cfg Config
	_, err = toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalf("unable to unmarshal toml config file, ended with error: %s", err)
	}
	return &cfg
}
