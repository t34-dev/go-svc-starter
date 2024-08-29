package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/t34-dev/go-svc-starter/internal/servers"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channels for errors and signals
	errChan := make(chan error, 2)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("===========================================")

	go func() {
		errChan <- servers.GrpcServe(ctx)
	}()

	go func() {
		errChan <- servers.SwaggerServe(ctx)
	}()

	time.Sleep(100 * time.Millisecond)
	fmt.Println("===========================================")

	// Waiting for termination signal or error
	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("Error occurred: %v", err)
		}
	case sig := <-sigChan:
		log.Printf("Received signal: %v", sig)
	}

	// Starting graceful shutdown process
	log.Println("==============")
	log.Println("Starting graceful shutdown...")
	log.Println("==============")
	cancel()

	// Creating context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Waiting for both servers to shutdown or timeout
	shutdownErrChan := make(chan error, 2)
	go func() {
		shutdownErrChan <- servers.ShutdownGrpcServe(shutdownCtx)
	}()
	go func() {
		shutdownErrChan <- servers.ShutdownSwaggerServe(shutdownCtx)
	}()

	for i := 0; i < 2; i++ {
		if err := <-shutdownErrChan; err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}

	log.Println("Graceful shutdown completed. Exiting.")
}
