package command

import (
	"strings"
	"testing"

	"github.com/mitchellh/cli"
)

func TestWorkspaceContextCommand(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceContextCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{"-org=testorg", "-workspace=testworkspace"})
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}

	output := ui.OutputWriter.String()
	if !strings.Contains(output, "testorg/testworkspace") {
		t.Errorf("expected output to contain workspace path, got: %s", output)
	}
	if !strings.Contains(output, "Runs:") {
		t.Errorf("expected output to contain Runs section, got: %s", output)
	}
	if !strings.Contains(output, "Variables:") {
		t.Errorf("expected output to contain Variables section, got: %s", output)
	}
	if !strings.Contains(output, "State:") {
		t.Errorf("expected output to contain State section, got: %s", output)
	}
}

func TestWorkspaceContextCommandNoOrgOrWorkspace(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceContextCommand{
		Meta: newTestMeta(ui),
	}

	code := cmd.Run([]string{})
	if code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}

	errOutput := ui.ErrorWriter.String()
	if !strings.Contains(errOutput, "No organization or workspace specified") {
		t.Errorf("expected error about missing org/workspace, got: %s", errOutput)
	}
}

func TestWorkspaceContextHelp(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceContextCommand{
		Meta: newTestMeta(ui),
	}

	help := cmd.Help()
	if help == "" {
		t.Error("expected non-empty help text")
	}
}

func TestWorkspaceContextSynopsis(t *testing.T) {
	ui := cli.NewMockUi()
	cmd := &WorkspaceContextCommand{
		Meta: newTestMeta(ui),
	}

	synopsis := cmd.Synopsis()
	if synopsis == "" {
		t.Error("expected non-empty synopsis")
	}
}
