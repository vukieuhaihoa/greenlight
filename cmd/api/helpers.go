package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

var (
	ErrInvalidIDParameter = fmt.Errorf("invalid id parameter")
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, ErrInvalidIDParameter
	}
	return int64(id), nil
}
