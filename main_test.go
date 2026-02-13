package main

import (
	"testing"
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
			// Save original values
			origVersion := Version
			origPrerelease := VersionPrerelease

			// Set test values
			Version = tt.version
			VersionPrerelease = tt.versionPrerelease

			// Test GetVersion
			result := GetVersion()

			if result != tt.expected {
				t.Errorf("GetVersion() = %q, expected %q", result, tt.expected)
			}

			// Restore original values
			Version = origVersion
			VersionPrerelease = origPrerelease
		})
	}
}

func TestGetVersion_CurrentValues(t *testing.T) {
	// Test with actual current values to ensure they're set
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() returned empty string")
	}

	// Version should contain at least the Version value
	if Version == "" {
		t.Error("Version global variable is empty")
	}
}
