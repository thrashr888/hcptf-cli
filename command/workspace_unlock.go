package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type workspaceUnlocker interface {
	Unlock(ctx context.Context, workspaceID string) (*tfe.Workspace, error)
}

// WorkspaceUnlockCommand is a command to unlock a workspace.
type WorkspaceUnlockCommand struct {
	Meta
	organization string
	name         string
	format       string
	workspaceSvc workspaceReader
	unlockSvc    workspaceUnlocker
}

// Run executes the workspace unlock command.
func (c *WorkspaceUnlockCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace unlock")
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

	unlockedWorkspace, err := c.unlockService(client).Unlock(client.Context(), workspace.ID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error unlocking workspace: %s", err))
		return 1
	}

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Workspace '%s' unlocked successfully", unlockedWorkspace.Name))
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

func (c *WorkspaceUnlockCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *WorkspaceUnlockCommand) unlockService(client *client.Client) workspaceUnlocker {
	if c.unlockSvc != nil {
		return c.unlockSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace unlock command.
func (c *WorkspaceUnlockCommand) Help() string {
	helpText := `
Usage: hcptf workspace unlock [options]

  Unlock a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace unlock -org=my-org -name=my-workspace
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace unlock command.
func (c *WorkspaceUnlockCommand) Synopsis() string {
	return "Unlock a workspace"
}
