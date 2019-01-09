package main

import (
	"net/http"

	"github.com/leenin/card-generator-service/server"
)

func main() {
	http.HandleFunc("/image/png", server.ImagePngHandler)
	http.HandleFunc("/image/url", server.ImageURLHandler)

	http.ListenAndServe(":8000", nil)
}
