package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hilmiikhsan/library-book-service/cmd"
	"github.com/hilmiikhsan/library-book-service/helpers"
)

func main() {
	// Setup configuration
	helpers.SetupConfig()

	// Setup logging
	helpers.SetupLogger()

	// Setup PostgreSQL connection
	helpers.SetupPostgres()

	// Setup Redis connection
	helpers.SetupRedis()

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Run gRPC server
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd.ServeGRPC()
	}()

	// Run HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		cmd.ServeHTTP()
	}()

	// Graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	<-signalChan
	log.Println("Shutting down servers...")

	// Wait for all goroutines to finish
	wg.Wait()
	log.Println("All servers stopped. Exiting...")
}
