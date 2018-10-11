package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/rvedam/go-password-service/server"
)

func main() {
	stop := make(chan bool, 1)

	srv := &http.Server{Addr: ":8080", Handler: http.DefaultServeMux}
	go func() {
		<-stop
		log.Println("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not shutdown gracefully: %v\n", err)
		}
	}()
	http.HandleFunc("/hash", server.ComputePasswordHash)
	http.HandleFunc("/shutdown", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, http.StatusOK)
		stop <- true
	})
	log.Fatal(srv.ListenAndServe())

}
