package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kaahvote/backend-service-api/internal/data"
	"github.com/kaahvote/backend-service-api/internal/validator"
)

func (app *application) getSessionHandler(w http.ResponseWriter, r *http.Request) {

	sessionPublicId := app.readStringParam(r, "session_public_id")
	if sessionPublicId == "" {
		app.notFoundResponse(w, r)
		return
	}

	session, err := app.models.Sessions.GetFullDetail(sessionPublicId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) postSessionHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name               string        `json:"name"`
		VotingPolicyID     int64         `json:"votingPolicyId"`
		VotersPolicyID     int64         `json:"votersPolicyId"`
		CandidatesPolicyID int64         `json:"candidatesPolicyId"`
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
		errors := make(map[string]string)
		errors["expiresAt"] = err.Error()
		app.failedValidationResponse(w, r, errors)
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

	v := validator.New()

	if data.ValidateSession(v, session); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
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

func (app *application) updateSessionHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	var input struct {
		Name               *string        `json:"name"`
		VotingPolicyID     *int64         `json:"votingPolicyId"`
		VotersPolicyID     *int64         `json:"votersPolicyId"`
		CandidatesPolicyID *int64         `json:"candidatesPolicyId"`
		ExpiresAt          *data.DateTime `json:"expiresAt"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		session.Name = *input.Name
	}

	if input.VotingPolicyID != nil {
		session.VotingPolicyID = *input.VotingPolicyID
	}

	if input.VotersPolicyID != nil {
		session.VotersPolicyID = *input.VotersPolicyID
	}

	if input.CandidatesPolicyID != nil {
		session.CandidatesPolicyID = *input.CandidatesPolicyID
	}

	if input.ExpiresAt != nil {
		dt, err := input.ExpiresAt.ToTime()
		if err != nil {
			errors := make(map[string]string)
			errors["expiresAt"] = err.Error()
			app.failedValidationResponse(w, r, errors)
			return
		}

		session.ExpiresAt = dt
	}

	v := validator.New()

	if data.ValidateSession(v, session); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Sessions.Update(session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) deleteSessionHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	err = app.models.Sessions.Delete(session.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getSession(r *http.Request) (*data.Session, error) {

	sessionPublicId := app.readStringParam(r, "session_public_id")

	if sessionPublicId == "" {
		return nil, data.ErrRecordNotFound
	}

	session, err := app.models.Sessions.Get(sessionPublicId)
	if err != nil {
		return nil, err
	}

	return session, nil

}
