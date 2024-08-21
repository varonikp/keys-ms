package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/varonikp/keys-ms/internal/domain"
)

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Error      string `json:"error,omitempty"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}

func BadRequest(slug string, err error, w http.ResponseWriter) {
	httpRespondWithError(err, slug, w, "bad request", http.StatusBadRequest)
}

func Unauthorized(w http.ResponseWriter) {
	httpSimpleErrorRespond("unathorized", http.StatusUnauthorized, w)
}

func InternalError(err error, w http.ResponseWriter) {
	httpRespondWithError(err, "internal-server-error", w, "internal error", http.StatusInternalServerError)
}

func RespondWithError(err error, w http.ResponseWriter) {
	if errors.Is(err, domain.ErrNotFound) {
		httpRespondWithError(err, "user not found", w, "user not found", http.StatusNotFound)
		return
	}

	InternalError(err, w)
}

func httpSimpleErrorRespond(slug string, status int, w http.ResponseWriter) {
	response := ErrorResponse{
		Slug:       slug,
		httpStatus: status,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, msg string, status int) {
	slog.Error(
		"http error",
		slog.String("error", err.Error()), slog.String("slug", slug), slog.String("msg", msg),
	)

	response := ErrorResponse{
		Slug:       slug,
		httpStatus: status,
	}

	if os.Getenv("DEBUG_ERRORS") != "" && err != nil {
		response.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
