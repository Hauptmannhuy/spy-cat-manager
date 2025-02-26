package main

import "net/http"

func (app *application) getCat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app *application) createCat(w http.ResponseWriter, r *http.Request) {}
