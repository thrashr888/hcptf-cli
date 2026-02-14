package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type ProjectTeamAccessCreateCommand struct {
	Meta
	projectID            string
	teamID               string
	access               string
	format               string
	projectTeamAccessSvc projectTeamAccessCreator
}

// Run executes the project team access create command
func (c *ProjectTeamAccessCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("projectteamaccess create")
	flags.StringVar(&c.projectID, "project-id", "", "Project ID (required)")
	flags.StringVar(&c.teamID, "team-id", "", "Team ID (required)")
	flags.StringVar(&c.access, "access", "", "Access level: read, write, maintain, admin, or custom (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.projectID == "" {
		c.Ui.Error("Error: -project-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.teamID == "" {
		c.Ui.Error("Error: -team-id flag is required")
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

	// Create project team access
	options := tfe.TeamProjectAccessAddOptions{
		Access: accessLevel,
		Team: &tfe.Team{
			ID: c.teamID,
		},
		Project: &tfe.Project{
			ID: c.projectID,
		},
	}

	projectTeamAccess, err := c.projectTeamAccessService(client).Add(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating project team access: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Project team access created successfully"))

	// Show project team access details
	data := map[string]interface{}{
		"ID":        projectTeamAccess.ID,
		"TeamID":    projectTeamAccess.Team.ID,
		"ProjectID": projectTeamAccess.Project.ID,
		"Access":    string(projectTeamAccess.Access),
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ProjectTeamAccessCreateCommand) projectTeamAccessService(client *client.Client) projectTeamAccessCreator {
	if c.projectTeamAccessSvc != nil {
		return c.projectTeamAccessSvc
	}
	return client.TeamProjectAccess
}

// Help returns help text for the project team access create command
func (c *ProjectTeamAccessCreateCommand) Help() string {
	helpText := `
Usage: hcptf projectteamaccess create [options]

  Grant team access to a project.

Options:

  -project-id=<id>  Project ID (required)
  -team-id=<id>     Team ID (required)
  -access=<level>   Access level: read, write, maintain, admin, or custom (required)
  -output=<format>  Output format: table (default) or json

Access Levels:
  read     - Read project and Read workspace access on project workspaces
  write    - Read project and Write workspace access on project workspaces
  maintain - Read project and Admin workspace access on project workspaces
  admin    - Admin project, Admin workspace access, create/move workspaces, manage team access
  custom   - Custom permissions (use update command to configure)

Example:

  hcptf projectteamaccess create -project-id=prj-123abc -team-id=team-456def -access=write
  hcptf projectteamaccess create -project-id=prj-123abc -team-id=team-456def -access=read -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project team access create command
func (c *ProjectTeamAccessCreateCommand) Synopsis() string {
	return "Grant team access to a project"
}
