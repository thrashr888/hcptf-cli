package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestWorkspaceResourceListCommand_implements(t *testing.T) {
	var _ cli.Command = &WorkspaceResourceListCommand{}
}

func TestWorkspaceResourceListCommand_Help(t *testing.T) {
	c := &WorkspaceResourceListCommand{}
	help := c.Help()

	if help == "" {
		t.Fatal("help text should not be empty")
	}

	if len(help) < 50 {
		t.Fatal("help text should be descriptive")
	}
}

func TestWorkspaceResourceListCommand_Synopsis(t *testing.T) {
	c := &WorkspaceResourceListCommand{}
	synopsis := c.Synopsis()

	if synopsis == "" {
		t.Fatal("synopsis should not be empty")
	}

	expected := "List resources in workspace state"
	if synopsis != expected {
		t.Errorf("Expected synopsis %q, got %q", expected, synopsis)
	}
}

func TestWorkspaceResourceListCommand_Run_MissingWorkspaceID(t *testing.T) {
	ui := cli.NewMockUi()
	c := &WorkspaceResourceListCommand{
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
