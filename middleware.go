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

type responseWriter struct {
	writer http.ResponseWriter
	status int
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
	middleware.logMiddleware.logRequest(r)
	rw := &responseWriter{
		writer: w,
	}
	middleware.routeMiddleware.router.ServeHTTP(rw, r)
	middleware.logMiddleware.logResponse(rw, r)

}

func (logMid *loggingMiddleware) logRequest(r *http.Request) {
	fmt.Printf("New request, method: %s, request URI: %s\n", r.Method, r.URL)
}

func (logMid *loggingMiddleware) logResponse(rw *responseWriter, r *http.Request) {
	fmt.Printf("Response - method: %s, request URI: %s, status - %d\n", r.Method, r.URL, rw.status)
}

func (rw responseWriter) Header() http.Header {
	return rw.writer.Header()
}
func (rw *responseWriter) Write(p []byte) (int, error) {
	n, err := rw.writer.Write(p)
	if err != nil {
		fmt.Println(err)
	}
	return n, nil
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.writer.WriteHeader(code)
	rw.status = code
}
