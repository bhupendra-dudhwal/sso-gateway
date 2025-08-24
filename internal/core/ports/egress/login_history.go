package egress

import (
	"context"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
)

type LoginHistoryPorts interface {
	Add(ctx context.Context, loginHistory *models.LoginHistory) error
	GetByID(ctx context.Context, id int) (*models.LoginHistory, error)
	GetByIDAndLoginAt(ctx context.Context, id int, loginAt time.Time) ([]models.LoginHistory, error)
}
