package server

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	validator "gopkg.in/validator.v2"
)

type ImageParam struct {
	URL    string `json:"url" validate:"nonzero"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  uint   `json:"width" validate:"nonzero"`
	Height uint   `json:"height" validate:"nonzero"`
	Clip   byte   `json:"clip" validate:"min=0,max=1"`
}

type TextParam struct {
	Content string  `json:"content" validate:"nonzero"`
	Size    float64 `json:"size" validate:"nonzero"`
	Font    string  `json:"font" validate:"nonzero"`
	X       int     `json:"x"`
	Y       int     `json:"y"`
	Anchor  byte    `json:"anchor" validate:"min=0,max=1"`
	Color   []byte  `json:"color" validate:"min=3,max=4"`
}

type Model struct {
	BaseImage string       `json:"base_image" validate:"nonzero"`
	Images    []ImageParam `json:"images"`
	Texts     []TextParam  `json:"texts"`
}

func requestToModel(r *http.Request) (m Model, err error) {
	if r.Header.Get("Content-Type") != "application/json" {
		err = errors.New("it only accept application/json content")
		return
	}
	rBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if !json.Valid(rBody) {
		err = errors.New("incorrect json format")
		return
	}
	json.Unmarshal(rBody, &m)
	err = validator.Validate(m)
	return
}
