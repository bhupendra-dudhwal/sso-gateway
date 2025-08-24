package database

import (
	"context"
	"errors"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"gorm.io/gorm"
)

type role struct {
	client *gorm.DB
}

func NewRoleRepository(client *gorm.DB) egress.RoleRepositoryPorts {
	return &role{
		client: client,
	}
}

func (r *role) Add(ctx context.Context, role *models.Role) error {
	err := r.client.WithContext(ctx).Create(role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrDocumentNotFound
	}
	return err
}

func (r *role) GetByID(ctx context.Context, id constants.Roles) (*models.Role, error) {
	var role models.Role
	err := r.client.WithContext(ctx).First(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}
