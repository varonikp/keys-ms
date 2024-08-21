package models

type User struct {
	ID           int    `db:"id"`
	Login        string `db:"login"`
	Password     string `db:"password"`
	HasAdminRole bool   `db:"has_admin_role"`
}
