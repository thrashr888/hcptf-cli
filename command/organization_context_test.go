package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestOrganizationContextCommand(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationContextCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-org=testorg"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "testorg") {
		t.Errorf("expected output to contain organization name, got: %s", output)
	}
	if !strings.Contains(output, "Workspaces:") {
		t.Errorf("expected output to contain Workspaces section, got: %s", output)
	}
	if !strings.Contains(output, "Projects:") {
		t.Errorf("expected output to contain Projects section, got: %s", output)
	}
	if !strings.Contains(output, "Teams:") {
		t.Errorf("expected output to contain Teams section, got: %s", output)
	}
}

func TestOrganizationContextCommandNoOrg(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationContextCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	errOutput := ui.ErrorWriter.String()
	if !strings.Contains(errOutput, "No organization specified") {
		t.Errorf("expected error about missing organization, got: %s", errOutput)
	}
}

func TestOrganizationContextHelp(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationContextCommand{
		Meta: newTestMeta(ui),
	}

	help := cmd.Help()
	if help == "" {
		t.Error("expected non-empty help text")
	}
}

func TestOrganizationContextSynopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &OrganizationContextCommand{
		Meta: newTestMeta(ui),
	}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Error("expected non-empty synopsis")
	}
}
