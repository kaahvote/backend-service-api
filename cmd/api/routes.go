package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	routes := httprouter.New()

	routes.NotFound = http.HandlerFunc(app.notFoundResponse)
	routes.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	routes.HandlerFunc("GET", "/v1/health", app.healthCheckHandler)
	routes.HandlerFunc("GET", "/v1/sessions/:session_public_id", app.getSessionHandler)
	return routes
}
