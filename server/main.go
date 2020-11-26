package main

import (
	"encoding/json"
	"json-to-golang/lib"
	"log"
	"net/http"
	"os"
)

const defaultPort = "4000"

type ConvertRequest struct {
	Json string `json:"json"`
}

func main() {
	// any unmatched request is sent to serve files from /build
	fs := http.FileServer(http.Dir("./build"))
	http.Handle("/", fs)

	// handle request to convert JSON to a Golang struct
	http.HandleFunc("/convertJson", handleConvertRequest)

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

func handleConvertRequest(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		log.Println("Got a request at /convertJson, but it was not a POST request")
		http.Error(res, "Request to /convertJson must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var convertRequest ConvertRequest
	err := decoder.Decode(&convertRequest)
	if err != nil {
		log.Println("Could not get request body", err)
		http.Error(res, "Could not get request body", http.StatusBadRequest)
		return
	}

	if convertRequest.Json == "" {
		log.Println("Empty JSON field on request")
		http.Error(res, "JSON field was empty", http.StatusBadRequest)
		return
	}

	jsonStr, err := lib.NewJsonStr(convertRequest.Json)
	if err != nil {
		log.Println("Invalid JSON", err.Error())
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	resultingStruct, err := jsonStr.GetAsGolangString()
	if err != nil {
		log.Println("Error converting to struct", err)
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	_, _ = res.Write([]byte(resultingStruct))
}
