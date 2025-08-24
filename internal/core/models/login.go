package models

import (
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type LoginHistory struct {
	ID         int              `json:"id"`
	UserID     int              `json:"user_id"`
	Status     constants.Status `json:"status"`
	Token      string           `json:"token"`
	Reason     string           `json:"reason"`
	Permission []string         `json:"permissions"`
	LoginAt    time.Time        `json:"login_at"`
}

type LoginRequest struct {
	IsUsingPassword bool   `json:"is_using_password"`
	IsUsingMobile   bool   `json:"is_using_mobile"`
	IsUsingEmail    bool   `json:"is_using_email"`
	Email           string `json:"email,omitempty"`
	Password        string `json:"password,omitempty"`
	MobileNumber    int64  `json:"mobile_number,omitempty"`
	DeviceHash      string `json:"device_hash,omitempty"`
}

func (l *LoginRequest) Sanitize() {
	l.Email = utils.SanitizeLower(l.Email)
	l.Password = utils.Sanitize(l.Password)
	l.DeviceHash = utils.Sanitize(l.DeviceHash)
}

func (l LoginRequest) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.When(l.IsUsingEmail, validation.Required, is.Email)),

		// MobileNumber validation using reusable rule
		validation.Field(&l.MobileNumber,
			validation.When(l.IsUsingMobile,
				utils.MobileNumberValidation(l.IsUsingMobile),
			),
		),

		validation.Field(&l.Password,
			validation.When(l.IsUsingPassword,
				validation.Required,
				utils.PasswordStrengthValidation(8, 20),
			),
		),
	)
}
