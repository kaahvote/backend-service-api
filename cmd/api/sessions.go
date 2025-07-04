package main

import "net/http"

func (app *application) getSessionHandler(w http.ResponseWriter, r *http.Request) {

	// empty slice
	sessions := make([]any, 0)

	err := app.writeJSON(w, http.StatusOK, envelope{"sessions": sessions}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
