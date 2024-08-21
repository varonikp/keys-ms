package domain

type User struct {
	id           int
	login        string
	password     string
	hasAdminRole bool
}

type NewUserData struct {
	ID           int
	Login        string
	Password     string
	HasAdminRole bool
}

func NewUser(data NewUserData) User {
	return User{
		id:           data.ID,
		login:        data.Login,
		password:     data.Password,
		hasAdminRole: data.HasAdminRole,
	}
}

func (u User) ID() int {
	return u.id
}

func (u User) Login() string {
	return u.login
}

func (u User) Password() string {
	return u.password
}

func (u User) HasAdminRole() bool {
	return u.hasAdminRole
}
