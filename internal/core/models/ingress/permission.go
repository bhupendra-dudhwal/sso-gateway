package ingress

import (
	"regexp"
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Permission struct {
	ID          string           `json:"id" bson:"_id"`
	Description string           `json:"description" bson:"description"`
	Status      constants.Status `json:"status" bson:"status"`
	CreatedBy   int              `json:"created_by,omitempty" bson:"created_by,omitempty"`
	UpdatedBy   int              `json:"updated_by,omitempty" bson:"updated_by"`
	CreatedAt   time.Time        `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at,omitempty" bson:"updated_at"`
}

func (p *Permission) Sanitize(operation constants.Operations, userID int) {
	now := time.Now()
	p.ID = utils.SanitizeLower(p.ID)
	p.Description = utils.Sanitize(p.Description)

	if operation == constants.Create {
		p.Status = constants.StatusActive
		p.CreatedBy = userID
		p.CreatedAt = now
	}

	p.UpdatedBy = userID
	p.UpdatedAt = now
}

func (p Permission) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.ID, validation.Required, validation.Length(3, 20), validation.Match(regexp.MustCompile(`^[a-z_]+$`)).Error("Permission must contain only lowercase letters and underscores")),
		validation.Field(&p.Description, validation.Required, validation.Length(3, 200)),
		validation.Field(&p.Status, validation.Required, validation.In(constants.StatusActive, constants.StatusInactive)),
	)
}
