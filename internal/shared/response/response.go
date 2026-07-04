package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"

	"github.com/0xrinful/reddit-clone/internal/shared/request"
)

type envelope map[string]any

type Responder struct {
	logger *slog.Logger
}

func NewResponder(logger *slog.Logger) *Responder {
	return &Responder{logger}
}

func (r *Responder) JSON(w http.ResponseWriter, status int, data any, headers ...http.Header) {
	js, err := json.Marshal(data)
	if err != nil {
		r.logger.Error("json marshal failed", "err", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error","code":"internal_error"}` + "\n"))
		return
	}

	js = append(js, '\n')

	for _, header := range headers {
		maps.Copy(w.Header(), header)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func (r *Responder) Error(w http.ResponseWriter, status int, code string, message any) {
	r.JSON(w, status, envelope{"error": message, "code": code})
}

func (r *Responder) NotFound(w http.ResponseWriter, rq *http.Request) {
	r.Error(w, http.StatusNotFound, "not_found", "resource not found")
}

func (r *Responder) MethodNotAllowed(w http.ResponseWriter, rq *http.Request) {
	r.Error(w, http.StatusMethodNotAllowed, "method_not_allowed",
		fmt.Sprintf("%s method not allowed", rq.Method))
}

func (r *Responder) ServerError(w http.ResponseWriter, err error) {
	r.logger.Error("internal server error", "err", err)
	r.Error(w, http.StatusInternalServerError, "internal_error", "internal server error")
}

func (r *Responder) DecodeError(w http.ResponseWriter, err error) {
	var decodeErr *request.DecodeError
	if errors.As(err, &decodeErr) {
		r.Error(w, decodeErr.Status, "bad_request", decodeErr.Message)
		return
	}
	r.ServerError(w, err)
}

func (r *Responder) ValidationError(w http.ResponseWriter, errors map[string]string) {
	r.Error(w, http.StatusUnprocessableEntity, "validation_error", errors)
}

func (r *Responder) TooManyRequests(w http.ResponseWriter) {
	r.Error(w, http.StatusTooManyRequests, "too_many_requests", "too many requests")
}

func (r *Responder) InvalidCredentials(w http.ResponseWriter) {
	r.Error(w, http.StatusUnauthorized, "invalid_credentials", "invalid credentials")
}

func (r *Responder) InvalidToken(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	r.Error(w, http.StatusUnauthorized, "invalid_token", "invalid or missing token")
}

func (r *Responder) Unauthorized(w http.ResponseWriter) {
	r.Error(w, http.StatusUnauthorized, "unauthorized", "authentication required")
}

func (r *Responder) Forbidden(w http.ResponseWriter) {
	r.Error(w, http.StatusForbidden, "forbidden", "access denied")
}

func (r *Responder) ForbiddenMsg(w http.ResponseWriter, code, message string) {
	r.Error(w, http.StatusForbidden, code, message)
}

func (r *Responder) Conflict(w http.ResponseWriter) {
	r.Error(w, http.StatusConflict, "conflict", "resource conflict")
}

func (r *Responder) ConflictMsg(w http.ResponseWriter, code, message string) {
	r.Error(w, http.StatusConflict, code, message)
}
