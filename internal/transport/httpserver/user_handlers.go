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

type responseUser struct {
	ID           int    `json:"id"`
	Login        string `json:"login"`
	HasAdminRole bool   `json:"has_admin_role"`
}

func toResponseUser(user domain.User) responseUser {
	return responseUser{
		ID:           user.ID(),
		Login:        user.Login(),
		HasAdminRole: user.HasAdminRole(),
	}
}

func (h HttpServer) GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserByIdOrTag(w, r)
	if err != nil {
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseUser(user),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.GetUsers(r.Context())
	if err != nil {
		server.InternalError(err, w)
		return
	}

	out := make([]responseUser, 0, len(users))
	for _, user := range users {
		responseUser := toResponseUser(user)

		out = append(out, responseUser)
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   out,
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var request UserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.BadRequest("invalid json", err, w)
		return
	}

	user, err := h.getUserByIdOrTag(w, r)
	if err != nil {
		return
	}

	userData := domain.NewUserData{
		ID:           user.ID(),
		Login:        user.Login(),
		Password:     user.Password(),
		HasAdminRole: user.HasAdminRole(),
	}

	if request.Login != "" {
		userData.Login = request.Login
	}

	if request.Password != "" {
		encryptedPassword, err := encryptPassword(request.Password)
		if err != nil {
			server.InternalError(err, w)
			return
		}

		userData.Password = encryptedPassword
	}

	if request.HasAdminRole != nil {
		userData.HasAdminRole = *request.HasAdminRole
	}

	newUser := domain.NewUser(userData)

	_, err = h.userService.UpdateUser(r.Context(), newUser)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseUser(newUser),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserByIdOrTag(w, r)
	if err != nil {
		return
	}

	err = h.userService.DeleteUser(r.Context(), user.ID())
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := simpleResponse{
		Status: StatusOk,
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) getUserByIdOrTag(w http.ResponseWriter, r *http.Request) (domain.User, error) {
	tag := mux.Vars(r)["tag"]

	var (
		user domain.User
		err  error
	)

	if userID, err := strconv.Atoi(tag); err == nil { // tag is userID
		user, err = h.userService.GetUserByID(r.Context(), userID)
	} else { // tag is login
		user, err = h.userService.GetUser(r.Context(), tag)
	}

	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		server.InternalError(err, w)
		return domain.User{}, err
	}

	if err != nil && errors.Is(err, domain.ErrNotFound) {
		server.BadRequest("user doesn't exist", domain.ErrNotFound, w)
		return domain.User{}, nil
	}

	return user, nil

}
