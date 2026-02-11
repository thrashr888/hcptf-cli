package command

import (
	"fmt"
	"strings"
)

// RunTaskDetachCommand is a command to detach a run task from a workspace
type RunTaskDetachCommand struct {
	Meta
	organization       string
	workspace          string
	workspaceRunTaskID string
	force              bool
}

// Run executes the run task detach command
func (c *RunTaskDetachCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("runtask detach")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.workspaceRunTaskID, "workspace-runtask-id", "", "Workspace run task ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force detach without confirmation")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspaceRunTaskID == "" {
		c.Ui.Error("Error: -workspace-runtask-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace to obtain its ID
	workspace, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Read workspace run task to get details for confirmation
	workspaceRunTask, err := client.WorkspaceRunTasks.Read(client.Context(), workspace.ID, c.workspaceRunTaskID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace run task: %s", err))
		return 1
	}

	runTaskName := "unknown"
	if workspaceRunTask.RunTask != nil {
		runTaskName = workspaceRunTask.RunTask.Name
	}

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to detach run task '%s' from workspace '%s'? (yes/no): ", runTaskName, c.workspace))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Detach cancelled")
			return 0
		}
	}

	// Detach run task from workspace
	err = client.WorkspaceRunTasks.Delete(client.Context(), workspace.ID, c.workspaceRunTaskID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error detaching run task from workspace: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Run task '%s' detached from workspace '%s' successfully", runTaskName, c.workspace))
	return 0
}

// Help returns help text for the run task detach command
func (c *RunTaskDetachCommand) Help() string {
	helpText := `
Usage: hcptf runtask detach [options]

  Detach a run task from a workspace. This removes the workspace-task association
  but does not delete the run task itself.

  To find the workspace-runtask-id, you can list workspace run tasks using the
  Terraform Cloud API or UI.

Options:

  -organization=<name>        Organization name (required)
  -org=<name>                Alias for -organization
  -workspace=<name>          Workspace name (required)
  -workspace-runtask-id=<id> Workspace run task ID (required)
  -force                     Force detach without confirmation

Example:

  hcptf runtask detach -org=my-org -workspace=prod \
    -workspace-runtask-id=wr-ABC123

  hcptf runtask detach -org=my-org -workspace=dev \
    -workspace-runtask-id=wr-XYZ789 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the run task detach command
func (c *RunTaskDetachCommand) Synopsis() string {
	return "Detach a run task from a workspace"
}
