package auth

import (
	"github.com/brissonwilliam/ihavefriends/backend/config"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/api/auth"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/models"
	"github.com/brissonwilliam/ihavefriends/backend/pkg/storage/user"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	JWT_VALIDITY = time.Hour * 24 * 7
)

type Usecase interface {
	Authenticate(form models.AuthForm) (*models.UserWithCredentials, error)
	GetPublicUsers() ([]string, error)
}

func NewUsecase(userRepo user.UserRepository) Usecase {
	return defaultUsecase{
		userRepo: userRepo,
	}
}

type defaultUsecase struct {
	userRepo user.UserRepository
}

func (u defaultUsecase) Authenticate(form models.AuthForm) (*models.UserWithCredentials, error) {
	user, err := u.userRepo.GetByUsername(form.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, core.NewErrNotFound("user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
	if err != nil {
		// don't let the client know the user exists, return not found on password mismatch
		return nil, core.NewErrNotFound("user")
	}

	permissions, err := u.userRepo.GetUserPermissions(user.Id)
	if err != nil {
		return nil, err
	}

	jwtExpiration := time.Now().Add(JWT_VALIDITY)

	jwt, err := newJWT(*user, permissions, jwtExpiration)
	if err != nil {
		// unexpected error, don't return not found here
		return nil, err
	}

	userWithCreds := models.UserWithCredentials{
		User:          *user,
		JWT:           jwt,
		Permissions:   permissions,
		JWTExpiration: jwtExpiration,
	}

	return &userWithCreds, nil
}

func newJWT(user models.User, permissions []string, jwtExpiration time.Time) (string, error) {
	var err error

	claims := &auth.JWTClaims{
		Id:          user.Id,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtExpiration.Unix(),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSigningKey := config.GetWeb().JwtKey
	token, err := at.SignedString([]byte(jwtSigningKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u defaultUsecase) GetPublicUsers() ([]string, error) {
	return u.userRepo.GetPublicUsernames()
}
