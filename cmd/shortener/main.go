package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
)

const portNumber = "8080"

var urlsMap = make(map[string]string)

func getURLHash(body []byte) string {

	// TODO:
	// var hash uint32 = crc32.ChecksumIEEE([]byte(data))

	hasher := sha1.New()

	// Write the data to the hasher
	hasher.Write(body)

	// Compute the hash value
	hashBytes := hasher.Sum(nil)

	// Convert the hash bytes to a hexadecimal string for easier display
	hashString := fmt.Sprintf("%x", hashBytes)
	substringLength := min(10, len(hashString)-1)

	return hashString[:substringLength]
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		fmt.Println(r.Method)

		// Read the entire body content
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		hash := getURLHash(body)
		urlsMap[hash] = string(body)

		fmt.Fprintf(w, "http://localhost:8080/"+hash)
		fmt.Println("URL:", "http://localhost:8080/"+hash)
	default:
		fmt.Println("Not served method:", r.Method)
	}
}

func handlerId(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		id := r.PathValue("id")

		if id == "" {
			http.Error(w, "ID must be specified", http.StatusBadRequest)
			return
		}

		url := urlsMap[id]

		fmt.Println("ID:", id)
		fmt.Println("URL:", url)
		// Write a response to the client
		// http.Redirect(w, r, "http://localhost:8080", http.StatusTemporaryRedirect)
		// w.Write([]byte(url))
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		fmt.Println("Not served method:", r.Method)
	}
}

func emptyHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, http.HandlerFunc(handlerRoot))
	mux.HandleFunc(`/{id}`, http.HandlerFunc(handlerId))

	// empty handler for automatic browser favicon request
	mux.HandleFunc(`/favicon.ico`, http.HandlerFunc(emptyHandler))

	fmt.Println("Server starting at:", portNumber)

	err := http.ListenAndServe(":"+portNumber, mux)

	if err != nil {
		fmt.Println("Error starting server:", err)
	} else {
		fmt.Println("Server started at:", portNumber)
	}
}
