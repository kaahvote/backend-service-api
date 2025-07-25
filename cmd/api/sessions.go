package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kaahvote/backend-service-api/internal/data"
)

func (app *application) getSessionHandler(w http.ResponseWriter, r *http.Request) {

	sessionPublicId := app.readStringParam(r, "session_public_id")

	if sessionPublicId == "" {
		app.notFoundResponse(w, r)
		return
	}

	session, err := app.models.Sessions.Get(sessionPublicId)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) postSessionHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name               string        `json:"name"`
		VotingPolicyID     int64         `json:"votingPolicyID"`
		VotersPolicyID     int64         `json:"votersPolicyID"`
		CandidatesPolicyID int64         `json:"candidatesPolicyID"`
		CreatedBy          int64         `json:"createdBy"`
		ExpiresAt          data.DateTime `json:"expiresAt"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	expiresAt, err := input.ExpiresAt.ToTime()
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	uid, _ := uuid.NewV7()

	session := &data.Session{
		Name:               input.Name,
		PublicID:           uid.String(),
		CreatedAt:          time.Now(),
		ExpiresAt:          expiresAt,
		VotingPolicyID:     input.VotingPolicyID,
		VotersPolicyID:     input.VotersPolicyID,
		CandidatesPolicyID: input.CandidatesPolicyID,
		CreatedBy:          input.CreatedBy,
	}

	err = app.models.Sessions.Insert(session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/sessions/%s", session.PublicID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"session": session}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
