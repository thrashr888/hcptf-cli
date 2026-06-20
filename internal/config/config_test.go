package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

func chdir(t *testing.T, dir string) {
	t.Helper()
	previous, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previous); err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	})
}

func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	previous := make(map[string]string, len(keys))
	present := make(map[string]bool, len(keys))
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if ok {
			previous[key] = value
			present[key] = true
		}
		if err := os.Unsetenv(key); err != nil {
			t.Fatalf("failed to unset %s: %v", key, err)
		}
	}
	t.Cleanup(func() {
		for _, key := range keys {
			var err error
			if present[key] {
				err = os.Setenv(key, previous[key])
			} else {
				err = os.Unsetenv(key)
			}
			if err != nil {
				t.Fatalf("failed to restore %s: %v", key, err)
			}
		}
	})
}

func TestLoadMergesConfigAndTerraformCredentials(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_ORG", "HCPTF_ORG")
	home := t.TempDir()
	chdir(t, t.TempDir())
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
	unsetEnv(t, EnvFileVariable, "TFE_ORG", "HCPTF_ORG")
	home := t.TempDir()
	chdir(t, t.TempDir())
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

func TestLoadUsesDefaultOrganizationFromDotEnv(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_ORG", "HCPTF_ORG")
	home := t.TempDir()
	dir := t.TempDir()
	chdir(t, dir)
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(dir, ".env"), "TFE_ORG=dotenv-org\n")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DefaultOrganization != "dotenv-org" {
		t.Fatalf("expected default organization dotenv-org, got %s", cfg.DefaultOrganization)
	}
}

func TestLoadEnvironmentOrganizationOverridesConfig(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_ORG", "HCPTF_ORG")
	home := t.TempDir()
	chdir(t, t.TempDir())
	t.Setenv("HOME", home)
	t.Setenv("TFE_ORG", "env-org")

	writeFile(t, filepath.Join(home, ".hcptfrc"), `
default_organization = "config-org"
`)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DefaultOrganization != "env-org" {
		t.Fatalf("expected default organization env-org, got %s", cfg.DefaultOrganization)
	}
}

func TestLoadDotEnvLoadsDefaultFile(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN", "HCPTF_TOKEN", "TFE_ADDRESS", "HCPTF_ADDRESS")
	dir := t.TempDir()
	chdir(t, dir)
	writeFile(t, filepath.Join(dir, ".env"), `
# Terraform Enterprise connection
TFE_TOKEN=dotenv-token
TFE_ADDRESS=https://tfe.example.com
`)

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("TFE_TOKEN"); got != "dotenv-token" {
		t.Fatalf("expected TFE_TOKEN from .env, got %q", got)
	}
	if got := os.Getenv("TFE_ADDRESS"); got != "https://tfe.example.com" {
		t.Fatalf("expected TFE_ADDRESS from .env, got %q", got)
	}
}

func TestLoadDotEnvDoesNotOverrideExistingEnvironment(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	dir := t.TempDir()
	chdir(t, dir)
	writeFile(t, filepath.Join(dir, ".env"), "TFE_TOKEN=dotenv-token\n")
	t.Setenv("TFE_TOKEN", "exported-token")

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("TFE_TOKEN"); got != "exported-token" {
		t.Fatalf("expected exported TFE_TOKEN to win, got %q", got)
	}
}

func TestLoadDotEnvLoadsAncestorFile(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	home := t.TempDir()
	root := filepath.Join(home, "project")
	nested := filepath.Join(root, "lab", "b194")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("failed to create nested directory: %v", err)
	}
	chdir(t, nested)
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(root, ".env"), "TFE_TOKEN=ancestor-token\n")

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("TFE_TOKEN"); got != "ancestor-token" {
		t.Fatalf("expected TFE_TOKEN from ancestor .env, got %q", got)
	}
}

func TestLoadDotEnvLoadsUserDefaults(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	home := t.TempDir()
	work := filepath.Join(home, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatalf("failed to create work directory: %v", err)
	}
	chdir(t, work)
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(home, ".hcptf.env"), "TFE_TOKEN=user-token\n")

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("TFE_TOKEN"); got != "user-token" {
		t.Fatalf("expected TFE_TOKEN from user defaults, got %q", got)
	}
}

