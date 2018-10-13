package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/rvedam/go-password-service/server"
)

func main() {
	stop := make(chan bool)
	srv := &http.Server{Addr: ":8080", Handler: server.NewServer(stop)}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown gracefully: %v\n", err)
	}

	fmt.Println("Server shutdown complete")
}
