package user

import (
	"errors"
	"fmt"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/db"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage"
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
			deleted_ut,
			is_public
		FROM user
		WHERE %s
	`

	queryGetPublicUsernames = `
		SELECT username FROM user WHERE is_public = 1
	`

	queryGetUserPermissions = `
		SELECT permission FROM user_permission WHERE user_id = ?
	`

	queryCreateUser = `
		INSERT INTO user (id, username, password, email, is_public) VALUES (:id, :username, :password, :email, :is_public)
	`
)

var (
	queryGetUserByUsername = fmt.Sprintf(baseQueryGetUser, "username = ?")
	queryGetUserById       = fmt.Sprintf(baseQueryGetUser, "id = ?")
)

type UserRepository interface {
	GetById(id uuid.OrderedUUID) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetPublicUsernames() ([]string, error)
	GetUserPermissions(userId uuid.OrderedUUID) ([]string, error)

	Create(user models.CreaterUserForm) error

	WithUnitOfWork(uow storage.UnitOfWork) UserRepository
}

type defaultUserRepository struct {
	db sqlx.Ext
}

func NewUserRepository(db sqlx.Ext) UserRepository {
	return defaultUserRepository{
		db: db,
	}
}

func (r defaultUserRepository) WithUnitOfWork(uow storage.UnitOfWork) UserRepository {
	tx := storage.UnitAsTransaction(uow)
	return NewUserRepository(tx)
}

func (r defaultUserRepository) GetByUsername(username string) (*models.User, error) {
	users := []models.User{}
	err := sqlx.Select(r.db, &users, queryGetUserByUsername, username)
	if err != nil {
		return nil, err
	}

	if len(users) > 1 {
		return nil, errors.New("Got more than one user. Consider using a unique key")
	}

	if len(users) < 1 {
		return nil, core.NewErrNotFound("user")
	}

	return &users[0], nil
}

func (r defaultUserRepository) GetById(id uuid.OrderedUUID) (*models.User, error) {
	users := []models.User{}
	err := sqlx.Select(r.db, &users, queryGetUserById, id)
	if err != nil {
		return nil, err
	}

	if len(users) < 1 {
		return nil, core.NewErrNotFound("user")
	}

	return &users[0], nil
}

func (r defaultUserRepository) GetPublicUsernames() ([]string, error) {
	usernames := []string{}
	err := sqlx.Select(r.db, &usernames, queryGetPublicUsernames)
	if err != nil {
		return nil, err
	}
	return usernames, nil
}

func (r defaultUserRepository) GetUserPermissions(userId uuid.OrderedUUID) ([]string, error) {
	permissions := []string{}
	err := sqlx.Select(r.db, &permissions, queryGetUserPermissions, userId)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

func (r defaultUserRepository) Create(user models.CreaterUserForm) error {
	_, err := sqlx.NamedExec(r.db, queryCreateUser, &user)

	if db.CollidesWithUniqueIndex(err) {
		return core.NewErrConflict("username or email")
	}

	return err
}
