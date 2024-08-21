package httpserver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/varonikp/keys-ms/internal/common/server"
	"github.com/varonikp/keys-ms/internal/domain"
)

type responseLicense struct {
	ID         int       `json:"id"`
	SoftwareID int       `json:"software_id"`
	UserID     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	ExpireAt   time.Time `json:"expire_at"`
}

func toResponseLicense(license domain.License) responseLicense {
	return responseLicense{
		ID:         license.ID(),
		SoftwareID: license.SoftwareID(),
		UserID:     license.UserID(),
		CreatedAt:  license.CreatedAt(),
		ExpireAt:   license.ExpireAt(),
	}
}

func (h HttpServer) GetLicenses(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest("user_id", w, r)
	if err != nil {
		return
	}

	user := r.Context().Value(ContextUserKey).(domain.User)
	if user.ID() != userID && !user.HasAdminRole() {
		server.Unauthorized(w)
		return
	}

	licenses, err := h.licenseService.GetLicensesByUserID(r.Context(), userID)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	out := make([]responseLicense, 0, len(licenses))
	for _, license := range licenses {
		respLicense := toResponseLicense(license)
		out = append(out, respLicense)
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   out,
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) GetLicense(w http.ResponseWriter, r *http.Request) {
	licenseID, err := getIDFromRequest("license_id", w, r)
	if err != nil {
		return
	}

	license, err := h.licenseService.GetLicense(r.Context(), licenseID)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseLicense(license),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) CreateLicense(w http.ResponseWriter, r *http.Request) {
	var request LicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.BadRequest("invalid json", err, w)
		return
	}
	if err := request.Validate(); err != nil {
		server.BadRequest("invalid data", err, w)
		return
	}

	domainLicense := domain.NewLicense(domain.NewLicenseData{
		SoftwareID: request.SoftwareID,
		UserID:     request.UserID,
		CreatedAt:  time.Now(),
		ExpireAt:   time.Unix(request.ExpireAt, 0),
	})

	license, err := h.licenseService.CreateLicense(r.Context(), domainLicense)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseLicense(license),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) UpdateLicense(w http.ResponseWriter, r *http.Request) {
	licenseID, err := getIDFromRequest("license_id", w, r)
	if err != nil {
		return
	}

	var request LicenseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		server.BadRequest("invalid json", err, w)
		return
	}

	license, err := h.licenseService.GetLicense(r.Context(), licenseID)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	licenseData := domain.NewLicenseData{
		ID:         licenseID,
		SoftwareID: license.SoftwareID(),
		UserID:     license.UserID(),
		CreatedAt:  license.CreatedAt(),
		ExpireAt:   license.ExpireAt(),
	}

	if request.ExpireAt != 0 {
		licenseData.ExpireAt = time.Unix(request.ExpireAt, 0)
	}
	if request.SoftwareID != 0 {
		licenseData.SoftwareID = request.SoftwareID
	}
	if request.UserID != 0 {
		licenseData.UserID = request.UserID
	}

	updatedLicense := domain.NewLicense(licenseData)

	_, err = h.licenseService.UpdateLicense(r.Context(), updatedLicense)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := responseWithData{
		Status: StatusOk,
		Data:   toResponseLicense(updatedLicense),
	}

	server.RespondOK(resp, w)
}

func (h HttpServer) DeleteLicense(w http.ResponseWriter, r *http.Request) {
	licenseID, err := getIDFromRequest("license_id", w, r)
	if err != nil {
		return
	}

	err = h.licenseService.DeleteLicense(r.Context(), licenseID)
	if err != nil {
		server.InternalError(err, w)
		return
	}

	resp := simpleResponse{
		Status: StatusOk,
	}

	server.RespondOK(resp, w)
}
