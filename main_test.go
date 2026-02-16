package main

import (
	"reflect"
	"testing"

	"github.com/mitchellh/cli"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name              string
		version           string
		versionPrerelease string
		expected          string
	}{
		{
			name:              "Release version",
			version:           "1.0.0",
			versionPrerelease: "",
			expected:          "1.0.0",
		},
		{
			name:              "Prerelease version",
			version:           "1.0.0",
			versionPrerelease: "beta",
			expected:          "1.0.0-beta",
		},
		{
			name:              "Dev version",
			version:           "0.1.0",
			versionPrerelease: "dev",
			expected:          "0.1.0-dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origVersion := Version
			origPrerelease := VersionPrerelease

			Version = tt.version
			VersionPrerelease = tt.versionPrerelease

			result := GetVersion()
			if result != tt.expected {
				t.Errorf("GetVersion() = %q, expected %q", result, tt.expected)
			}

			Version = origVersion
			VersionPrerelease = origPrerelease
		})
	}
}

func TestGetVersion_CurrentValues(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() returned empty string")
	}

	if Version == "" {
		t.Error("Version global variable is empty")
	}
}

func TestBuildGetVerbIndex(t *testing.T) {
	commands := map[string]cli.CommandFactory{
		"workspace list":          nil,
		"workspace read":          nil,
		"organization show":       nil,
		"organization member":     nil,
		"organization token list": nil,
		"team create":             nil,
	}

	index := buildGetVerbIndex(commands)

	workspace := index["workspace"]
	if !workspace.hasList || !workspace.hasRead || workspace.hasShow {
		t.Fatalf("unexpected workspace get verbs: %+v", workspace)
	}

	organization := index["organization"]
	if organization.hasList || organization.hasRead || !organization.hasShow {
		t.Fatalf("unexpected organization get verbs: %+v", organization)
	}

	if _, ok := index["team"]; ok {
		t.Fatalf("did not expect team in get verb index")
	}
}

func TestInferImplicitGetVerb(t *testing.T) {
	index := map[string]getVerbAvailability{
		"workspace":          {hasList: true, hasRead: true},
		"organization":       {hasList: true, hasShow: true},
		"user":               {hasRead: true},
		"featureset":         {hasList: true},
		"organization token": {hasList: true, hasRead: true},
	}

	tests := []struct {
		name      string
		input     []string
		expected  []string
		expectErr bool
	}{
		{
			name:     "empty args",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "help is unchanged",
			input:    []string{"workspace", "-h"},
			expected: []string{"workspace", "-h"},
		},
		{
			name:     "explicit get verb unchanged",
			input:    []string{"workspace", "list", "-org=myorg"},
			expected: []string{"workspace", "list", "-org=myorg"},
		},
		{
			name:     "workspace list inferred with collection selector",
			input:    []string{"workspace", "-org=myorg"},
			expected: []string{"workspace", "list", "-org=myorg"},
		},
		{
			name:     "workspace id infers read",
			input:    []string{"workspace", "-id=ws-123"},
			expected: []string{"workspace", "read", "-id=ws-123"},
		},
		{
			name:     "workspace name infers read",
			input:    []string{"workspace", "-name=my-workspace"},
			expected: []string{"workspace", "read", "-name=my-workspace"},
		},
		{
			name:      "workspace no selectors is ambiguous",
			input:     []string{"workspace"},
			expectErr: true,
		},
		{
			name:     "organization id infers show",
			input:    []string{"organization", "-id=org-123"},
			expected: []string{"organization", "show", "-id=org-123"},
		},
		{
			name:      "organization no selectors is ambiguous",
			input:     []string{"organization"},
			expectErr: true,
		},
		{
			name:     "single read operation",
			input:    []string{"user"},
			expected: []string{"user", "read"},
		},
		{
			name:     "single list operation",
			input:    []string{"featureset"},
			expected: []string{"featureset", "list"},
		},
		{
			name:     "nested namespace defaults to list with collection selector",
			input:    []string{"organization", "token", "-org=myorg"},
			expected: []string{"organization", "token", "list", "-org=myorg"},
		},
		{
			name:     "nested namespace id infers read",
			input:    []string{"organization", "token", "-id=ot-123"},
			expected: []string{"organization", "token", "read", "-id=ot-123"},
		},
		{
			name:     "unknown namespace unchanged",
			input:    []string{"notacommand"},
			expected: []string{"notacommand"},
		},
		{
			name:     "extra positional token unchanged",
			input:    []string{"workspace", "abc"},
			expected: []string{"workspace", "abc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := inferImplicitGetVerb(tt.input, index)
			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("inferImplicitGetVerb(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeDeleteConfirmationFlags(t *testing.T) {
	commands := map[string]cli.CommandFactory{
		"workspace delete":          nil,
		"organization token delete": nil,
		"workspace read":            nil,
	}

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "f becomes force on delete command",
			input:    []string{"workspace", "delete", "-f"},
			expected: []string{"workspace", "delete", "-force"},
		},
		{
			name:     "y becomes force on delete command",
			input:    []string{"workspace", "delete", "-y"},
			expected: []string{"workspace", "delete", "-force"},
		},
		{
			name:     "nested delete command normalized",
			input:    []string{"organization", "token", "delete", "-y"},
			expected: []string{"organization", "token", "delete", "-force"},
		},
		{
			name:     "non-delete command unchanged",
			input:    []string{"workspace", "read", "-y"},
			expected: []string{"workspace", "read", "-y"},
		},
		{
			name:     "help unchanged",
			input:    []string{"workspace", "delete", "-h"},
			expected: []string{"workspace", "delete", "-h"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeDeleteConfirmationFlags(tt.input, commands)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Fatalf("normalizeDeleteConfirmationFlags(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
