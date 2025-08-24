package ingress

import (
	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
)

type TokenServicePorts interface {
	GenerateToken(roleID constants.Roles, permissions []string, userInfo *models.User) (string, error)
	GetTokenInfo(token string) (*models.Token, error)
	HavePermission(token, permission string) bool
}
