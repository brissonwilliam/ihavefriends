package models

type AuthForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserWithCredentials struct {
	User        User     `json:"user"`
	JWT         string   `json:"jwt"`
	Permissions []string `json:"permissions"`
}
