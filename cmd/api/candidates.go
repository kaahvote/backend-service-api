package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kaahvote/backend-service-api/internal/data"
	"github.com/kaahvote/backend-service-api/internal/validator"
)

func (app *application) postSessionCandidatesHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name   string `json:"name"`
		UserID int64  `json:"userId"`
	}

	if input.UserID == 0 {
		// TODO: Create an anonymous user here
		input.UserID = 2
	}

	session, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	c := &data.Candidate{
		Name:      input.Name,
		UserID:    input.UserID,
		SessionID: session.ID,
		CreatedAt: time.Now(),
	}

	// TODO: add validation for candidate fields

	err = app.models.Candidate.Insert(c)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/sessions/%s/candidates/%d", session.PublicID, c.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"candidate": c}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getSessionCandidatesHandler(w http.ResponseWriter, r *http.Request) {

	session, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	qs := r.URL.Query()

	name := app.readString(qs, "name", "")

	v := validator.New()
	page := app.readInt(qs, "currentPage", 1, v)
	pageSize := app.readInt(qs, "pageSize", 5, v)
	sort := app.readString(qs, "sort", "createdAt")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	filters := &data.CandidateFilters{
		Name:      name,
		SessionID: session.ID,
		Filters: data.Filters{
			Page:         page,
			PageSize:     pageSize,
			Sort:         sort,
			SortSafeList: []string{"name", "-name", "createdAt", "-createdAt"},
		},
	}

	//TODO: validate the filters, sort safe list

	candidates, metadata, err := app.models.Candidate.ListFiltering(filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"candidates": candidates, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteSessionCandidateHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getSingleCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	_, err := app.getSession(r)
	if err != nil {
		app.handleErrToNotFound(w, r, err)
		return
	}

	candidatedID, err := strconv.Atoi(app.readStringParam(r, "candidate_id"))

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if candidatedID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	candidate, err := app.models.Candidate.Get(int64(candidatedID))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"candidate": candidate}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
