package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestVersionHelp(t *testing.T) {
	cmd := &VersionCommand{}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !strings.Contains(help, "hcptf version") {
		t.Error("Help should contain usage")
	}
	if !strings.Contains(help, "Prints the version") {
		t.Error("Help should describe version command")
	}
}

func TestVersionSynopsis(t *testing.T) {
	cmd := &VersionCommand{}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
	if synopsis != "Prints the version" {
		t.Errorf("expected 'Prints the version', got %q", synopsis)
	}
}

func TestVersionFlagParsing(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "no flags",
			args: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &VersionCommand{}

			// Version command has no flags, but we still test the FlagSet can be created
			flags := cmd.Meta.FlagSet("version")

			if err := flags.Parse(tt.args); err != nil {
				t.Fatalf("flag parsing failed: %v", err)
			}
		})
	}
}

func TestVersionRun(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VersionCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "hcptf version") {
		t.Fatalf("expected version output, got %q", output)
	}

	if !strings.Contains(output, "0.1.0") {
		t.Fatalf("expected version number, got %q", output)
	}
}

func TestVersionRunWithArgs(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VersionCommand{
		Meta: newTestMeta(ui),
	}

	// Version command ignores arguments
	code := cmd.Run([]string{"--help"})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "hcptf version") {
		t.Fatalf("expected version output even with args, got %q", output)
	}
}
