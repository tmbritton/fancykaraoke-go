package main

import (
	"fancykaraoke/handlers"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", handlers.GetIndex)
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// @TODO 404 handler
	// router.HandleFunc("/", handlers.Get404) // Catch everything else
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}
