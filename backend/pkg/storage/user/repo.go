package user

import (
	"errors"
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/jmoiron/sqlx"
)

const (
	baseQueryGetUser = `
		SELECT
			id,
			username,
			password,
			email,
			created,
			last_updated,
			deleted_ut
		FROM user
		WHERE %s
	`
)

var (
	queryGetUserByUsername = fmt.Sprintf(baseQueryGetUser, "username = ?")
)

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
}

type defaultUserRepository struct {
	db sqlx.Ext
}

func NewUserRepository(db sqlx.Ext) UserRepository {
	return defaultUserRepository{
		db: db,
	}
}

func (r defaultUserRepository) GetByUsername(username string) (*models.User, error) {
	var users []models.User
	err := sqlx.Select(r.db, &users, queryGetUserByUsername, username)
	if err != nil {
		return nil, err
	}

	if len(users) > 1 {
		return nil, errors.New("Got more than one user to authenticate. Consider using a unique key")
	}

	if len(users) < 1 {
		return nil, core.NewErrNotFound("user")
	}

	return &users[0], nil
}
