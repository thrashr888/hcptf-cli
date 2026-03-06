package validate

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var idPattern = regexp.MustCompile(`^[a-zA-Z]+-[a-zA-Z0-9]+$`)

// ID validates a TFE-like resource ID.
func ID(value, flagName string) error {
	if value == "" {
		return nil
	}

	if len(value) > 64 {
		return fmt.Errorf("%s must be at most 64 characters", flagName)
	}

	if strings.ContainsAny(value, `/\\`) {
		return fmt.Errorf("%s contains path separator", flagName)
	}

	if !idPattern.MatchString(value) {
		return fmt.Errorf("%s must match pattern %q", flagName, idPattern.String())
	}

	return nil
}

// Name validates resource names and free-form identifiers.
func Name(value, flagName string) error {
	if value == "" {
		return nil
	}

	if len(value) > 256 {
		return fmt.Errorf("%s must be at most 256 characters", flagName)
	}

	if strings.Contains(value, "../") || strings.Contains(value, "..\\") {
		return fmt.Errorf("%s contains path traversal", flagName)
	}

	if strings.ContainsAny(value, `/\\`) {
		return fmt.Errorf("%s contains path separator", flagName)
	}

	if strings.ContainsAny(value, "?&") {
		return fmt.Errorf("%s contains invalid query character", flagName)
	}

	if hasURLEncodedSequence(value) {
		return fmt.Errorf("%s contains URL-encoded sequence", flagName)
	}

	if hasControlChars(value) {
		return fmt.Errorf("%s contains control characters", flagName)
	}

	return nil
}

// SafeString validates user-entered textual data such as descriptions.
func SafeString(value, flagName string) error {
	if value == "" {
		return nil
	}

	if len(value) > 4096 {
		return fmt.Errorf("%s must be at most 4096 characters", flagName)
	}

	if hasUnsafeControlChars(value) {
		return fmt.Errorf("%s contains control characters", flagName)
	}

	return nil
}

func hasURLEncodedSequence(value string) bool {
	for i := 0; i < len(value)-2; i++ {
		if value[i] != '%' {
			continue
		}

		h := value[i+1]
		t := value[i+2]
		if isHexDigit(h) && isHexDigit(t) {
			return true
		}
	}
	return false
}

func isHexDigit(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
}

func hasControlChars(value string) bool {
	for _, r := range value {
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}

func hasUnsafeControlChars(value string) bool {
	for _, r := range value {
		if r == '\n' || r == '\r' || r == '\t' {
			continue
		}
		if unicode.IsControl(r) {
			return true
		}
	}
	return false
}
