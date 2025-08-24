package utils

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidCacheKey    error = fmt.Errorf("key does not exists")
	ErrDuplicate          error = errors.New("document already exists")
	ErrDocumentNotFound   error = errors.New("document not found")
	ErrInvalidCredentials error = errors.New("Please enter a valid credentials")
)
