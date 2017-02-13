package main

import (
	"io"
	"log"
	"net/http"

	"encoding/json"

	"github.com/enrichman/gofield"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/reduce", reduceHandler).Methods("POST")
	log.Println("Ready")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func reduceHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var obj interface{}
	err := decoder.Decode(&obj)
	switch {
	case err == io.EOF:
		log.Println("Empty body")
	case err != nil:
		panic(err)
	}
	defer r.Body.Close()

	log.Println("Received:", obj)

	w.Header().Set("Content-Type", "application/json")

	fields := r.URL.Query().Get("fields")
	if fields == "" {
		log.Println("Missing fields")
		json.NewEncoder(w).Encode(obj)

	} else {
		log.Println("Received fields", fields)
		json.NewEncoder(w).Encode(gofield.Reduce(obj, fields))
	}
}
