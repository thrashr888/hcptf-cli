package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestWorkspaceResourceReadCommand_implements(t *testing.T) {
	var _ cli.Command = &WorkspaceResourceReadCommand{}
}

func TestWorkspaceResourceReadCommand_Help(t *testing.T) {
	c := &WorkspaceResourceReadCommand{}
	help := c.Help()

	if help == "" {
		t.Fatal("help text should not be empty")
	}

	if len(help) < 50 {
		t.Fatal("help text should be descriptive")
	}
}

func TestWorkspaceResourceReadCommand_Synopsis(t *testing.T) {
	c := &WorkspaceResourceReadCommand{}
	synopsis := c.Synopsis()

	if synopsis == "" {
		t.Fatal("synopsis should not be empty")
	}

	expected := "Show resource details"
	if synopsis != expected {
		t.Errorf("Expected synopsis %q, got %q", expected, synopsis)
	}
}

func TestWorkspaceResourceReadCommand_Run_MissingWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	c := &WorkspaceResourceReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	// Run without workspace ID should fail
	code := c.Run([]string{})
	if code != 1 {
		t.Errorf("Expected exit code 1, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if output == "" {
		t.Fatal("Expected error output")
	}
}

func TestWorkspaceResourceReadCommand_Run_MissingResourceID(t *testing.T) {
	ui := cli.NewMockUi()
	c := &WorkspaceResourceReadCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	// Run with workspace ID but without resource ID should fail
	code := c.Run([]string{"-workspace-id=ws-test"})
	if code != 1 {
		t.Errorf("Expected exit code 1, got %d", code)
	}

	output := ui.ErrorWriter.String()
	if output == "" {
		t.Fatal("Expected error output")
	}
}
