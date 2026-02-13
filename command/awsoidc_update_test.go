package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestAWSOIDCUpdateRequiresID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestAWSOIDCUpdateValidatesEmptyID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &AWSoidcUpdateCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-id=", "-role-arn=arn:aws:iam::123456789012:role/test"})
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}

	if out := ui.ErrorWriter.String(); !strings.Contains(out, "-id") {
		t.Fatalf("expected id error, got %q", out)
	}
}

func TestAWSOIDCUpdateHelp(t *testing.T) {
	cmd := &AWSoidcUpdateCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf awsoidc update") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "-id") {
		t.Error("Help should mention -id flag")
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
	if !strings.Contains(help, "AWS OIDC configuration") {
		t.Error("Help should describe AWS OIDC configuration")
	}
}

func TestAWSOIDCUpdateSynopsis(t *testing.T) {
	cmd := &AWSoidcUpdateCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Update AWS OIDC configuration settings" {
		t.Errorf("expected 'Update AWS OIDC configuration settings', got %q", synopsis)
	}
}

func TestAWSOIDCUpdateFlagParsing(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectedID  string
		expectedArn string
		expectedFmt string
	}{
		{
			name:        "id flag only",
			args:        []string{"-id=awsoidc-ABC123"},
			expectedID:  "awsoidc-ABC123",
			expectedArn: "",
			expectedFmt: "table",
		},
		{
			name:        "id and role-arn flags",
			args:        []string{"-id=awsoidc-XYZ789", "-role-arn=arn:aws:iam::123456789012:role/new-role"},
			expectedID:  "awsoidc-XYZ789",
			expectedArn: "arn:aws:iam::123456789012:role/new-role",
			expectedFmt: "table",
		},
		{
			name:        "with json output",
			args:        []string{"-id=awsoidc-DEF456", "-role-arn=arn:aws:iam::987654321098:role/test", "-output=json"},
			expectedID:  "awsoidc-DEF456",
			expectedArn: "arn:aws:iam::987654321098:role/test",
			expectedFmt: "json",
		},
		{
			name:        "with explicit table output",
			args:        []string{"-id=awsoidc-GHI789", "-output=table"},
			expectedID:  "awsoidc-GHI789",
			expectedArn: "",
			expectedFmt: "table",
		},
		{
			name:        "update only role arn",
			args:        []string{"-id=awsoidc-12345", "-role-arn=arn:aws:iam::111122223333:role/updated-role"},
			expectedID:  "awsoidc-12345",
			expectedArn: "arn:aws:iam::111122223333:role/updated-role",
			expectedFmt: "table",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &AWSoidcUpdateCommand{}

			flags := cmd.Meta.FlagSet("awsoidc update")
			flags.StringVar(&cmd.id, "id", "", "AWS OIDC configuration ID (required)")
			flags.StringVar(&cmd.roleArn, "role-arn", "", "AWS IAM role ARN")
			flags.StringVar(&cmd.format, "output", "table", "Output format: table or json")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}

			// Verify the ID was set correctly
			if cmd.id != tt.expectedID {
				t.Errorf("expected ID %q, got %q", tt.expectedID, cmd.id)
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
