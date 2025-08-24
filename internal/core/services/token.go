package services

import (
	"fmt"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/ingress"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type tokenService struct {
	config *models.Config
	logger ports.Logger
}

func NewTokenService(config *models.Config, logger ports.Logger) ingress.TokenServicePorts {
	return &tokenService{
		config: config,
		logger: logger,
	}
}

func (tk *tokenService) GenerateToken(roleID constants.Roles, permissions []string, userInfo *models.User) (string, error) {
	var (
		userID int
		now    = time.Now()
	)

	if userInfo != nil {
		userID = userInfo.ID
	}

	expiryAt := now.Add(tk.config.Jwt.LifeSpan)
	claims := models.Token{
		UserID:      userID,
		Role:        roleID,
		Permissions: sliceStringToMapStruct(permissions),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    tk.config.Jwt.Issuer,
			Subject:   tk.config.Jwt.Subject,
			Audience:  tk.config.Jwt.Audience,
			ExpiresAt: jwt.NewNumericDate(expiryAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tk.config.Jwt.SecretKey))
}

// HavePermission checks if the token contains the required permission
func (tk *tokenService) HavePermission(token, permission string) bool {
	tokenInfo, err := tk.GetTokenInfo(token)
	if err != nil || tokenInfo == nil {
		return false
	}

	_, found := tokenInfo.Permissions[permission]
	return found
}

// GetTokenInfo parses the JWT token and extracts the token claims
func (tk *tokenService) GetTokenInfo(token string) (*models.Token, error) {
	var claims models.Token

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(tk.config.Jwt.SecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Additional type assertion to ensure safe usage
	if _, ok := parsedToken.Claims.(*models.Token); !ok {
		return nil, fmt.Errorf("invalid token claims type")
	}

	return &claims, nil
}
