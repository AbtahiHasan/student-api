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
)

func main() {
	// load config 
	cfg := config.MustLoad()

	// start server
	router := http.NewServeMux();

	router.HandleFunc("POST /api/students", student.New())

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

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}