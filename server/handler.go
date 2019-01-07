package server

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

func composeImage(r *http.Request) (resultImg *image.RGBA, err error) {
	// get model
	m, err := requestToModel(r)
	if err != nil {
		return
	}

	// get base image
	baseImg, err := getImageFromURL(m.BaseImage)
	if err != nil {
		return
	}
	resultImg = image.NewRGBA(baseImg.Bounds())
	drawBg(resultImg, baseImg)

	// draw image
	var errchs = make([]chan error, len(m.Images)+len(m.Texts))
	for n, imageParam := range m.Images {
		errchs[n] = make(chan error)
		go drawImage(resultImg, imageParam, errchs[n])
	}

	// draw text
	for i, textParam := range m.Texts {
		errchs[len(m.Images)+i] = make(chan error)
		go drawText(resultImg, textParam, errchs[len(m.Images)+i])
	}

	for _, errch := range errchs {
		if err = <-errch; err != nil {
			return
		}
	}

	return
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	// error
	defer func() {
		if err := recover(); err != nil {
			new(result).fail(-1, err.(error).Error()).responseWrite(w)
		}
	}()

	// compose image
	resultImg, err := composeImage(r)
	if err != nil {
		panic(err)
	}

	// write png to buffer
	buffer := new(bytes.Buffer)
	err = png.Encode(buffer, resultImg)
	if err != nil {
		panic(err)
	}

	// response
	w.Header().Add("Content-Type", "image/png")
	w.Write(buffer.Bytes())
}

func ImageToQiniuHandler(w http.ResponseWriter, r *http.Request) {
	// compose image
	resultImg, err := composeImage(r)
	if err != nil {
		panic(err)
	}

	// upload to qiniu
	fmt.Println(resultImg)

	// response
	var data = map[string]string{"url": ""}
	new(result).successData(data)

	// error
	defer func() {
		if err := recover(); err != nil {
			new(result).fail(-1, err.(error).Error()).responseWrite(w)
		}
	}()
}
