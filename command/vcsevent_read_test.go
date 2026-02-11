package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestVCSEventReadCommand_Help(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventReadCommand{
		Meta: testMeta(t, ui),
	}

	help := cmd.Help()
	if help == "" {
		t.Fatal("Help should not be empty")
	}

	// Check for key help elements
	if !contains(help, "Usage: hcptf vcsevent read") {
		t.Error("Help should contain usage")
	}
	if !contains(help, "-id") {
		t.Error("Help should mention -id flag")
	}
	if !contains(help, "VCS event") {
		t.Error("Help should describe VCS event")
	}
}

func TestVCSEventReadCommand_Synopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventReadCommand{
		Meta: testMeta(t, ui),
	}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Fatal("Synopsis should not be empty")
	}
}

func TestVCSEventReadCommand_Run_NoID(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &VCSEventReadCommand{
		Meta: testMeta(t, ui),
	}

	code := cmd.Run([]string{})
	if code == 0 {
		t.Fatal("Expected non-zero exit code when ID is missing")
	}

	output := ui.ErrorWriter.String()
	if !contains(output, "id") {
		t.Error("Error should mention missing id flag")
	}
}
