package main

import (
	"fmt"
	"net/http"
)

type middleware struct {
	logMiddleware   *loggingMiddleware
	routeMiddleware *routeMiddleware
}

type loggingMiddleware struct{}

type routeMiddleware struct {
	router http.Handler
}

func initMiddleware(app *application) *middleware {

	return &middleware{
		logMiddleware: &loggingMiddleware{},
		routeMiddleware: &routeMiddleware{
			router: initRoutes(app),
		},
	}
}

func initRoutes(app *application) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /spy/{id}/", app.getCat)
	router.HandleFunc("POST /spy", app.createCat)
	return router
}

func (middleware *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	middleware.logMiddleware.log(r)
	middleware.routeMiddleware.router.ServeHTTP(w, r)
}

func (logMid *loggingMiddleware) log(r *http.Request) {
	if r.Response == nil {
		fmt.Printf("New request, method: %s, request URI: %s\n", r.Method, r.URL)
	} else {
		fmt.Printf("Response: status - %s", r.Response.Status)
	}

}
