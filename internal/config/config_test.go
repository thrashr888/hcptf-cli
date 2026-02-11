package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}

func TestLoadMergesConfigAndTerraformCredentials(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	configContent := `
credentials "app.terraform.io" {
  token = "hcptf-token"
}

default_organization = "hashicorp"
output_format = "json"
`
	writeFile(t, filepath.Join(home, ".hcptfrc"), configContent)

	tfCreds := TerraformCredentials{Credentials: map[string]TerraformCredential{
		"app.terraform.io":    {Token: "tf-token"},
		"private.example.com": {Token: "other-token"},
	}}
	data, err := json.Marshal(tfCreds)
	if err != nil {
		t.Fatalf("failed to marshal terraform creds: %v", err)
	}
	writeFile(t, filepath.Join(home, ".terraform.d", "credentials.tfrc.json"), string(data))

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.OutputFormat != "json" {
		t.Fatalf("expected output format json, got %s", cfg.OutputFormat)
	}

	if cfg.DefaultOrganization != "hashicorp" {
		t.Fatalf("expected default organization hashicorp, got %s", cfg.DefaultOrganization)
	}

	if got := cfg.Credentials["app.terraform.io"].Token; got != "hcptf-token" {
		t.Fatalf("expected hcptf token to override terraform token, got %s", got)
	}

	if got := cfg.Credentials["private.example.com"].Token; got != "other-token" {
		t.Fatalf("expected terraform credential to be added, got %s", got)
	}
}

func TestLoadDefaultsWhenConfigMissing(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.OutputFormat != "table" {
		t.Fatalf("expected default output format table, got %s", cfg.OutputFormat)
	}

	if len(cfg.Credentials) != 0 {
		t.Fatalf("expected no credentials, got %d", len(cfg.Credentials))
	}
}

func TestGetTokenPriority(t *testing.T) {
	cfg := &Config{Credentials: map[string]*Credential{
		"app.terraform.io": {Hostname: "app.terraform.io", Token: "config-token"},
	}}

	t.Setenv("TFE_TOKEN", "tfe-env")
	if token := cfg.GetToken("app.terraform.io"); token != "tfe-env" {
		t.Fatalf("expected TFE_TOKEN to be used, got %s", token)
	}

	t.Setenv("TFE_TOKEN", "")
	t.Setenv("HCPTF_TOKEN", "hcptf-env")
	if token := cfg.GetToken("app.terraform.io"); token != "hcptf-env" {
		t.Fatalf("expected HCPTF_TOKEN to be used, got %s", token)
	}

	t.Setenv("HCPTF_TOKEN", "")
	if token := cfg.GetToken("app.terraform.io"); token != "config-token" {
		t.Fatalf("expected config credential to be used, got %s", token)
	}
}

func TestGetAddress(t *testing.T) {
	t.Setenv("HCPTF_ADDRESS", "https://example.com")
	if addr := GetAddress(); addr != "https://example.com" {
		t.Fatalf("expected custom address, got %s", addr)
	}

	t.Setenv("HCPTF_ADDRESS", "")
	if addr := GetAddress(); addr != "https://app.terraform.io" {
		t.Fatalf("expected default address, got %s", addr)
	}
}

func TestSaveAndRemoveCredential(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SaveCredential("app.terraform.io", "token-one"); err != nil {
		t.Fatalf("SaveCredential error = %v", err)
	}

	if err := SaveCredential("private.example.com", "token-two"); err != nil {
		t.Fatalf("SaveCredential error = %v", err)
	}

	credsPath := filepath.Join(home, ".terraform.d", "credentials.tfrc.json")
	data, err := os.ReadFile(credsPath)
	if err != nil {
		t.Fatalf("failed to read credentials: %v", err)
	}

	var stored TerraformCredentials
	if err := json.Unmarshal(data, &stored); err != nil {
		t.Fatalf("failed to unmarshal credentials: %v", err)
	}

	if stored.Credentials["app.terraform.io"].Token != "token-one" {
		t.Fatalf("expected token-one, got %s", stored.Credentials["app.terraform.io"].Token)
	}

	if err := RemoveCredential("app.terraform.io"); err != nil {
		t.Fatalf("RemoveCredential error = %v", err)
	}

	data, err = os.ReadFile(credsPath)
	if err != nil {
		t.Fatalf("failed to read credentials after removal: %v", err)
	}

	stored = TerraformCredentials{}
	if err := json.Unmarshal(data, &stored); err != nil {
		t.Fatalf("failed to unmarshal credentials after removal: %v", err)
	}

	if _, exists := stored.Credentials["app.terraform.io"]; exists {
		t.Fatalf("expected app.terraform.io to be removed")
	}

	if err := RemoveCredential("private.example.com"); err != nil {
		t.Fatalf("RemoveCredential error = %v", err)
	}

	if _, err := os.Stat(credsPath); !os.IsNotExist(err) {
		t.Fatalf("expected credentials file to be removed when empty")
	}
}

func TestGetConfigPathEnvOverride(t *testing.T) {
	t.Setenv("HCPTF_CONFIG", "/tmp/custom")
	if path := GetConfigPath(); path != "/tmp/custom" {
		t.Fatalf("expected HCPTF_CONFIG to override path, got %s", path)
	}
}
