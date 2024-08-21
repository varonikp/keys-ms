package httpserver

import (
	"fmt"

	"github.com/varonikp/keys-ms/internal/domain"
)

type Status string

var (
	StatusOk Status = "ok"
)

type responseWithData struct {
	Status Status `json:"status"`
	Data   any    `json:"data"`
}

type simpleResponse struct {
	Status Status `json:"status"`
}

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *AuthRequest) Validate() error {
	if r.Login == "" {
		return fmt.Errorf("%w: username", domain.ErrRequired)
	}

	if r.Password == "" {
		return fmt.Errorf("%w: password", domain.ErrRequired)
	}

	return nil
}

type GrantAdminRequest struct {
	Login string `json:"login,omitempty"`
	ID    int    `json:"id,omitempty"`
}

func (r *GrantAdminRequest) Validate() error {
	if r.Login == "" && r.ID == 0 {
		return fmt.Errorf("%w: %s", domain.ErrRequired, "`login` or `id`")
	}

	return nil
}

type SoftwareRequest struct {
	Name string `json:"name"`
}

func (r *SoftwareRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: %s", domain.ErrRequired, "name")
	}

	return nil
}

func setAdminRole(user domain.User, adminRole bool) domain.User {
	return domain.NewUser(domain.NewUserData{
		ID:           user.ID(),
		Login:        user.Login(),
		Password:     user.Password(),
		HasAdminRole: adminRole,
	})
}

type LicenseRequest struct {
	SoftwareID int   `json:"software_id"`
	UserID     int   `json:"user_id"`
	ExpireAt   int64 `json:"expire_at"`
}

func (r *LicenseRequest) Validate() error {
	if r.SoftwareID == 0 {
		return fmt.Errorf("%w: %s", domain.ErrRequired, "software_id")
	}
	if r.UserID == 0 {
		return fmt.Errorf("%w: %s", domain.ErrRequired, "user_id")
	}
	if r.ExpireAt == 0 {
		return fmt.Errorf("%w: %s", domain.ErrRequired, "expire_at")
	}

	return nil
}

type UserRequest struct {
	Login        string `json:"login,omitempty"`
	Password     string `json:"password,omitempty"`
	HasAdminRole *bool  `json:"has_admin_role,omitempty"`
}
