package models

import (
	"time"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/utils"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type User struct {
	ID           int              `json:"id"`
	Mobile       int              `json:"mobile"`
	Role         constants.Roles  `json:"role"`
	Permissions  []string         `json:"permissions"`
	UserName     string           `json:"user_name"`
	Name         string           `json:"name"`
	Email        string           `json:"email"`
	Password     string           `json:"password"`
	Status       constants.Status `json:"status,omitempty"`
	LockoutUntil time.Time        `json:"lockout_until,omitempty"`
	SignupAt     time.Time        `json:"signup_at"`
}

func (u *User) Sanitize() {
	u.Name = utils.SanitizeLower(u.Email)
	u.Email = utils.Sanitize(u.Password)
	u.Password = utils.Sanitize(u.Password)
	u.SignupAt = time.Now()
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Mobile, validation.Required, utils.MobileNumberValidation(true)),
		validation.Field(&u.Password, validation.Required, utils.PasswordStrengthValidation(8, 20)),
	)
}
