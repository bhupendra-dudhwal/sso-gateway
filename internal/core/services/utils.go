package services

import (
	"context"
	"time"
)

// Slice of stling to map[strng]struct{}
func sliceStringToMapStruct(perms []string) map[string]struct{} {
	result := make(map[string]struct{}, len(perms))
	for _, p := range perms {
		result[p] = struct{}{}
	}
	return result
}

func withTimeout(ctx context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, duration)
}
