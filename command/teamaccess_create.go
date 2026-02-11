package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type TeamAccessCreateCommand struct {
	Meta
	workspaceID   string
	teamID        string
	access        string
	format        string
	teamAccessSvc teamAccessCreator
}

// Run executes the team access create command
func (c *TeamAccessCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamaccess create")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.teamID, "team-id", "", "Team ID (required)")
	flags.StringVar(&c.access, "access", "", "Access level: read, plan, write, admin, or custom (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.workspaceID == "" {
		c.Ui.Error("Error: -workspace-id flag is required")
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
	var accessLevel tfe.AccessType
	switch strings.ToLower(c.access) {
	case "read":
		accessLevel = tfe.AccessRead
	case "plan":
		accessLevel = tfe.AccessPlan
	case "write":
		accessLevel = tfe.AccessWrite
	case "admin":
		accessLevel = tfe.AccessAdmin
	case "custom":
		accessLevel = tfe.AccessCustom
	default:
		c.Ui.Error(fmt.Sprintf("Error: invalid access level '%s'. Must be one of: read, plan, write, admin, custom", c.access))
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create team access
	options := tfe.TeamAccessAddOptions{
		Access: &accessLevel,
		Team: &tfe.Team{
			ID: c.teamID,
		},
		Workspace: &tfe.Workspace{
			ID: c.workspaceID,
		},
	}

	teamAccess, err := c.teamAccessService(client).Add(client.Context(), options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating team access: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Team access created successfully"))

	// Show team access details
	data := map[string]interface{}{
		"ID":          teamAccess.ID,
		"TeamID":      teamAccess.Team.ID,
		"WorkspaceID": teamAccess.Workspace.ID,
		"Access":      string(teamAccess.Access),
	}

	formatter.KeyValue(data)
	return 0
}

func (c *TeamAccessCreateCommand) teamAccessService(client *client.Client) teamAccessCreator {
	if c.teamAccessSvc != nil {
		return c.teamAccessSvc
	}
	return client.TeamAccess
}

// Help returns help text for the team access create command
func (c *TeamAccessCreateCommand) Help() string {
	helpText := `
Usage: hcptf teamaccess create [options]

  Grant team access to a workspace.

Options:

  -workspace-id=<id>  Workspace ID (required)
  -team-id=<id>       Team ID (required)
  -access=<level>     Access level: read, plan, write, admin, or custom (required)
  -output=<format>    Output format: table (default) or json

Access Levels:
  read   - Read-only access to workspace
  plan   - Can create plans
  write  - Can create plans and apply runs
  admin  - Full administrative access
  custom - Custom permissions (use update command to configure)

Example:

  hcptf teamaccess create -workspace-id=ws-123abc -team-id=team-456def -access=write
  hcptf teamaccess create -workspace-id=ws-123abc -team-id=team-456def -access=read -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team access create command
func (c *TeamAccessCreateCommand) Synopsis() string {
	return "Grant team access to a workspace"
}
