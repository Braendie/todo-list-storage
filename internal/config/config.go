package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	Server      Server `yaml:"server"`
}

type Server struct {
	Port    int           `yaml:"port" env-default:"8081"`
	Address string        `yaml:"address" env-default:"localhost`
	Timeout time.Duration `yaml:"timeout" env-default:"60s"`
}

func MustLoad() *Config {
	cfg_path := flag.String("config", "config/local.yaml", "Path to the configuration file")
	flag.Parse()

	if _, err := os.Stat(*cfg_path); err != nil {
		log.Fatalf("Configuration file does not exist: %s", *cfg_path)
	}

	if *cfg_path == "" {
		log.Fatal("Configuration file path must be provided")
	}

	var cfg Config
	err := cleanenv.ReadConfig(*cfg_path, &cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return &cfg
}
