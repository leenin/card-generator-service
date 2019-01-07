package server

import (
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"

	"github.com/qiniu/api.v7/storage"
)

type resultData struct {
	img []byte
	url string
}

func service(r *http.Request, onlyURL bool) (rData resultData, err error) {
	// get request body
	if r.Header.Get("Content-Type") != "application/json" {
		err = errors.New("it only accept application/json content")
		return
	}
	rBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// json to model
	m, err := jsonToModel(rBody)
	if err != nil {
		return
	}

	// model json to sha1
	hexString, err := modelToSha1(&m)
	if err != nil {
		return
	}
	key := "card/" + hexString + ".png"

	// get qiniu config
	qcfg := QiniuCfg{}
	err = getQiniuCfg(&qcfg)
	if err != nil {
		return
	}

	// get file info
	fileInfo, err := getFileInfoByKey(&qcfg, key)
	if err != nil {
		return
	}

	domain := qcfg.Domain
	publicAccessURL := storage.MakePublicURL(domain, key)
	rData.url = publicAccessURL

	var resultImg image.Image
	if fileInfo.PutTime == 0 {
		resultImg, err = composeImage(&m)
		// write png to buffer
		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, resultImg)
		rData.img = buffer.Bytes()

		_, err = uploadFile(&qcfg, rData.img, key)
	} else if !onlyURL {
		resultImg, err = getImageFromURL(rData.url)
		if err != nil {
			return
		}
		// write png to buffer
		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, resultImg)
		rData.img = buffer.Bytes()
	}
	return
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
