package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DecodeError struct {
	Status  int
	Message string
}

func (e *DecodeError) Error() string {
	return e.Message
}

func decodeErr(status int, msg string) *DecodeError {
	return &DecodeError{Status: status, Message: msg}
}

func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
		if mediaType != "application/json" {
			return decodeErr(
				http.StatusUnsupportedMediaType,
				"Content-Type header is not application/json",
			)
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		return mapDecodeError(err)
	}

	if dec.More() {
		return decodeErr(
			http.StatusBadRequest,
			"request body must only contain a single JSON object",
		)
	}

	return nil
}

func mapDecodeError(err error) error {
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError
	var maxBytesError *http.MaxBytesError

	switch {
	case errors.As(err, &syntaxError):
		return decodeErr(
			http.StatusBadRequest,
			fmt.Sprintf(
				"Request body contains badly-formed JSON (at position %d)",
				syntaxError.Offset,
			),
		)

	case errors.Is(err, io.ErrUnexpectedEOF):
		return decodeErr(http.StatusBadRequest, "request body contains badly-formed JSON")

	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return decodeErr(
				http.StatusBadRequest,
				fmt.Sprintf("incorrect type for field `%s`, expected `%s`",
					unmarshalTypeError.Field, unmarshalTypeError.Type),
			)
		}
		return decodeErr(
			http.StatusBadRequest,
			fmt.Sprintf("incorrect JSON type at character %d", unmarshalTypeError.Offset),
		)

	case errors.Is(err, io.EOF):
		return decodeErr(http.StatusBadRequest, "request body must not be empty")

	case strings.HasPrefix(err.Error(), "json: unknown field"):
		field := strings.TrimPrefix(err.Error(), "json: unknown field ")
		return decodeErr(http.StatusBadRequest, fmt.Sprintf("unknown field `%s`", field))

	case errors.As(err, &maxBytesError):
		return decodeErr(
			http.StatusRequestEntityTooLarge,
			fmt.Sprintf("request body must not be larger than %d bytes", maxBytesError.Limit),
		)

	case errors.As(err, &invalidUnmarshalError):
		panic(err)

	default:
		return err
	}
}
