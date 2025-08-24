package models

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/golang-jwt/jwt/v4"
)

type Token struct {
	UserID      int                 `json:"user_id"`
	Role        constants.Roles     `json:"role"`
	Permissions map[string]struct{} `json:"permission"`
	jwt.RegisteredClaims
}
