package utils

import "fmt"

var (
	ErrInvalidCacheKey error = fmt.Errorf("key does not exists")
)
