package validate

import (
	"strings"
	"testing"
)

func TestIDValid(t *testing.T) {
	valid := []string{
		"ws-abc123",
		"run-xyz789",
		"var-ABC",
		"pol-a1b2c3",
		"",
	}
	for _, id := range valid {
		if err := ID(id, "id"); err != nil {
			t.Errorf("ID(%q) = %v, want nil", id, err)
		}
	}
}

func TestIDInvalid(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"ws_abc", "must match pattern"},
		{"../../etc/passwd", "path separator"},
		{"ws/abc", "path separator"},
		{"ws\\abc", "path separator"},
		{"no-prefix-number-123", "must match pattern"},
		{"justletters", "must match pattern"},
		{strings.Repeat("a", 65), "at most 64"},
		{"ws-valid" + strings.Repeat("x", 60), "at most 64"},
	}
	for _, tt := range tests {
		err := ID(tt.value, "id")
		if err == nil {
			t.Errorf("ID(%q) = nil, want error containing %q", tt.value, tt.want)
			continue
		}
		if !strings.Contains(err.Error(), tt.want) {
			t.Errorf("ID(%q) = %v, want error containing %q", tt.value, err, tt.want)
		}
	}
}

func TestNameValid(t *testing.T) {
	valid := []string{
		"workspace-1",
		"my_workspace",
		"Prod-Environment",
		"ws.with.dots",
		"simple",
		"",
	}
	for _, name := range valid {
		if err := Name(name, "name"); err != nil {
			t.Errorf("Name(%q) = %v, want nil", name, err)
		}
	}
}

func TestNameInvalid(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"../../../etc/passwd", "path traversal"},
		{"foo/../bar", "path traversal"},
		{"..\\windows\\system", "path traversal"},
		{"bad/name", "path separator"},
		{"bad\\name", "path separator"},
		{"name?x=1", "query character"},
		{"name&y=2", "query character"},
		{"%2fadmin", "URL-encoded"},
		{"%00null", "URL-encoded"},
		{string([]byte{'n', 'a', 0x00, 'm', 'e'}), "control characters"},
		{string([]byte{'t', 'e', 0x1b, 's', 't'}), "control characters"},
		{strings.Repeat("a", 257), "at most 256"},
	}
	for _, tt := range tests {
		err := Name(tt.value, "name")
		if err == nil {
			t.Errorf("Name(%q) = nil, want error containing %q", tt.value, tt.want)
			continue
		}
		if !strings.Contains(err.Error(), tt.want) {
			t.Errorf("Name(%q) = %v, want error containing %q", tt.value, err, tt.want)
		}
	}
}

func TestSafeStringValid(t *testing.T) {
	valid := []string{
		"simple value",
		"line1\nline2",
		"tabs\tare\tok",
		"carriage\rreturn",
		"",
		strings.Repeat("x", 4096),
	}
	for _, s := range valid {
		if err := SafeString(s, "desc"); err != nil {
			t.Errorf("SafeString(%q) = %v, want nil", s, err)
		}
	}
}

func TestSafeStringInvalid(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{string([]byte{'a', 0x00, 'b'}), "control characters"},
		{string([]byte{'a', 0x1b, 'b'}), "control characters"},
		{string([]byte{'a', 0x07, 'b'}), "control characters"},
		{strings.Repeat("x", 4097), "at most 4096"},
	}
	for _, tt := range tests {
		err := SafeString(tt.value, "desc")
		if err == nil {
			t.Errorf("SafeString(%q) = nil, want error containing %q", tt.value, tt.want)
			continue
		}
		if !strings.Contains(err.Error(), tt.want) {
			t.Errorf("SafeString(%q) = %v, want error containing %q", tt.value, err, tt.want)
		}
	}
}

func TestURLEncodedDetection(t *testing.T) {
	tests := []struct {
		value  string
		expect bool
	}{
		{"%2f", true},
		{"%00", true},
		{"%7A", true},
		{"%2Fadmin", true},
		{"normal", false},
		{"100%", false},
		{"%%", false},
		{"%zz", false},
		{"%", false},
		{"%a", false},
	}
	for _, tt := range tests {
		got := hasURLEncodedSequence(tt.value)
		if got != tt.expect {
			t.Errorf("hasURLEncodedSequence(%q) = %v, want %v", tt.value, got, tt.expect)
		}
	}
}

func TestFlagNameInErrorMessage(t *testing.T) {
	err := ID("bad_value", "-my-flag")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "-my-flag") {
		t.Errorf("error should contain flag name, got %q", err.Error())
	}
}
