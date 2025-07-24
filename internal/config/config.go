package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env                string `yaml:"env" env-default:"local"`
	StoragePostgresCon string `yaml:"storage_path" env-required:"true"`
	MigrationsPath     string `yaml:"migrations_path" env-default:"./migrations"`
	MigrationsTestPath string `yaml:"migrations_test_path" env-default:"./tests/migrations"`
	Server             Server `yaml:"server"`
}

type Server struct {
	Port    int           `yaml:"port" env-default:"8081"`
	Address string        `yaml:"address" env-default:"localhost"`
	Timeout time.Duration `yaml:"timeout" env-default:"60s"`
}

func MustLoad() *Config {
	ptr := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	cfg_path := *ptr
	if cfg_path == "" {
		cfg_path = os.Getenv("CONFIG_PATH")
		if cfg_path == "" {
			log.Fatal("Configuration file path must be provided via -config flag or CONFIG_PATH environment variable")
		}
	}

	if _, err := os.Stat(cfg_path); err != nil {
		log.Fatalf("Configuration file does not exist: %s", cfg_path)
	}

	var cfg Config
	err := cleanenv.ReadConfig(cfg_path, &cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	return &cfg
}
