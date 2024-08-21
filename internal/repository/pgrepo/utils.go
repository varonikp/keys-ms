package pgrepo

import (
	"github.com/varonikp/keys-ms/internal/domain"
	"github.com/varonikp/keys-ms/internal/repository/models"
)

func domainToSoftware(software domain.Software) models.Software {
	return models.Software{
		ID:   software.ID(),
		Name: software.Name(),
	}
}

func softwareToDomain(software models.Software) domain.Software {
	return domain.NewSoftware(domain.NewSoftwareData{
		ID:   software.ID,
		Name: software.Name,
	})
}

func domainToLicense(license domain.License) models.License {
	return models.License{
		ID:         license.ID(),
		SoftwareID: license.SoftwareID(),
		UserID:     license.UserID(),
		CreatedAt:  license.CreatedAt(),
		ExpireAt:   license.ExpireAt(),
	}
}

func licenseToDomain(license models.License) domain.License {
	return domain.NewLicense(domain.NewLicenseData{
		ID:         license.ID,
		SoftwareID: license.SoftwareID,
		UserID:     license.UserID,
		CreatedAt:  license.CreatedAt,
		ExpireAt:   license.ExpireAt,
	})
}

func domainToUser(user domain.User) models.User {
	return models.User{
		ID:           user.ID(),
		Login:        user.Login(),
		Password:     user.Password(),
		HasAdminRole: user.HasAdminRole(),
	}
}

func userToDomain(user models.User) domain.User {
	return domain.NewUser(domain.NewUserData{
		ID:           user.ID,
		Login:        user.Login,
		Password:     user.Password,
		HasAdminRole: user.HasAdminRole,
	})
}
