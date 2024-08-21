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

type responseSoftware struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func toResponseSoftware(software domain.Software) responseSoftware {
	return responseSoftware{
		Name: software.Name(),
		ID:   software.ID(),
	}
}

func (h HttpServer) CreateSoftware(w http.ResponseWriter, r *http.Request) {
	var req SoftwareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.BadRequest("invalid-json", err, w)
		return
	}

	if err := req.Validate(); err != nil {
		server.BadRequest("invalid-request", err, w)
		return
	}

	domainSoftware := domain.NewSoftware(domain.NewSoftwareData{
		Name: req.Name,
	})

	software, err := h.softwareService.CreateSoftware(r.Context(), domainSoftware)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseSoftware(software),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) UpdateSoftware(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	softwareIDValue, ok := vars["software_id"]
	if !ok {
		server.BadRequest("id not provided", domain.ErrRequired, w)
		return
	}

	softwareID, err := strconv.Atoi(softwareIDValue)
	if err != nil {
		server.BadRequest("invalid request", err, w)
		return
	}

	var req SoftwareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.BadRequest("invalid-json", err, w)
		return
	}

	if err := req.Validate(); err != nil {
		server.BadRequest("invalid-request", err, w)
		return
	}

	domainSoftware := domain.NewSoftware(domain.NewSoftwareData{
		ID:   softwareID,
		Name: req.Name,
	})

	software, err := h.softwareService.UpdateSoftware(r.Context(), domainSoftware)
	if err != nil && errors.Is(err, domain.ErrNotFound) {
		server.InternalError(errors.New("software not found"), w)
		return
	}
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseSoftware(software),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) DeleteSoftware(w http.ResponseWriter, r *http.Request) {
	softwareID, err := getIDFromRequest("software_id", w, r)
	if err != nil {
		return
	}

	err = h.softwareService.DeleteSoftware(r.Context(), softwareID)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := simpleResponse{
		Status: StatusOk,
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) GetSoftware(w http.ResponseWriter, r *http.Request) {
	softwareID, err := getIDFromRequest("software_id", w, r)
	if err != nil {
		return
	}

	software, err := h.softwareService.GetSoftware(r.Context(), softwareID)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseSoftware(software),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) GetSoftwares(w http.ResponseWriter, r *http.Request) {
	softwares, err := h.softwareService.GetSoftwares(r.Context())
	if err != nil {
		server.InternalError(err, w)
		return
	}

	out := make([]responseSoftware, 0, len(softwares))

	for _, software := range softwares {
		respSoftware := toResponseSoftware(software)
		out = append(out, respSoftware)
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   out,
	}

	server.RespondOK(resp, w)

}
