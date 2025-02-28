package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func sendResponse(w http.ResponseWriter, status int, body ...[]byte) {
	w.WriteHeader(status)
	if len(body) > 0 {
		w.Write(body[0])
	}
}

func encodeBody(body interface{}) ([]byte, error) {
	res, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func getURLvar(r *http.Request, name string) int {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars[name])
	return id
}

func decodeBody(body io.ReadCloser, destination interface{}) error {
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, destination)
	if err != nil {
		return err
	}
	return nil
}
