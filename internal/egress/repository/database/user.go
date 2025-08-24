package database

import (
	"context"
	"errors"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports/egress"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	"gorm.io/gorm"
)

type user struct {
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB) egress.UserRepositoryPorts {
	return &user{
		client: client,
	}
}

func (r *user) Add(ctx context.Context, user *models.User) error {
	err := r.client.WithContext(ctx).Create(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrDocumentNotFound
	}
	return err
}

func (r *user) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := r.client.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *user) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.client.WithContext(ctx).First(&user, email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *user) LockByID(ctx context.Context, id int, lockout_until time.Time) error {
	return r.client.WithContext(ctx).Where("id=?", id).Update("lockout_until", lockout_until).Error
}
