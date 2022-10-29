package app

import (
	"log"
	"net/http"
	. "kafkamongo/internal/publisher/rest-kafka"
)

const (
	POST   = "POST"
)

func (a *App) setRouters() {
	a.Post("/jobs", JobsPostHandler)
	a.Post("/jobs/list", JobsGetHandler)
}

// Wrap the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.StrictSlash(true).HandleFunc(path, f).Methods(POST)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
