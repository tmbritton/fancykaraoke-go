package main

import (
	"context"
	"fancykaraoke/db"
	"fancykaraoke/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var server *http.Server

func main() {
	dbase := db.GetConnection()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for sig := range channel {
			if sig == syscall.SIGINT {
				server.Shutdown(ctx)
				dbase.Close()
			}
		}
	}()

	if err := db.DoMigrations(dbase); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", handlers.GetIndex)
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// @TODO 404 handler
	// router.HandleFunc("/", handlers.Get404)

	log.Println("Starting server on port 8080")
	server = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
