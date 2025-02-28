package main

import (
	"net/http"
)

type rootHandler func(w http.ResponseWriter, r *http.Request) error

func (fn rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	dbErr, ok := err.(dbError)
	if ok {
		w.WriteHeader(dbErr.code)
		w.Write([]byte(dbErr.reason))
	}
	serverError, ok := err.(serverError)
	if ok {
		w.WriteHeader(serverError.code)
		w.Write([]byte(serverError.reason))
	}
	// proceed logic
}