func TestLoadDotEnvProjectFileWinsOverUserDefaults(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN", "TFE_ORG")
	home := t.TempDir()
	project := filepath.Join(home, "project")
	if err := os.MkdirAll(project, 0o755); err != nil {
		t.Fatalf("failed to create project directory: %v", err)
	}
	chdir(t, project)
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(home, ".hcptf.env"), "TFE_TOKEN=user-token\nTFE_ORG=user-org\n")
	writeFile(t, filepath.Join(project, ".env"), "TFE_ORG=project-org\n")

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("TFE_TOKEN"); got != "user-token" {
		t.Fatalf("expected TFE_TOKEN from user defaults, got %q", got)
	}
	if got := os.Getenv("TFE_ORG"); got != "project-org" {
		t.Fatalf("expected project TFE_ORG to win, got %q", got)
	}
}

func TestLoadDotEnvExplicitFileTakesPrecedenceOverDefaultFile(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	dir := t.TempDir()
	chdir(t, dir)
	writeFile(t, filepath.Join(dir, ".env"), "TFE_TOKEN=default-token\n")
	explicitPath := filepath.Join(dir, "prod.env")
	writeFile(t, explicitPath, "TFE_TOKEN=explicit-token\n")
	t.Setenv(EnvFileVariable, explicitPath)

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}

	if got := os.Getenv("TFE_TOKEN"); got != "explicit-token" {
		t.Fatalf("expected explicit env file token to win, got %q", got)
	}
}

func TestLoadDotEnvMissingDefaultFileIsIgnored(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	dir := t.TempDir()
	chdir(t, dir)

	if err := LoadDotEnv(); err != nil {
		t.Fatalf("LoadDotEnv() error = %v", err)
	}
}

func TestLoadDotEnvMissingExplicitFileReturnsError(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	dir := t.TempDir()
	chdir(t, dir)
	t.Setenv(EnvFileVariable, filepath.Join(dir, "missing.env"))

	err := LoadDotEnv()
	if err == nil {
		t.Fatal("expected missing explicit env file to return error")
	}
	if !strings.Contains(err.Error(), "failed to load env file") {
		t.Fatalf("expected env file error, got %v", err)
	}
}

func TestLoadDotEnvMalformedFileReturnsError(t *testing.T) {
	unsetEnv(t, EnvFileVariable, "TFE_TOKEN")
	dir := t.TempDir()
	chdir(t, dir)
	envPath := filepath.Join(dir, "bad.env")
	writeFile(t, envPath, "not valid dotenv\n")
	t.Setenv(EnvFileVariable, envPath)

	err := LoadDotEnv()
	if err == nil {
		t.Fatal("expected malformed env file to return error")
	}
	if !strings.Contains(err.Error(), "failed to load env file") {
		t.Fatalf("expected env file error, got %v", err)
	}
}

func TestLoadInvalidTerraformCredentialsReturnsError(t *testing.T) {
	unsetEnv(t, EnvFileVariable)
	home := t.TempDir()
	chdir(t, t.TempDir())
	t.Setenv("HOME", home)

	writeFile(t, filepath.Join(home, ".terraform.d", "credentials.tfrc.json"), `{"credentials": invalid-json}`)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error loading invalid terraform credentials")
	}

	if !strings.Contains(err.Error(), "failed to load Terraform credentials") {
		t.Fatalf("expected parse-related error, got %v", err)
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
	tests := []struct {
		name            string
		hcptfAddress    string
		tfeAddress      string
		expectedAddress string
	}{
		{
			name:            "HCPTF_ADDRESS takes precedence",
			hcptfAddress:    "https://hcptf.example.com",
			tfeAddress:      "https://tfe.example.com",
			expectedAddress: "https://hcptf.example.com",
		},
		{
			name:            "TFE_ADDRESS used as fallback",
			hcptfAddress:    "",
			tfeAddress:      "https://tfe.example.com",
			expectedAddress: "https://tfe.example.com",
		},
		{
			name:            "default address when neither set",
			hcptfAddress:    "",
			tfeAddress:      "",
			expectedAddress: "https://app.terraform.io",
		},
		{
			name:            "HCPTF_ADDRESS alone works",
			hcptfAddress:    "https://custom.example.com",
			tfeAddress:      "",
			expectedAddress: "https://custom.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hcptfAddress != "" {
				t.Setenv("HCPTF_ADDRESS", tt.hcptfAddress)
			} else {
				t.Setenv("HCPTF_ADDRESS", "")
			}

			if tt.tfeAddress != "" {
				t.Setenv("TFE_ADDRESS", tt.tfeAddress)
			} else {
				t.Setenv("TFE_ADDRESS", "")
			}

			addr := GetAddress()
			if addr != tt.expectedAddress {
				t.Errorf("expected address %q, got %q", tt.expectedAddress, addr)
			}
		})
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
