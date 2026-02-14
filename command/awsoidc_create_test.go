package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAWSOIDCCreateRequiresOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestAWSOIDCCreateRequiresRoleArn(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-role-arn") {
		t.Fatalf("expected role-arn error, got %q", out)
	}
}

func TestAWSOIDCCreateRequiresEmptyOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=", "-role-arn=arn:aws:iam::123:role/test"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-organization") {
		t.Fatalf("expected organization error, got %q", out)
	}
}

func TestAWSOIDCCreateRequiresEmptyRoleArn(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcCreateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-organization=my-org", "-role-arn="})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-role-arn") {
		t.Fatalf("expected role-arn error, got %q", out)
	}
}

func TestAWSOIDCCreateHelp(t *testing.T) {
	cmd := &AWSoidcCreateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf awsoidc create") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !strings.Contains(help, "-org") {
		t.Error("Help should mention -org flag alias")
	}
	if !strings.Contains(help, "-role-arn") {
		t.Error("Help should mention -role-arn flag")
	}
	if !strings.Contains(help, "-output") {
		t.Error("Help should mention -output flag")
	}
	if !strings.Contains(help, "required") {
		t.Error("Help should indicate required flags")
	}
	if !strings.Contains(help, "AWS OIDC") {
		t.Error("Help should describe AWS OIDC configuration")
	}
}

func TestAWSOIDCCreateSynopsis(t *testing.T) {
	cmd := &AWSoidcCreateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Create an AWS OIDC configuration for dynamic credentials" {
		t.Errorf("expected 'Create an AWS OIDC configuration for dynamic credentials', got %q", synopsis)
	}
}

func TestAWSOIDCCreateFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedOrg string
		expectedArn string
		expectedFmt string
	}{
		{
			name:        "organization and role-arn flags",
			args:        []string{"-organization=my-org", "-role-arn=arn:aws:iam::123456789012:role/terraform-role"},
			expectedOrg: "my-org",
			expectedArn: "arn:aws:iam::123456789012:role/terraform-role",
			expectedFmt: "table",
		},
		{
			name:        "org alias flag",
			args:        []string{"-org=test-org", "-role-arn=arn:aws:iam::987654321098:role/test-role"},
			expectedOrg: "test-org",
			expectedArn: "arn:aws:iam::987654321098:role/test-role",
			expectedFmt: "table",
		},
		{
			name:        "with json output",
			args:        []string{"-organization=my-org", "-role-arn=arn:aws:iam::123456789012:role/role", "-output=json"},
			expectedOrg: "my-org",
			expectedArn: "arn:aws:iam::123456789012:role/role",
			expectedFmt: "json",
		},
		{
			name:        "with table output",
			args:        []string{"-org=test-org", "-role-arn=arn:aws:iam::111122223333:role/myRole", "-output=table"},
			expectedOrg: "test-org",
			expectedArn: "arn:aws:iam::111122223333:role/myRole",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AWSoidcCreateCommand{}

			flags := cmd.Meta.FlagSet("awsoidc create")
			flags.StringVar(&cmd.organization, "organization", "", "Organization name (required)")
			flags.StringVar(&cmd.organization, "org", "", "Organization name (alias)")
			flags.StringVar(&cmd.roleArn, "role-arn", "", "AWS IAM role ARN (required)")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the organization was set correctly
			if cmd.organization != tt.expectedOrg {
				t.Errorf("expected organization %q, got %q", tt.expectedOrg, cmd.organization)
			}

			// Verify the role ARN was set correctly
			if cmd.roleArn != tt.expectedArn {
				t.Errorf("expected role ARN %q, got %q", tt.expectedArn, cmd.roleArn)
			}

			// Verify the format was set correctly
			if cmd.format != tt.expectedFmt {
				t.Errorf("expected format %q, got %q", tt.expectedFmt, cmd.format)
			}
		})
	}
}
