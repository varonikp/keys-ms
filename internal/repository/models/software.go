package models

type Software struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
