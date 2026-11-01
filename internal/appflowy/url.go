package appflowy

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ValidateBaseURL checks scheme (http/https), requires a host, forbids userinfo,
// and strips query and fragment from the normalized URL (no trailing slash).
func ValidateBaseURL(raw string) (normalized string, insecureHTTP bool, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false, errors.New("APPFLOWY_URL is empty")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return "", false, fmt.Errorf("APPFLOWY_URL is not a valid URL: %w", err)
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "http", "https":
	default:
		return "", false, fmt.Errorf("APPFLOWY_URL must use http or https (got %q)", u.Scheme)
	}
	if u.Host == "" {
		return "", false, errors.New("APPFLOWY_URL must include a host")
	}
	if u.User != nil {
		return "", false, errors.New("APPFLOWY_URL must not include userinfo")
	}
	u.RawQuery = ""
	u.Fragment = ""
	normalized = strings.TrimRight(u.String(), "/")
	return normalized, scheme == "http", nil
}
