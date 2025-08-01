package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	routes := httprouter.New()

	routes.NotFound = http.HandlerFunc(app.notFoundResponse)
	routes.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	routes.HandlerFunc(http.MethodGet, "/v1/health", app.healthCheckHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/sessions/:session_public_id", app.getSessionHandler)
	routes.HandlerFunc(http.MethodPost, "/v1/sessions", app.postSessionHandler)

	routes.HandlerFunc(http.MethodPatch, "/v1/sessions/:session_public_id", app.updateSessionHandler)
	routes.HandlerFunc(http.MethodDelete, "/v1/sessions/:session_public_id", app.deleteSessionHandler)

	routes.HandlerFunc(http.MethodPost, "/v1/sessions/:session_public_id/flows", app.postSessionFlowHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/sessions/:session_public_id/flows", app.getSessionFlowHandler)

	return routes
}
