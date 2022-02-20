package models

import "time"

type AuthForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserWithCredentials struct {
	User          User      `json:"user"`
	JWT           string    `json:"jwt"`
	JWTExpiration time.Time `json:"jwtExpiration"`
	Permissions   []string  `json:"permissions"`
}
