package main

import (
	"net/http"

	"github.com/kaahvote/backend-service-api/internal/data"
	"github.com/kaahvote/backend-service-api/internal/validator"
)

func (app *application) postSessionFlowHandler(w http.ResponseWriter, r *http.Request) {

	s, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	var input struct {
		StateID int64  `json:"stateId"`
		Comment string `json:"comment"`
	}

	app.readJSON(w, r, &input)

	flow := &data.Flow{
		SessionID:       s.ID,
		StateID:         input.StateID,
		Comment:         input.Comment,
		SessionPublicID: s.PublicID,
	}

	currentFlow, err := app.models.Flows.GetCurrentState(flow.SessionID)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	currentFlow.SessionPublicID = s.PublicID

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

func (app *application) getSessionFlowHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	qs := r.URL.Query()
	v := validator.New()

	page := app.readInt(qs, "currentPage", 1, v)
	pageSize := app.readInt(qs, "pageSize", 5, v)
	sort := app.readString(qs, "sort", "createdAt")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	filters := data.FlowFilters{
		SessionID: session.ID,
		Filters: data.Filters{
			Page:         page,
			PageSize:     pageSize,
			SortSafeList: []string{"state", "-state", "createdAt", "-createdAt"},
			Sort:         sort,
		},
	}

	flows, metadata, err := app.models.Flows.GetFullHistory(filters)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "flows": flows}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
