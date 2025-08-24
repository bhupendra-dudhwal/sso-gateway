package egress

import (
	"context"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"gorm.io/gorm"
)

type DatabaseConnectionPorts interface {
	Connect() (*gorm.DB, error)
	Close() error
}

type RoleRepositoryPorts interface {
	Add(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id constants.Roles) (*models.Role, error)
}

type UserRepositoryPorts interface {
	Add(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	LockByID(ctx context.Context, id int, lockout_until time.Time) error
}
