package server

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

type resultData struct {
	img []byte
	url string
}

func ImagePngHandler(w http.ResponseWriter, r *http.Request) {
	// error
	defer func() {
		if err := recover(); err != nil {
			new(result).fail(-1, err.(error).Error()).responseWrite(w)
		}
	}()

	rData, err := service(r, false)
	if err != nil {
		panic(err)
	}

	// response
	w.Header().Add("Content-Type", "image/png")
	w.Write(rData.img)
}

func ImageURLHandler(w http.ResponseWriter, r *http.Request) {
	// error
	defer func() {
		if err := recover(); err != nil {
			new(result).fail(-1, err.(error).Error()).responseWrite(w)
		}
	}()

	rData, err := service(r, true)
	if err != nil {
		panic(err)
	}

	// response
	var data = map[string]string{"url": rData.url}
	new(result).successData(data).responseWrite(w)
}

func ApiDocsHandler(w http.ResponseWriter, r *http.Request) {
	// error
	defer func() {
		if err := recover(); err != nil {
			new(result).fail(-1, err.(error).Error()).responseWrite(w)
		}
	}()

	filePath, err := filepath.Abs("./swagger.json")
	if err != nil {
		panic(errors.New("Swagger json font not exists"))
	}

	swaggerJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(swaggerJSON)
}
