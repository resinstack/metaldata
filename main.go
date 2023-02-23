package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"

	"github.com/resinstack/metaldata/pkg/http"
	"github.com/resinstack/metaldata/pkg/source/fs"
)

func main() {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  "metaldata",
		Level: hclog.LevelFromString("TRACE"),
	})

	s := http.New(
		http.WithLogger(appLogger),
		http.WithInfoSource(fs.New(os.Getenv("MD_FSBASE"))),
	)

	bind := os.Getenv("MD_BIND")
	if bind == "" {
		bind = ":3030"
	}

	go func() {
		s.Serve(bind)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	appLogger.Info("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		appLogger.Error("Error during shutdown", "error", err)
		os.Exit(2)
	}
}
