package models

import (
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Role struct {
	ID          constants.Roles  `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Description string           `json:"description" gorm:"type:varchar(255)"`
	Permissions []string         `json:"permissions" gorm:"type:text[]"`
	Status      constants.Status `json:"status" gorm:"type:varchar(20);not null"`
	CreatedBy   int              `json:"created_by,omitempty"`
	UpdatedBy   int              `json:"updated_by,omitempty"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty"`
}

func (r *Role) Sanitize(operation constants.Operations, userID int) {
	now := time.Now()
	r.Description = utils.Sanitize(r.Description)
	r.Permissions = utils.SanitizeLowerSlice(r.Permissions)

	if operation == constants.Create {
		r.Status = constants.StatusActive
		r.CreatedBy = userID
		r.CreatedAt = now
	}

	r.UpdatedBy = userID
	r.UpdatedAt = now
}

func (r Role) Validate() error {
	return validation.ValidateStruct(&r,
		// validation.Field(&r.ID, validation.Required, validation.Length(3, 20), validation.Match(regexp.MustCompile(`^[a-z_]+$`)).Error("ID must contain only lowercase letters and underscores")),
		validation.Field(&r.ID, validation.Required, validation.In(constants.RoleAnonymousUser, constants.RoleSessionUser, constants.RoleSystemAdmin, constants.Roleadmin, constants.Roleuser)),
		validation.Field(&r.Permissions, validation.Required, validation.Each(validation.Length(3, 30))),
		validation.Field(&r.Description, validation.Required, validation.Length(3, 200)),
		validation.Field(&r.Status, validation.Required, validation.In(constants.StatusActive, constants.StatusInactive)),
	)
}

// We can use this if we need to allow all type of status
// func validateStatus(value interface{}) error {
// 	if s, ok := value.(constants.Status); ok && s.IsValid() {
// 		return nil
// 	}
// 	return errors.New("invalid status")
// }

// We can use this if we need to allow all type of roles
// func validateRoles(value interface{}) error {
// 	if s, ok := value.(constants.Roles); ok && s.IsValid() {
// 		return nil
// 	}
// 	return errors.New("invalid status")
// }
