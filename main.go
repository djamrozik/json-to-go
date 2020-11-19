package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
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

// this function takes a json map as map[string]interface{} and will return a list of strings
// representing that map as a golang struct
func buildStringRepList(jsonMap map[string]interface{}, indentLevel int) []string {
	keyReps := make([]string, 0)
	for k, v := range jsonMap {
		typeName := getValueType(v)
		if typeName == "map[string]interface {}" {
			valueConverted := v.(map[string]interface{})
			structKeyReps := buildStructKeyReps(k, valueConverted, indentLevel)
			keyReps = append(keyReps, structKeyReps...)
			continue
		}
		keyRep := buildTerminalKeyRep(k, typeName, indentLevel)
		keyReps = append(keyReps, keyRep)
	}
	return keyReps
}

func buildStructKeyReps(key string, jsonMap map[string]interface{}, indentLevel int) []string {
	structStart := strings.Repeat(" ", indentLevel*4)
	structStart += strings.Title(key)
	structStart += " struct {"

	innerKeyReps := buildStringRepList(jsonMap, indentLevel+1)

	structEnd := strings.Repeat(" ", indentLevel*4)
	structEnd += fmt.Sprintf("} `json:\"%s\"", key)

	keyReps := []string{structStart}
	keyReps = append(keyReps, innerKeyReps...)
	keyReps = append(keyReps, structEnd)
	return keyReps
}

func buildTerminalKeyRep(key string, typeName string, indentLevel int) string {
	keyRep := strings.Repeat(" ", indentLevel*4)
	keyRep += strings.Title(key)
	keyRep += " " + typeName
	keyRep += " " + fmt.Sprintf("`json:\"%s\"`", key)
	return keyRep
}

func convertJsonToGolang(jsonStr string) (string, error) {
	var jsonMap map[string]interface{}

	// convert json string to map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return "", errors.New("JSON is invalid: " + err.Error())
	}

	// build string representation of each key
	keyReps := buildStringRepList(jsonMap, 1)

	// build final string for the golang struct
	finalString := "type Generated struct {\n"
	for _, v := range keyReps {
		finalString += v + "\n"
	}
	finalString += "}"

	return finalString, nil
}

func handleConvertRequest(res http.ResponseWriter, req *http.Request) {
	// only allowing POST requests
	if req.Method != http.MethodPost {
		log.Println("Got a request at /convertJson, but it was not a POST request")
		http.Error(res, "Request to /convertJson must be a POST request", http.StatusMethodNotAllowed)
		return
	}

	// decode request body to a ConvertRequest struct
	// if the 'json' field is not present, will be an empty string in resulting struct
	// other fields in the body are ignored
	// only was this is an error is if the body is not present or body is malformed
	decoder := json.NewDecoder(req.Body)
	var convertRequest ConvertRequest
	err := decoder.Decode(&convertRequest)
	if err != nil {
		log.Println("Could not get request body", err)
		http.Error(res, "Could not get request body", http.StatusBadRequest)
		return
	}

	// some quick validation that the json field is not empty
	if convertRequest.Json == "" {
		log.Println("Empty JSON field on request")
		http.Error(res, "JSON field was empty", http.StatusBadRequest)
		return
	}

	// where the actual work is taking place, converting json string to golang struct string
	// if there is an error, return the error message
	// if successful, then return golang struct string as data
	resultingStruct, err := convertJsonToGolang(convertRequest.Json)
	if err != nil {
		log.Println("Error converting to struct:", err)
		http.Error(res, err.Error(), http.StatusBadRequest) // TODO : differentiate between internal/external error
	}
	_, _ = res.Write([]byte(resultingStruct))
}

func getValueType(v interface{}) string {
	typeName := reflect.TypeOf(v).String()
	if typeName == "float64" {
		value := v.(float64)
		if value == float64(int(value)) {
			typeName = "int"
		}
	}
	return typeName
}
