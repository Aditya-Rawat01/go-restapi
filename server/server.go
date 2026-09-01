package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aditya-Rawat01/go-restapi/internal/config"
	"github.com/Aditya-Rawat01/go-restapi/internal/storage/sqlite"
	"github.com/Aditya-Rawat01/go-restapi/server/handlers/ping"
	"github.com/Aditya-Rawat01/go-restapi/server/handlers/student"
)

func ServerHandling(cfg *config.Config, storage *sqlite.Sqlite) {
	router := http.NewServeMux()

	router.HandleFunc("GET /ping", ping.New()) // health endpoint
	router.HandleFunc("POST /api/students", student.New(storage)) // create new student
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage)) // get student by id
	router.HandleFunc("GET /api/students/", student.GetList(storage)) // get all students
	router.HandleFunc("PUT /api/students/{id}", student.UpdateStudent(storage)) // update the student by id
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteStudentById(storage)) // delete the student by id
	
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	gracefulShutdownHandling(&server, cfg)
}

func gracefulShutdownHandling(server *http.Server, cfg *config.Config) {
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {

		slog.Info("Server started", slog.String("address", cfg.Addr))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start the server")
		}
	}()
	<-done // blocking call, channel will be filled by signal.Notify

	slog.Info("Shutting down the server gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // frees the resources immediately as main ends, doesnt care about timeout. without this, the resources will be free only when timeout ends.

	if err := server.Shutdown(ctx); err != nil { // server will shutdown regardless, this is just server didnt shutdown gracefully as expected (some active req had to suffer)
		slog.Error("failed to shutdown server gracefully", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")
}
