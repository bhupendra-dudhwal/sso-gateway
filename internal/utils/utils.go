package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
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
