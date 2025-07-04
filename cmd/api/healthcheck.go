package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {

	body := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, body, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered an error, please try again later.", http.StatusInternalServerError)
	}
}
