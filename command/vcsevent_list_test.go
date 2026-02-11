package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestVCSEventListCommand_Help(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventListCommand{
		Meta: testMeta(t, ui),
	}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !contains(help, "Usage: hcptf vcsevent list") {
		t.Error("Help should contain usage")
	}
	if !contains(help, "-organization") {
		t.Error("Help should mention -organization flag")
	}
	if !contains(help, "VCS events") {
		t.Error("Help should describe VCS events")
	}
}

func TestVCSEventListCommand_Synopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventListCommand{
		Meta: testMeta(t, ui),
	}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
}

func TestVCSEventListCommand_Run_NoOrganization(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventListCommand{
		Meta: testMeta(t, ui),
	}

	code := cmd.Run([]string{})
	if code == 0 {
		t.Fatal("Expected non-zero exit code when organization is missing")
	}

	output := ui.ErrorWriter.String()
	if !contains(output, "organization") {
		t.Error("Error should mention missing organization flag")
	}
}

func TestVCSEventListCommand_Run_InvalidLevel(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventListCommand{
		Meta: testMeta(t, ui),
	}

	code := cmd.Run([]string{"-org=test-org", "-level=invalid"})
	if code == 0 {
		t.Fatal("Expected non-zero exit code when level is invalid")
	}

	output := ui.ErrorWriter.String()
	if !contains(output, "level must be either 'info' or 'error'") {
		t.Error("Error should mention valid level values")
	}
}
