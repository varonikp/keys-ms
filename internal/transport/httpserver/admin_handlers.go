package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/varonikp/keys-ms/internal/common/server"
	"github.com/varonikp/keys-ms/internal/domain"
)

func (h HttpServer) GrantAdmin(w http.ResponseWriter, r *http.Request) {
	var req GrantAdminRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.BadRequest("invalid-json", err, w)
		return
	}

	if err := req.Validate(); err != nil {
		server.BadRequest("bad request", err, w)
		return
	}

	var targetUser domain.User
	if req.Login != "" {
		user, err := h.userService.GetUser(r.Context(), req.Login)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			server.InternalError(err, w)
			return
		}
		if err != nil && errors.Is(err, domain.ErrNotFound) {
			server.BadRequest("user doesn't exist", domain.ErrNotFound, w)
			return
		}

		targetUser = user
	} else {
		user, err := h.userService.GetUserByID(r.Context(), req.ID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			server.InternalError(err, w)
			return
		}

		if err != nil && errors.Is(err, domain.ErrNotFound) {
			server.BadRequest("user doesn't exist", domain.ErrNotFound, w)
			return
		}

		targetUser = user
	}

	user := setAdminRole(targetUser, true)

	_, err := h.userService.UpdateUser(r.Context(), user)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	server.RespondOK(resp{"status": "ok"}, w)
}

func (h HttpServer) RevokeAdmin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adminTag := vars["tag"]

	var targetUser domain.User
	if adminID, err := strconv.Atoi(adminTag); err == nil { // check if tag is ID
		user, err := h.userService.GetUserByID(r.Context(), adminID)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			server.InternalError(err, w)
			return
		}
		if err != nil && errors.Is(err, domain.ErrNotFound) {
			server.BadRequest("user doesn't exist", domain.ErrNotFound, w)
			return
		}

		targetUser = user
	} else { // tag is login
		user, err := h.userService.GetUser(r.Context(), adminTag)
		if err != nil && !errors.Is(err, domain.ErrNotFound) {
			server.InternalError(err, w)
			return
		}
		if err != nil && errors.Is(err, domain.ErrNotFound) {
			server.BadRequest("user doesn't exist", domain.ErrNotFound, w)
			return
		}

		targetUser = user
	}

	user := setAdminRole(targetUser, false)

	_, err := h.userService.UpdateUser(r.Context(), user)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	server.RespondOK(resp{"status": "ok"}, w)
}
