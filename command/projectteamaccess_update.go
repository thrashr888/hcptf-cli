package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type ProjectTeamAccessUpdateCommand struct {
	Meta
	id                     string
	access                 string
	format                 string
	projectTeamAccessSvc   projectTeamAccessUpdater
}

// Run executes the project team access update command
func (c *ProjectTeamAccessUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("projectteamaccess update")
	flags.StringVar(&c.id, "id", "", "Project team access ID (required)")
	flags.StringVar(&c.access, "access", "", "Access level: read, write, maintain, admin, or custom (required)")
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

	if c.access == "" {
		c.Ui.Error("Error: -access flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate access level
	var accessLevel tfe.TeamProjectAccessType
	switch strings.ToLower(c.access) {
	case "read":
		accessLevel = tfe.TeamProjectAccessRead
	case "write":
		accessLevel = tfe.TeamProjectAccessWrite
	case "maintain":
		accessLevel = tfe.TeamProjectAccessMaintain
	case "admin":
		accessLevel = tfe.TeamProjectAccessAdmin
	case "custom":
		accessLevel = tfe.TeamProjectAccessCustom
	default:
		c.Ui.Error(fmt.Sprintf("Error: invalid access level '%s'. Must be one of: read, write, maintain, admin, custom", c.access))
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Update project team access
	options := tfe.TeamProjectAccessUpdateOptions{
		Access: &accessLevel,
	}

	projectTeamAccess, err := c.projectTeamAccessService(client).Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating project team access: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Project team access updated successfully"))

	// Show project team access details
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

func (c *ProjectTeamAccessUpdateCommand) projectTeamAccessService(client *client.Client) projectTeamAccessUpdater {
	if c.projectTeamAccessSvc != nil {
		return c.projectTeamAccessSvc
	}
	return client.TeamProjectAccess
}

// Help returns help text for the project team access update command
func (c *ProjectTeamAccessUpdateCommand) Help() string {
	helpText := `
Usage: hcptf projectteamaccess update [options]

  Update team project permissions.

Options:

  -id=<id>         Project team access ID (required)
  -access=<level>  Access level: read, write, maintain, admin, or custom (required)
  -output=<format> Output format: table (default) or json

Access Levels:
  read     - Read project and Read workspace access on project workspaces
  write    - Read project and Write workspace access on project workspaces
  maintain - Read project and Admin workspace access on project workspaces
  admin    - Admin project, Admin workspace access, create/move workspaces, manage team access
  custom   - Custom permissions

Example:

  hcptf projectteamaccess update -id=tprj-123abc -access=write
  hcptf projectteamaccess update -id=tprj-123abc -access=admin -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project team access update command
func (c *ProjectTeamAccessUpdateCommand) Synopsis() string {
	return "Update team project permissions"
}
