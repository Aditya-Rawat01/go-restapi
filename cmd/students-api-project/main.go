package main

import (
	"log"
	"log/slog"

	"github.com/Aditya-Rawat01/go-restapi/internal/config"
	"github.com/Aditya-Rawat01/go-restapi/internal/storage/sqlite"
	"github.com/Aditya-Rawat01/go-restapi/server"
)

func main() {
	cfg := config.MustLoad()
	storage, err := sqlite.New(cfg)
	if err != nil  {
		log.Fatalf("Error occured during db connection %s", err.Error())
	}

	slog.Info("storage initialized")
	server.ServerHandling(cfg, storage)
}
