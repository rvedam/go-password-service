package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/rvedam/go-password-service/server"
)

func main() {
	portPtr := flag.Int("port", 8080, "specify a port number to listen on")
	helpPtr := flag.Bool("help", false, "generate help")

	flag.Parse()

	if *helpPtr == true {
		flag.PrintDefaults()
		return
	}

	stop := make(chan bool)
	srv := &http.Server{Addr: ":" + strconv.Itoa(*portPtr), Handler: server.NewServer(stop)}
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			log.Printf("HTTP server ListenAndServe: %v", err)
		}
	}()

	fmt.Println("Server listening on 0.0.0.0:" + strconv.Itoa(*portPtr))
	<-stop
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Could not shutdown gracefully: %v\n", err)
	}

	fmt.Println("Server shutdown complete")
}
