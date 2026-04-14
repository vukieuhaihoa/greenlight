package data

import (
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidRuntimeFormat = fmt.Errorf("invalid runtime format")

type Runtime int

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d mins", r)

	qoutedJSONValue := strconv.Quote(jsonValue)

	return []byte(qoutedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquoteJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquoteJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrInvalidRuntimeFormat
	}

	minutes, err := strconv.Atoi(parts[0])
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(minutes)

	return nil
}
