package main

import (
	"log"
	"net/http"
	"sdt-upload-filters/internal"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", internal.IndexHandler)
	mux.HandleFunc("/upload", internal.UploadHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
