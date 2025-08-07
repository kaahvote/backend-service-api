package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kaahvote/backend-service-api/internal/data"
	"github.com/kaahvote/backend-service-api/internal/validator"
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

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	// set the body limit to be 1MB
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})

	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil

}

func (app *application) handleErrToNotFound(w http.ResponseWriter, r *http.Request, err error) {

	switch {
	case errors.Is(err, data.ErrRecordNotFound):
		app.notFoundResponse(w, r)
		return
	default:
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (app *application) readDate(qs url.Values, key string) *string {
	value := qs.Get(key)
	if value == "" {
		return nil
	}
	return &value
}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	value := qs.Get(key)

	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, "must be a positive integer")
		return defaultValue
	}

	if i < 0 {
		v.AddError(key, "must be a positive integer")
		return defaultValue
	}

	return i
}
