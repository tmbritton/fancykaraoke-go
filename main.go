package main

import (
	"context"
	"fancykaraoke/db"
	"fancykaraoke/handlers"
	"fancykaraoke/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	server *http.Server
	dbase  *db.SQLiteStore
)

func main() {
	dbase = db.GetConnection()

	if err := db.DoMigrations(dbase); err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("No command specified, please pass import or serve as command line argument")
	}

	command := os.Args[1]

	listenForShutdown()

	if command == "import" {
		doSongImport()
	}

	if command == "serve" {
		startServer()
	}
}

func doSongImport() {
	utils.ImportSongs()
}

func listenForShutdown() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for sig := range channel {
			if sig == syscall.SIGINT {
				if server != nil {
					server.Shutdown(ctx)
				}
				dbase.Close()
			}
		}
	}()
}

func startServer() {
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
