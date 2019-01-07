package main

import (
	"card-service/server"
	"net/http"
)

func main() {
	http.HandleFunc("/image", server.ImageHandler)
	http.HandleFunc("/image2qiniu", server.ImageToQiniuHandler)

	http.ListenAndServe(":8000", nil)
}
