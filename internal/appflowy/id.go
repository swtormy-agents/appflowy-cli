package appflowy

import (
	"errors"
	"regexp"
	"strconv"
)

var uuidRE = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// ParseUUID accepts only canonical UUID strings for use in API path segments.
func ParseUUID(s string) (string, error) {
	s = stringsTrimSpace(s)
	if s == "" {
		return "", errors.New("UUID is empty")
	}
	if !uuidRE.MatchString(s) {
		return "", errors.New("UUID must match 8-4-4-4-12 hex format")
	}
	return stringsToLowerASCII(s), nil
}

// ParseDepth accepts folder tree depth (1–10 per AppFlowy Cloud API behaviour).
func ParseDepth(s string) (string, error) {
	n, err := strconv.ParseUint(stringsTrimSpace(s), 10, 32)
	if err != nil || n == 0 || n > 10 {
		return "", errors.New("depth must be an integer from 1 to 10")
	}
	return strconv.FormatUint(n, 10), nil
}

func stringsTrimSpace(s string) string {
	for len(s) > 0 && (s[0] == ' ' || s[0] == '\t' || s[0] == '\n' || s[0] == '\r') {
		s = s[1:]
	}
	for len(s) > 0 {
		c := s[len(s)-1]
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}

func stringsToLowerASCII(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}
