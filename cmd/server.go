package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"volume.finance/pkg/route"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	http.HandleFunc("/track", route.FlightHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	go startListening(srv, port)

	awaitInterruption(srv)
}

func startListening(srv *http.Server, port string) {
	log.Println("listening on port", port)
	if err := srv.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalln("failed to listen, error:", err.Error())
		}
	}
}

func awaitInterruption(srv *http.Server) {
	deathRay := make(chan os.Signal, 1)
	signal.Notify(deathRay, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	<-deathRay
	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("shutdown error:", err.Error())
	}
}
