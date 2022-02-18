package models

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"time"
)

type User struct {
	Id          uuid.OrderedUUID `json:"id" db:"id"`
	Username    string           `json:"username" db:"username"`
	Password    string           `json:"-" db:"password"`
	Email       string           `json:"email" db:"email"`
	Created     time.Time        `json:"-" db:"created"`
	LastUpdated time.Time        `json:"-" db:"last_updated"`
	DeletedUnix uint64           `json:"-" db:"deleted_ut"`
	IsPublic    bool             `json:"-" db:"is_public"`
}
