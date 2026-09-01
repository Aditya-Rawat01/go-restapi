package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aditya-Rawat01/go-restapi/internal/config"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(*cfg)
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from the go server 😁"))
	})

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	done := make(chan os.Signal, 1)
	
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		
		slog.Info("Server started", slog.String("address", cfg.Addr))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start the server")
		}
	}()
	<-done

	slog.Info("Shutting down the server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // frees the resources immediately as main ends, doesnt care about timeout. without this, the resources will be free only when timeout ends.

	if err := server.Shutdown(ctx); err != nil { // server will shutdown regardless, this is just server didnt shutdown gracefully as expected (some active req had to suffer)
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
