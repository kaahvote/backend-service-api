package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type envelope map[string]any

func (app *application) readStringParam(r *http.Request, param string) string {
	params := httprouter.ParamsFromContext(r.Context())
	return params.ByName(param)
}

func (app *application) writeJSON(w http.ResponseWriter, status int, body envelope, headers http.Header) error {

	resp, err := json.MarshalIndent(body, "", "\t")
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(resp)

	return nil
}
