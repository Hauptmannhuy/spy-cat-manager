package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	app := &application{
		db: openAndMigrateDB(),
	}

	middleware := initMiddleware(app)

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
