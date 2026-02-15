package command

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestLoginHelp(t *testing.T) {
	cmd := &LoginCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	if !strings.Contains(help, "hcptf login") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-hostname") {
		t.Error("Help should mention -hostname flag")
	}
	if !strings.Contains(help, "app.terraform.io") {
		t.Error("Help should mention default hostname")
	}
	if !strings.Contains(help, "credentials.tfrc.json") {
		t.Error("Help should mention credentials file")
	}
	if !strings.Contains(help, "-show-token") {
		t.Error("Help should mention -show-token flag")
	}
}

func TestLoginSynopsis(t *testing.T) {
	cmd := &LoginCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Authenticate to HCP Terraform" {
		t.Errorf("expected 'Authenticate to HCP Terraform', got %q", synopsis)
	}
}

func TestLoginFlagParsing(t *testing.T) {
	tests := []struct {
		name              string
		args              []string
		expectedHostname  string
		expectedShowToken bool
	}{
		{
			name:              "default hostname",
			args:              []string{},
			expectedHostname:  "app.terraform.io",
			expectedShowToken: false,
		},
		{
			name:              "custom hostname",
			args:              []string{"-hostname=custom.terraform.io"},
			expectedHostname:  "custom.terraform.io",
			expectedShowToken: false,
		},
		{
			name:              "enterprise hostname",
			args:              []string{"-hostname=tfe.company.com"},
			expectedHostname:  "tfe.company.com",
			expectedShowToken: false,
		},
		{
			name:              "show token flag",
			args:              []string{"-show-token"},
			expectedHostname:  "app.terraform.io",
			expectedShowToken: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &LoginCommand{}

			flags := cmd.Meta.FlagSet("login")
			flags.StringVar(&cmd.hostname, "hostname", "app.terraform.io", "HCP Terraform hostname")
			flags.BoolVar(&cmd.showToken, "show-token", false, "Show token after successful login")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			if cmd.hostname != tt.expectedHostname {
				t.Errorf("expected hostname %q, got %q", tt.expectedHostname, cmd.hostname)
			}
			if cmd.showToken != tt.expectedShowToken {
				t.Errorf("expected show-token %v, got %v", tt.expectedShowToken, cmd.showToken)
			}
		})
	}
}

func TestLoginShowTokenFromConfig(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	credsPath := filepath.Join(tmpDir, ".terraform.d", "credentials.tfrc.json")
	if err := ensureDir(filepath.Dir(credsPath)); err != nil {
		t.Fatalf("failed to create terraform credential directory: %v", err)
	}

	raw := map[string]interface{}{
		"credentials": map[string]interface{}{
			"app.terraform.io": map[string]interface{}{
				"token": "show-token-value",
			},
		},
	}
	encoded, err := json.Marshal(raw)
	if err != nil {
		t.Fatalf("failed to marshal credentials: %v", err)
	}
	if err := os.WriteFile(credsPath, encoded, 0o600); err != nil {
		t.Fatalf("failed to write credentials: %v", err)
	}

	ui := cli.NewMockUi()
	cmd := &LoginCommand{Meta: newTestMeta(ui)}

	if got := cmd.Run([]string{"-show-token"}); got != 0 {
		t.Fatalf("expected exit 0, got %d", got)
	}

	if !strings.Contains(ui.OutputWriter.String(), "show-token-value") {
		t.Fatalf("expected raw token in output, got %q", ui.OutputWriter.String())
	}
}

func TestLoginShowTokenMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	ui := cli.NewMockUi()
	cmd := &LoginCommand{Meta: newTestMeta(ui)}

	got := cmd.Run([]string{"-show-token"})
	if got != 1 {
		t.Fatalf("expected exit 1, got %d", got)
	}
	if !strings.Contains(ui.ErrorWriter.String(), "No token found") {
		t.Fatalf("expected missing token error, got %q", ui.ErrorWriter.String())
	}
}

func ensureDir(path string) error {
	return os.MkdirAll(path, 0o700)
}
