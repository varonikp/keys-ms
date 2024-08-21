package httpserver

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/varonikp/keys-ms/internal/common/server"
)

const (
	authHeader   = "Authorization"
	bearerPrefix = "Bearer "
)

func (h HttpServer) CheckAuthorized(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get(authHeader)

		jwtToken := strings.TrimPrefix(header, bearerPrefix)
		if jwtToken == "" {
			server.Unauthorized(w)
			return
		}

		user, err := h.tokenService.GetUser(jwtToken)
		if err != nil {
			server.InternalError(err, w)
			return
		}
		if user.Login() == "" {
			server.BadRequest("invalid-token", err, w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next(w, r.WithContext(ctx))
	}
}

func (h HttpServer) CheckAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authHeader)

		jwtToken := strings.TrimPrefix(authorizationHeader, bearerPrefix)
		if jwtToken == "" {
			server.Unauthorized(w)
			return
		}

		user, err := h.tokenService.GetUser(jwtToken)
		if err != nil {
			server.InternalError(err, w)
			return
		}
		if user.Login() == "" {
			server.BadRequest("invalid-token", err, w)
			return
		}

		if !user.HasAdminRole() {
			server.InternalError(errors.New("access denied"), w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, user)
		next(w, r.WithContext(ctx))
	}
}
