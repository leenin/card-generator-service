package server

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"

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
	Color   []int   `json:"color" validate:"min=3,max=4"`
}

type Model struct {
	BaseImage string       `json:"base_image" validate:"nonzero"`
	Images    []ImageParam `json:"images"`
	Texts     []TextParam  `json:"texts"`
}

func jsonToModel(rBody []byte) (m Model, err error) {
	if !json.Valid(rBody) {
		err = errors.New("incorrect json format")
		fmt.Println(string(rBody))
		return
	}
	json.Unmarshal(rBody, &m)
	err = validator.Validate(m)
	return
}

func modelToSha1(m *Model) (hexString string, err error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return
	}
	sha1 := sha1.New()
	sha1.Write(jsonBytes)
	hexString = hex.EncodeToString(sha1.Sum([]byte("")))
	return
}
