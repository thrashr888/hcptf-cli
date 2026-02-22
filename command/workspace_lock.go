package command

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type workspaceLocker interface {
	Lock(ctx context.Context, workspaceID string, options tfe.WorkspaceLockOptions) (*tfe.Workspace, error)
}

// WorkspaceLockCommand is a command to lock a workspace.
type WorkspaceLockCommand struct {
	Meta
	organization string
	name         string
	reason       string
	format       string
	workspaceSvc workspaceReader
	lockSvc      workspaceLocker
}

// Run executes the workspace lock command.
func (c *WorkspaceLockCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace lock")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")
	flags.StringVar(&c.reason, "reason", "", "Reason for locking the workspace")
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

	options := tfe.WorkspaceLockOptions{}
	if c.reason != "" {
		options.Reason = tfe.String(c.reason)
	}

	lockedWorkspace, err := c.lockService(client).Lock(client.Context(), workspace.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error locking workspace: %s", err))
		return 1
	}

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Workspace '%s' locked successfully", lockedWorkspace.Name))
	}

	formatter := c.Meta.NewFormatter(c.format)
	formatter.KeyValue(map[string]interface{}{
		"ID":           lockedWorkspace.ID,
		"Name":         lockedWorkspace.Name,
		"Organization": c.organization,
		"Locked":       lockedWorkspace.Locked,
		"LockReason":   c.reason,
	})

	return 0
}

func (c *WorkspaceLockCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *WorkspaceLockCommand) lockService(client *client.Client) workspaceLocker {
	if c.lockSvc != nil {
		return c.lockSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace lock command.
func (c *WorkspaceLockCommand) Help() string {
	helpText := `
Usage: hcptf workspace lock [options]

  Lock a workspace to prevent new runs.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -reason=<text>       Optional reason for the lock
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace lock -org=my-org -name=my-workspace
  hcptf workspace lock -org=my-org -name=my-workspace -reason="maintenance window"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace lock command.
func (c *WorkspaceLockCommand) Synopsis() string {
	return "Lock a workspace"
}
