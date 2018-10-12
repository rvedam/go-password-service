package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rvedam/go-password-service/server"
)

func main() {
	stop := make(chan bool)
	sigTerm := make(chan os.Signal, 1)

	signal.Notify(sigTerm, os.Interrupt, syscall.SIGTERM)
	srv := &http.Server{Addr: ":8080", Handler: server.NewServer(stop)}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	go func() {
		<-stop
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not shutdown gracefully: %v\n", err)
		}
	}()

	<-sigTerm

	fmt.Println("Server shutdown complete")
}
