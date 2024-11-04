package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/abtahihasan/students-api/pkg/config"
	"github.com/abtahihasan/students-api/pkg/http/handlers/student"
	"github.com/abtahihasan/students-api/pkg/storage/sqlite"
)

func main() {
	// load config 
	cfg := config.MustLoad()

	// database setup 
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env))

	// start server
	router := http.NewServeMux();

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetStudentList(storage))

	server := http.Server{
		Addr:  cfg.Addr,
		Handler: router,
	}


	slog.Info("server started on :", slog.String("address",cfg.Addr))
	done := make(chan os.Signal, 1)


	signal.Notify(done, os.Interrupt,syscall.SIGINT,syscall.SIGTERM)

	go func ()  {
		err :=server.ListenAndServe()

		if(err != nil) {
			log.Fatal("failed start server")
		}
	}()

	<- done

	// graceful shutdown
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)

	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}