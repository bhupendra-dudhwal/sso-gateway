package egress

import (
	"context"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	ingressModel "github.com/bhupendra-dudhwal/sso-gateway/internal/core/models/ingress"
	"gorm.io/gorm"
)

type DatabaseConnectionPorts interface {
	Connect() (*gorm.DB, error)
	Close() error
}

type RoleRepositoryPorts interface {
	Add(ctx context.Context, role *models.Role) error
	GetByID(ctx context.Context, id constants.Roles) (*models.Role, error)
	DeleteByID(ctx context.Context, id constants.Roles) error
	GetByIDs(ctx context.Context, ids []constants.Roles) ([]models.Role, error)
	GetRolesWithoutPagination(ctx context.Context) ([]models.Role, error)
}

type UserRepositoryPorts interface {
	Add(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	LockByID(ctx context.Context, id int, lockout_until time.Time) error
}

type PermissionRepositoryPorts interface {
	Add(ctx context.Context, role *ingressModel.Permission) error
	GetByID(ctx context.Context, id string) (*ingressModel.Permission, error)
	DeleteByID(ctx context.Context, id string) error
	GetByIDs(ctx context.Context, ids []string) ([]ingressModel.Permission, error)
	GetPermissionWithoutPagination(ctx context.Context) ([]ingressModel.Permission, error)
}
