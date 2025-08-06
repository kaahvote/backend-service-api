package main

import (
	"net/http"

	"github.com/kaahvote/backend-service-api/internal/data"
	"github.com/kaahvote/backend-service-api/internal/validator"
)

func (app *application) getUserSessionsHandler(w http.ResponseWriter, r *http.Request) {

	userPID := app.readStringParam(r, "user_public_id")
	if userPID == "" {
		app.notFoundResponse(w, r)
	}

	user, err := app.models.Users.Get(userPID)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	v := validator.New()

	qs := r.URL.Query()
	name := app.readString(qs, "name", "")
	expFrom := app.readDate(qs, "expFrom")
	expTo := app.readDate(qs, "expTo")

	crtdFrom := app.readDate(qs, "createdFrom")
	crtdTo := app.readDate(qs, "createdTo")

	votingPolicyID := int64(app.readInt(qs, "votingPolicy", 0, v))
	votersPolicyID := int64(app.readInt(qs, "votersPolicy", 0, v))
	candidatePolicyID := int64(app.readInt(qs, "candidatePolicy", 0, v))

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	filters := data.SessionFilters{
		Name:              name,
		VotingPolicyID:    votingPolicyID,
		VotersPolicyID:    votersPolicyID,
		CandidatePolicyID: candidatePolicyID,
		CreatedBy:         user.ID,
		CreatedAtFrom:     crtdFrom,
		CreatedAtTo:       crtdTo,
		ExpiresAtFrom:     expFrom,
		ExpiresAtTo:       expTo,
	}

	sessions, err := app.models.Sessions.ListSessionsFiltering(filters)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"sessions": sessions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
