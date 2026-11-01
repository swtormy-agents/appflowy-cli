package appflowy

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ValidateAPIPath validates a raw API path for the `api` command.
// The path must start with /api/, must not contain query or fragment,
// and after percent-decoding must not contain "..".
func ValidateAPIPath(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", errors.New("path is empty")
	}
	if !strings.HasPrefix(raw, "/") {
		return "", errors.New("path must start with /")
	}
	if strings.Contains(raw, "?") || strings.Contains(raw, "#") {
		return "", errors.New("path must not contain query or fragment; use --query")
	}
	decoded, err := url.PathUnescape(raw)
	if err != nil {
		return "", fmt.Errorf("invalid path encoding: %w", err)
	}
	if strings.Contains(decoded, "..") {
		return "", errors.New("path must not contain ..")
	}
	if !strings.HasPrefix(decoded, "/api/") {
		return "", errors.New("path must start with /api/")
	}
	return decoded, nil
}
