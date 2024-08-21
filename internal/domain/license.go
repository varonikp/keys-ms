package domain

import (
	"time"
)

type License struct {
	id         int
	softwareID int
	userID     int
	createdAt  time.Time
	expireAt   time.Time
}

type NewLicenseData struct {
	ID         int
	SoftwareID int
	UserID     int
	CreatedAt  time.Time
	ExpireAt   time.Time
}

func NewLicense(data NewLicenseData) License {
	return License{
		id:         data.ID,
		softwareID: data.SoftwareID,
		userID:     data.UserID,
		createdAt:  data.CreatedAt,
		expireAt:   data.ExpireAt,
	}
}

func (l License) ID() int {
	return l.id
}

func (l License) SoftwareID() int {
	return l.softwareID
}

func (l License) UserID() int {
	return l.userID
}

func (l License) CreatedAt() time.Time {
	return l.createdAt
}

func (l License) ExpireAt() time.Time {
	return l.expireAt
}
