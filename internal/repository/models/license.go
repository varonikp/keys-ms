package models

import (
	"time"
)

type License struct {
	ID         int       `db:"id"`
	SoftwareID int       `db:"software_id"`
	UserID     int       `db:"user_id"`
	CreatedAt  time.Time `db:"created_at"`
	ExpireAt   time.Time `db:"expire_at"`
}
