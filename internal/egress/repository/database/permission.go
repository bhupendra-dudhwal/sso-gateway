package database

import (
	"context"
	"errors"

	ingressModel "github.com/bhupendra-dudhwal/sso-gateway/internal/core/models/ingress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"gorm.io/gorm"
)

type permission struct {
	client *gorm.DB
}

func NewPermissionRepository(client *gorm.DB) egress.PermissionRepositoryPorts {
	return &permission{
		client: client,
	}
}

func (r *permission) Add(ctx context.Context, permission *ingressModel.Permission) error {
	err := r.client.WithContext(ctx).Create(permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrDocumentNotFound
	}
	return err
}

func (r *permission) GetByID(ctx context.Context, id string) (*ingressModel.Permission, error) {
	var permission ingressModel.Permission
	err := r.client.WithContext(ctx).First(&permission, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrDocumentNotFound
	}
	return &permission, err
}

func (r *permission) GetByIDs(ctx context.Context, ids []string) ([]ingressModel.Permission, error) {
	var permission []ingressModel.Permission
	err := r.client.WithContext(ctx).Where("id IN (?)", ids).Find(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrDocumentNotFound
	}
	return permission, err
}

func (r *permission) GetPermissionWithoutPagination(ctx context.Context) ([]ingressModel.Permission, error) {
	var permission []ingressModel.Permission
	err := r.client.WithContext(ctx).Find(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrDocumentNotFound
	}
	return permission, err
}

func (r *permission) DeleteByID(ctx context.Context, id string) error {
	var permission ingressModel.Permission
	err := r.client.WithContext(ctx).Delete(&permission, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrDocumentNotFound
	}
	return err
}
