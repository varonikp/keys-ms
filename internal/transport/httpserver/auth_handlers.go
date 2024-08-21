package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/varonikp/keys-ms/internal/common/server"
	"github.com/varonikp/keys-ms/internal/domain"
)

func (h HttpServer) SignUp(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		server.BadRequest("invalid-request", err, w)
		return
	}

	if err := authRequest.Validate(); err != nil {
		server.BadRequest("invalid-json", err, w)
		return
	}

	userExists, err := h.userService.IsUserExists(r.Context(), authRequest.Login)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	if userExists {
		server.BadRequest("user already exists", errors.New("user with this login already exists"), w)
		return
	}

	encryptedPassword, err := encryptPassword(authRequest.Password)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	user := domain.NewUser(domain.NewUserData{
		Login:    authRequest.Login,
		Password: encryptedPassword,
	})

	_, err = h.userService.CreateUser(r.Context(), user)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	server.RespondOK(
		resp{"success": true}, w,
	)
}

func (h HttpServer) SignIn(w http.ResponseWriter, r *http.Request) {
	var authRequest AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&authRequest); err != nil {
		server.BadRequest("invalid-request", err, w)
		return
	}

	if err := authRequest.Validate(); err != nil {
		server.BadRequest("invalid-json", err, w)
		return
	}

	user, err := h.userService.GetUser(r.Context(), authRequest.Login)
	if err != nil {
		server.RespondWithError(err, w)
		return
	}

	if !isEqualPassword(authRequest.Password, user.Password()) {
		server.BadRequest("invalid password", errors.New("authorization failed"), w)
		return
	}

	token, err := h.tokenService.GenerateToken(user)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	server.RespondOK(resp{
		"token": token,
	}, w)
}
