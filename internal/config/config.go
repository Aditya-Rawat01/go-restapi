package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPSERVER struct {
	Addr string  `yaml:"address"`
}

type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPSERVER  `yaml:"http_server" env-required:"true"`
}

func MustLoad() *Config {

	
	configPath := os.Getenv("CONFIG_PATH")
	
	if configPath == "" {
		flags := flag.String("config", "", "path to config file")
		flag.Parse()
		
		configPath = *flags
		
		if configPath == "" {
			log.Fatalf("config file doesnt exist: %s", configPath)
		}
	}
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Error while reading config %s", configPath)
	}
	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("can not read config file: %s", err.Error())
	}
	return &cfg

}
