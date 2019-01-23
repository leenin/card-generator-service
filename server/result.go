package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type result struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *result) success() *result {
	return r
}

func (r *result) successData(data interface{}) *result {
	r.Data = data
	return r
}

func (r *result) fail(code int, message string) *result {
	r.Code = code
	r.Message = message
	return r
}

func (r *result) responseWrite(w http.ResponseWriter) {
	json, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(json)
}
