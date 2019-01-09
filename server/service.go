package server

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"net/http"

	"github.com/qiniu/api.v7/storage"
)

func service(r *http.Request, onlyURL bool) (rData resultData, err error) {
	rBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}

	// init model
	m := &Model{}

	// json to model
	err = m.fromJSONBytes(rBody)
	if err != nil {
		return
	}

	// model json to sha1
	hexString, err := m.toSha1()
	if err != nil {
		return
	}
	key := "card/" + hexString + ".png"

	// init qiniu
	qn := new(Qiniu)
	qn.init()

	// get file info
	fileInfo := qn.getFileInfoByKey(key)

	domain := qn.Domain
	publicAccessURL := storage.MakePublicURL(domain, key)
	rData.url = publicAccessURL

	var resultImg image.Image
	if fileInfo.PutTime == 0 {
		resultImg, err = composeImage(m)
		// write png to buffer
		buffer := new(bytes.Buffer)
		err = png.Encode(buffer, resultImg)
		rData.img = buffer.Bytes()

		if onlyURL {
			qn.uploadFile(rData.img, key)
		} else {
			go qn.uploadFile(rData.img, key)
		}
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
