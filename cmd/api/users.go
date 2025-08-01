package main

import (
	"net/http"
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

	sessions, err := app.models.Sessions.ListSessionsByUserID(user.ID)
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
