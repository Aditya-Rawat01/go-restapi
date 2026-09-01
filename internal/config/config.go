package config

import (
	"flag"
	"log"
	"os"
)

type HTTPSERVER struct {
	Addr string
}

type Config struct {
	Env         string `yml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yml: "storage_path" env-required:"true"`
	HTTPSERVER `yml: "http_server"`
}

func MustLoad() *Config {
	var cfg Config

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatalf("config file doesnt exist: %s", configPath)
		}
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Error while reading config %s", err.Error())
	}

	return &cfg

}