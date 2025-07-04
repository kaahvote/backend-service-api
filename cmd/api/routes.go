package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	routes := httprouter.New()
	routes.HandlerFunc("GET", "/v1/health", app.healthCheckHandler)
	return routes
}
