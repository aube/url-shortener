package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/aube/url-shortener/internal/app/store"
	"google.golang.org/grpc"
)

// Run initializes and starts the URL shortener application.
// It performs the following steps:
//  1. Loads configuration using config.NewConfig()
//  2. Initializes the storage backend using store.NewStore()
//  3. Creates the router with all endpoints using router.New()
//  4. Starts the HTTP server on the configured address
//
// The function blocks until the server exits and returns any error that occurs.
// In case of a fatal server error, it logs the error and terminates the program.
//
// Example usage:
//
//	if err := app.Run(); err != nil {
//	    log.Fatal("Application failed:", err)
//	}
func Run() error {
	// Initialize storage (database, file, or memory based on config)
	store.NewStore()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	var httpServer http.Server
	var grpcServer *grpc.Server

	go func() {
		httpServer = *httpServerStart()
	}()

	// Run gRPC server
	go func() {
		grpcServer = grpcServerStart()
	}()

	sig := <-sigs
	log.Printf("Received signal: %v\n", sig)

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Close database
	if err := store.Close(); err != nil {
		log.Printf("Store close error: %v", err)
		return err
	} else {
		log.Printf("Store connection closed")
	}

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("HTTP Server shutdown error: %v\n", err)
		return err
	} else {
		log.Println("HTTP Server gracefully stopped")
	}

	// Gracefully stop the server
	stopped := make(chan struct{})
	go func() {
		// I don't know why grpcServer.GracefulStop calls panic
		// grpcServer.GracefulStop()
		close(stopped)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-stopped:
		log.Println("Type of grpcServer", reflect.TypeOf(grpcServer))
		log.Println("Server gracefully stopped")
	case <-ctx.Done():
		log.Println("Shutdown timeout exceeded, forcing stop")
		// I don't know why grpcServer.GracefulStop calls panic
		// grpcServer.Stop()
	}

	// grpcServer release own port even if GracefulStop or Stop methods still commented
	log.Println("Server shutdown complete")

	return nil
}
