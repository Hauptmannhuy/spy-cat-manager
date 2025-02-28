package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
	router := mux.NewRouter()
	router.Handle("/spy/{id}", rootHandler(app.getSpy)).Methods("GET")
	router.Handle("/spy", rootHandler(app.createSpy)).Methods("POST")
	router.Handle("/spy/{id}", rootHandler(app.updateSpy)).Methods("PUT")
	router.Handle("/spy/{id}", rootHandler(app.deleteSpy)).Methods("DELETE")

	router.Handle("/mission/{id}", rootHandler(app.getMission)).Methods("GET")
	router.Handle("/mission", rootHandler(app.createMission)).Methods("POST")
	router.Handle("/mission/{id}", rootHandler(app.updateMission)).Methods("PUT")
	router.Handle("/mission/{id}", rootHandler(app.deleteMission)).Methods("DELETE")
	router.Handle("/mission/{id}/target", rootHandler(app.addTargetToMission)).Methods("POST")
	router.Handle("/mission/{mission_id}/target/{target_id}", rootHandler(app.updateMissionTarget)).Methods("PATCH")
	router.Handle("/mission/{mission_id}/target/{target_id}", rootHandler(app.deleteMissionTarget)).Methods("DELETE")
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
	fmt.Printf("Response - method: %s, request URI: %s, status - %d ", r.Method, r.URL, rw.status)
	if rw.body != "" {
		fmt.Printf("body - %s\n", rw.body)
	} else {
		fmt.Printf("\n")
	}
}

type responseWriter struct {
	writer http.ResponseWriter
	status int
	body   string
}

func (rw responseWriter) Header() http.Header {
	return rw.writer.Header()
}
func (rw *responseWriter) Write(p []byte) (int, error) {
	n, err := rw.writer.Write(p)
	rw.body = string(p)
	if err != nil {
		fmt.Println(err)
	}
	return n, nil
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.writer.WriteHeader(code)
	rw.status = code
}
