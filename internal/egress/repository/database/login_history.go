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

type loginHistory struct {
	client *gorm.DB
}

func NewloginHistoryRepository(client *gorm.DB) egress.LoginHistoryPorts {
	return &loginHistory{
		client: client,
	}
}

func (l *loginHistory) Add(ctx context.Context, loginHistory *models.LoginHistory) error {
	err := l.client.WithContext(ctx).Create(loginHistory).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return utils.ErrDocumentNotFound
	}
	return err
}

func (l *loginHistory) GetByID(ctx context.Context, id int) (*models.LoginHistory, error) {
	var loginHistory models.LoginHistory
	err := l.client.WithContext(ctx).First(&loginHistory, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrDocumentNotFound
	}
	return &loginHistory, err
}

func (l *loginHistory) GetByIDAndLoginAt(ctx context.Context, id int, loginAt time.Time) ([]models.LoginHistory, error) {
	var loginHistory []models.LoginHistory
	err := l.client.WithContext(ctx).First(&loginHistory, id).Error
	return loginHistory, err
}
