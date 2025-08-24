package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/valyala/fasthttp"
)

func GetField(ctx *fasthttp.RequestCtx, key constants.Context) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}

	val, ok := v.(string)
	if !ok {
		return ""
	}

	return val
}

// findProjectRoot walks up from the current directory to find a directory containing go.mod.
func FindProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	for {
		if FileExists(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // reached root "/"
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found in any parent directory")
}

// fileExists checks whether a file exists at the given path.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// SanitizeLower trims, lowercases, and sanitizes UTF-8
func SanitizeLower(str string) string {
	clean := strings.TrimSpace(str)
	clean = strings.ToLower(clean)
	return strings.ToValidUTF8(clean, "")
}

// Sanitize trims and converts to valid UTF-8
func Sanitize(str string) string {
	clean := strings.TrimSpace(str)
	return strings.ToValidUTF8(clean, "")
}

// SanitizeLowerSlice applies SanitizeLower to a slice
func SanitizeLowerSlice(strs []string) []string {
	sanitized := make([]string, 0, len(strs))
	for _, val := range strs {
		clean := SanitizeLower(val)
		if clean != "" {
			sanitized = append(sanitized, clean)
		}
	}
	return sanitized
}

// SanitizeSlice applies Sanitize to a slice
func SanitizeSlice(strs []string) []string {
	sanitized := make([]string, 0, len(strs))
	for _, val := range strs {
		clean := Sanitize(val)
		if clean != "" {
			sanitized = append(sanitized, clean)
		}
	}
	return sanitized
}

// Rule: Mobile number validation for int64 fields
func MobileNumberValidation(isRequired bool) validation.Rule {
	return validation.By(func(value any) error {
		// If required, check for nil / zero value
		if isRequired {
			if value == nil {
				return errors.New("Mobile number is required")
			}
			if v, ok := value.(int64); ok && v == 0 {
				return errors.New("Mobile number is required")
			}
		}

		// If value is zero and not required, skip further checks
		v, ok := value.(int64)
		if !ok {
			return errors.New("Invalid mobile number type")
		}

		// Only validate if value is set (non-zero)
		if v != 0 {
			mobileStr := fmt.Sprintf("%d", v)
			if !IsValidMobile(mobileStr) {
				return errors.New("Invalid mobile number format")
			}
		}

		return nil
	})
}

func IsValidMobileFromInt64(mobile int64) bool {
	return IsValidMobile(fmt.Sprintf("%d", mobile))
}

// Actual mobile number format check (e.g., 10-digit Indian mobile numbers)
func IsValidMobile(mobile string) bool {
	re := regexp.MustCompile(`^[6-9][0-9]{9}$`)
	return re.MatchString(mobile)
}

func PasswordStrengthValidation(min, max int) validation.Rule {
	return validation.By(func(value any) error {
		password, ok := value.(string)
		if !ok {
			return errors.New("Invalid password type")
		}

		length := len(password)
		if length < min || length > max {
			return errors.New("Password length must be between 8 and 20 characters")
		}

		var (
			hasUpper   = regexp.MustCompile(`[A-Z]`).MatchString(password)
			hasLower   = regexp.MustCompile(`[a-z]`).MatchString(password)
			hasDigit   = regexp.MustCompile(`[0-9]`).MatchString(password)
			hasSpecial = regexp.MustCompile(`[_\-@]`).MatchString(password) // Only _, -, @ allowed
		)

		if !hasUpper {
			return errors.New("Password must contain at least one uppercase letter")
		}
		if !hasLower {
			return errors.New("Password must contain at least one lowercase letter")
		}
		if !hasDigit {
			return errors.New("Password must contain at least one digit")
		}
		if !hasSpecial {
			return errors.New("Password must contain at least one special character")
		}

		return nil
	})
}
