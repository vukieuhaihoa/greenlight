package main

import (
	"fmt"
	"net/http"
)

// logError writes an error message to the errorLog, along with some additional context about the request.
//
// Parameters:
// - r: The HTTP request that caused the error.
// - err: The error to be logged.
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// errorResponse is a helper method for sending JSON-formatted error messages to the client.
//
// Parameters:
// - w: The http.ResponseWriter used to send the response.
// - r: The HTTP request that caused the error.
// - status: The HTTP status code to be sent in the response.
// - message: The error message to be included in the response body.
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// serverErrorResponse is a helper method for sending a 500 Internal Server Error response to the client, along with logging the error.
//
// Parameters:
// - w: The http.ResponseWriter used to send the response.
// - r: The HTTP request that caused the error.
// - err: The error to be logged and included in the response message.
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}
