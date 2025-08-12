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

	routes.HandlerFunc(http.MethodGet, "/v1/users/:user_public_id/sessions", app.getUserSessionsHandler)

	// Candidates
	routes.HandlerFunc(http.MethodGet, "/v1/sessions/:session_public_id/candidates", app.getSessionCandidatesHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/sessions/:session_public_id/candidates/:candidate_id", app.getSingleCandidatesHandler)
	routes.HandlerFunc(http.MethodPost, "/v1/sessions/:session_public_id/candidates", app.postSessionCandidatesHandler)
	routes.HandlerFunc(http.MethodDelete, "/v1/sessions/:session_public_id/candidates/:candidate_id", app.deleteSessionCandidateHandler)

	// Settings
	routes.HandlerFunc(http.MethodGet, "/v1/voting-policies", app.getVotingPoliciesHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/voter-policies", app.getVoterPoliciesHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/candidate-policies", app.getCandidatePoliciesHandler)
	routes.HandlerFunc(http.MethodGet, "/v1/session-states", app.getSessionStatesHandler)

	return routes
}
