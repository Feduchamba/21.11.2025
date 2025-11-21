package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"project/internal/handlers"
	"project/internal/storage"
	"syscall"
	"time"
)

func main() {
	storage := storage.New()

	http.HandleFunc("/", handlers.CheckLinks(storage))
	http.HandleFunc("/pastLinks", handlers.PastLinks(storage))
	http.HandleFunc("/html", handlers.Handler)

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-rootCtx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	log.Println("Server shutdown")
}
