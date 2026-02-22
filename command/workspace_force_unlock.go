package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type workspaceForceUnlocker interface {
	ForceUnlock(ctx context.Context, workspaceID string) (*tfe.Workspace, error)
}

// WorkspaceForceUnlockCommand is a command to force-unlock a workspace.
type WorkspaceForceUnlockCommand struct {
	Meta
	organization string
	name         string
	format       string
	workspaceSvc workspaceReader
	forceSvc     workspaceForceUnlocker
}

// Run executes the workspace force-unlock command.
func (c *WorkspaceForceUnlockCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace force-unlock")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	workspace, err := c.workspaceService(client).Read(client.Context(), c.organization, c.name)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	unlockedWorkspace, err := c.forceService(client).ForceUnlock(client.Context(), workspace.ID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error force-unlocking workspace: %s", err))
		return 1
	}

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Workspace '%s' force-unlocked successfully", unlockedWorkspace.Name))
	}

	formatter := c.Meta.NewFormatter(c.format)
	formatter.KeyValue(map[string]interface{}{
		"ID":           unlockedWorkspace.ID,
		"Name":         unlockedWorkspace.Name,
		"Organization": c.organization,
		"Locked":       unlockedWorkspace.Locked,
	})

	return 0
}

func (c *WorkspaceForceUnlockCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *WorkspaceForceUnlockCommand) forceService(client *client.Client) workspaceForceUnlocker {
	if c.forceSvc != nil {
		return c.forceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace force-unlock command.
func (c *WorkspaceForceUnlockCommand) Help() string {
	helpText := `
Usage: hcptf workspace force-unlock [options]

  Force-unlock a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace force-unlock -org=my-org -name=my-workspace
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace force-unlock command.
func (c *WorkspaceForceUnlockCommand) Synopsis() string {
	return "Force-unlock a workspace"
}
