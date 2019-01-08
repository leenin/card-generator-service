package server

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"

	"github.com/qiniu/api.v7/storage"
)

func service(r *http.Request, onlyURL bool) (rData resultData, err error) {
	// get request body
	if r.Header.Get("Content-Type") != "application/json" {
		fmt.Println(r.Header)
		err = errors.New("it only accept application/json content, but get " + r.Header.Get("Content-Type"))
		return
	}

	rBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

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
	qcfg, err := getQiniuCfg()
	if err != nil {
		return
	}

	// get file info
	fileInfo, err := getFileInfoByKey(qcfg, key)
	if err != nil {
		return
	}

	domain := qcfg.Domain
	publicAccessURL := storage.MakePublicURL(domain, key)
	rData.url = publicAccessURL

	var resultImg image.Image
	if fileInfo == nil {
		resultImg, err = composeImage(&m)
		// write png to buffer
		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, resultImg)
		rData.img = buffer.Bytes()

		_, err = uploadFile(qcfg, rData.img, key)
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
