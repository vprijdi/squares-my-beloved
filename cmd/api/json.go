package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "appllication/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

// readJSON parses JSON request body into the data stuct
func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func readJSONArray(w http.ResponseWriter, r *http.Request, slicePtr any) error {
	if reflect.TypeOf(slicePtr).Kind() != reflect.Ptr ||
		reflect.TypeOf(slicePtr).Elem().Kind() != reflect.Slice {
		return fmt.Errorf("destination must be a slice pointer")
	}
	return readJSON(w, r, slicePtr)
}

func writeJSONError(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}
	return writeJSON(w, status, &envelope{Error: message})
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &envelope{Data: data})
}
