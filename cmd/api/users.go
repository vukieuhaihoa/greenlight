package main

import (
	"errors"
	"net/http"

	"github.com/vukieuhaihoa/greenlight/internal/data"
	"github.com/vukieuhaihoa/greenlight/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var inputData input

	err := app.readJSON(w, r, &inputData)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      inputData.Name,
		Email:     inputData.Email,
		Activated: false,
	}

	err = user.Password.Set(inputData.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.background(func() {
		err := app.mailer.Send(user.Email, "user_welcome.tmpl", user)
		if err != nil {
			app.logger.Error("unable to send welcome email", "error", err)
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
