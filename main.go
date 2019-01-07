package main

import (
	"card-service/server"
	"net/http"
)

func main() {
	http.HandleFunc("/image/png", server.ImagePngHandler)
	http.HandleFunc("/image/url", server.ImageURLHandler)

	http.ListenAndServe(":8000", nil)
}
