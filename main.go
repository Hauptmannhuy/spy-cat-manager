package main

import (
	"fmt"
	"log"
	"net/http"
)

type application struct {
	middleware *middleware
}

func main() {
	app := &application{}
	middleware := initMiddleware(app)
	app.middleware = middleware

	server := http.Server{
		Addr:    "localhost:8090",
		Handler: middleware,
	}
	fmt.Println("Server is running on localhost:8090")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func (app *application) respond(w http.ResponseWriter, r *http.Request) {

}
