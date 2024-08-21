package domain

import "errors"

var (
	ErrRequired = errors.New("this field is required")
	ErrNotFound = errors.New("not found")
)
