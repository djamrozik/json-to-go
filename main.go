package main

import (
	"encoding/json"
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
	http.HandleFunc("/convertJson", func(res http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var convertRequest ConvertRequest
		err := decoder.Decode(&convertRequest)
		if err != nil {
			log.Println("Could not get request body", err)
			http.Error(res, "Could not get request body", http.StatusBadRequest)
			return
		}

		resultingStruct, err := convertJsonToGolang(convertRequest.Json)
		if err != nil {
			log.Println("Error converting to struct", err)
			http.Error(res, err.Error(), http.StatusBadRequest)
		}
		_, _ = res.Write([]byte(resultingStruct))
	})

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

func convertJsonToGolang(jsonStr string) (string, error) {
	var jsonMap map[string]interface{}

	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return "", err
	}

	return "type Generated struct {\n" + "    Info string `json:\"info\"`\n" + "}", nil
}
