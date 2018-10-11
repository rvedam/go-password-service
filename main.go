package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/rvedam/go-password-service/server"
)

func main() {
	stop := make(chan bool, 1)

	srv := &http.Server{Addr: ":8080", Handler: server.NewServer(stop)}
	go func() {
		<-stop
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not shutdown gracefully: %v\n", err)
		}
	}()
	log.Fatal(srv.ListenAndServe())

}
