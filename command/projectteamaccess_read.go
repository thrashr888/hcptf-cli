package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type ProjectTeamAccessReadCommand struct {
	Meta
	id                   string
	format               string
	projectTeamAccessSvc projectTeamAccessReader
}

// Run executes the project team access read command
func (c *ProjectTeamAccessReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("projectteamaccess read")
	flags.StringVar(&c.id, "id", "", "Project team access ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read project team access
	projectTeamAccess, err := c.projectTeamAccessService(client).Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading project team access: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":        projectTeamAccess.ID,
		"TeamID":    projectTeamAccess.Team.ID,
		"ProjectID": projectTeamAccess.Project.ID,
		"Access":    string(projectTeamAccess.Access),
	}

	// Add custom permissions if access is custom
	if projectTeamAccess.Access == tfe.TeamProjectAccessCustom {
		// Project access permissions
		if projectTeamAccess.ProjectAccess != nil {
			data["ProjectSettings"] = string(projectTeamAccess.ProjectAccess.ProjectSettingsPermission)
			data["ProjectTeams"] = string(projectTeamAccess.ProjectAccess.ProjectTeamsPermission)
		}
		// Workspace access permissions
		if projectTeamAccess.WorkspaceAccess != nil {
			data["WorkspaceRuns"] = string(projectTeamAccess.WorkspaceAccess.WorkspaceRunsPermission)
			data["WorkspaceVariables"] = string(projectTeamAccess.WorkspaceAccess.WorkspaceVariablesPermission)
			data["WorkspaceStateVersions"] = string(projectTeamAccess.WorkspaceAccess.WorkspaceStateVersionsPermission)
			data["WorkspaceSentinelMocks"] = string(projectTeamAccess.WorkspaceAccess.WorkspaceSentinelMocksPermission)
			data["WorkspaceCreate"] = projectTeamAccess.WorkspaceAccess.WorkspaceCreatePermission
			data["WorkspaceDelete"] = projectTeamAccess.WorkspaceAccess.WorkspaceDeletePermission
			data["WorkspaceMove"] = projectTeamAccess.WorkspaceAccess.WorkspaceMovePermission
			data["WorkspaceLocking"] = projectTeamAccess.WorkspaceAccess.WorkspaceLockingPermission
			data["WorkspaceRunTasks"] = projectTeamAccess.WorkspaceAccess.WorkspaceRunTasksPermission
		}
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ProjectTeamAccessReadCommand) projectTeamAccessService(client *client.Client) projectTeamAccessReader {
	if c.projectTeamAccessSvc != nil {
		return c.projectTeamAccessSvc
	}
	return client.TeamProjectAccess
}

// Help returns help text for the project team access read command
func (c *ProjectTeamAccessReadCommand) Help() string {
	helpText := `
Usage: hcptf projectteamaccess read [options]

  Show project team access details.

Options:

  -id=<id>         Project team access ID (required)
  -output=<format> Output format: table (default) or json

Example:

  hcptf projectteamaccess read -id=tprj-123abc
  hcptf projectteamaccess read -id=tprj-123abc -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project team access read command
func (c *ProjectTeamAccessReadCommand) Synopsis() string {
	return "Show project team access details"
}
