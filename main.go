package main

import (
	"fmt"
	"log"
	"net/http"
)

type application struct {
	db *sqlDB
}

func main() {
	app := &application{
		db: openAndMigrateDB(),
	}
	initValidBreedsMap()
	middleware := initMiddleware(app)

	server := http.Server{
		Addr:    ":8090",
		Handler: middleware,
	}
	fmt.Println("Server is running on localhost:8090")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
