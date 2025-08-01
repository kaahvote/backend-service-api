package main

import (
	"net/http"

	"github.com/kaahvote/backend-service-api/internal/data"
)

func (app *application) postSessionFlowHandler(w http.ResponseWriter, r *http.Request) {

	psi := app.readStringParam(r, "session_public_id")
	s, err := app.models.Sessions.Get(psi)

	if err != nil {
		switch err {
		case data.ErrRecordNotFound:
			app.notFoundResponse(w, r)
			return
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	var input struct {
		StateID int64  `json:"state_id"`
		Comment string `json:"comment"`
	}

	app.readJSON(w, r, &input)

	flow := &data.Flow{
		SessionID: s.ID,
		StateID:   input.StateID,
		Comment:   input.Comment,
	}

	err = app.models.Flows.Insert(flow)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"flow": flow}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
