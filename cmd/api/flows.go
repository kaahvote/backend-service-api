package main

import (
	"net/http"

	"github.com/kaahvote/backend-service-api/internal/data"
)

func (app *application) postSessionFlowHandler(w http.ResponseWriter, r *http.Request) {

	s, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
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

	currentFlow, err := app.models.Flows.GetCurrentState(flow.SessionID)
	app.handleErrToNotFound(w, r, err)

	if currentFlow.Equals(flow) {
		err = app.writeJSON(w, http.StatusCreated, envelope{"flow": currentFlow}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		return
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
