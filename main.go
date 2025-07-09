package main

import (
	"fancykaraoke/db"
	"fancykaraoke/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	dbase, err := db.InitDb()
	if err != nil {
		log.Fatal(err)
	}

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	go func() {
		for sig := range channel {
			// sig is a ^C, handle it
			if sig == syscall.SIGINT {
				dbase.Close()
				log.Fatal("Shutting down")
				// @TODO gracefully shut down http server here
			}
		}
	}()

	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", handlers.GetIndex)
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// @TODO 404 handler
	// router.HandleFunc("/", handlers.Get404)

	log.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
