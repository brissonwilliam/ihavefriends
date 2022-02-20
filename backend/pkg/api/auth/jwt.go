package auth

import (
	"github.com/brissonwilliam/ihavefriends/backend/pkg/core/uuid"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const CTX_KEY = "user"

type JWTClaims struct {
	Id          uuid.OrderedUUID `json:"id"`
	Permissions []string         `json:"permissions"`
	jwt.StandardClaims
}

func GetJWTClaimsFromContext(ctx echo.Context) *JWTClaims {
	u := ctx.Get(CTX_KEY)
	if u == nil {
		return nil
	}

	return u.(*jwt.Token).Claims.(*JWTClaims)
}
