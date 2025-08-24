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
		return nil, utils.ErrDocumentNotFound
	}
	return &role, err
}

func (r *role) GetByIDs(ctx context.Context, ids []constants.Roles) ([]models.Role, error) {
	var role []models.Role
	err := r.client.WithContext(ctx).Where("id IN(?)", ids).Find(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrDocumentNotFound
	}
	return role, err
}

func (r *role) DeleteByID(ctx context.Context, id constants.Roles) error {
	var role models.Role
	err := r.client.WithContext(ctx).Delete(&role, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrDocumentNotFound
	}
	return err
}

func (r *role) GetRolesWithoutPagination(ctx context.Context) ([]models.Role, error) {
	var role []models.Role
	err := r.client.WithContext(ctx).Find(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrDocumentNotFound
	}
	return role, err
}
