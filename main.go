package main

import (
	"log"
	"net/http"
	"os"
)

const defaultPort = "4000"

func main() {
	// any unmatched request is sent to serve files from /build
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	// get port to listen on
	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("PORT env var not set, using default port")
		port = defaultPort
	}
	log.Printf("Listening on :%s...", port)

	// start server
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
