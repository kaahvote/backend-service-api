package main

import (
	"net/http"
	"time"

	"github.com/kaahvote/backend-service-api/internal/data"
)

func (app *application) getSessionHandler(w http.ResponseWriter, r *http.Request) {

	sessionPublicId := app.readStringParam(r, "session_public_id")

	if len(sessionPublicId) == 0 {
		app.notFoundResponse(w, r)
	}

	session := data.Session{
		ID:                 1,
		Name:               "Eleição do representante de turma - 2026",
		PublicID:           sessionPublicId,
		ExpiresAt:          time.Now(),
		VotingPolicyID:     1,
		VotersPolicyID:     1,
		CandidatesPolicyID: 1,
		CreatedBy:          1,
		CreatedAt:          time.Now(),
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"session": session}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
