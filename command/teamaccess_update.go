package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type TeamAccessUpdateCommand struct {
	Meta
	id            string
	access        string
	format        string
	teamAccessSvc teamAccessUpdater
}

// Run executes the team access update command
func (c *TeamAccessUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamaccess update")
	flags.StringVar(&c.id, "id", "", "Team access ID (required)")
	flags.StringVar(&c.access, "access", "", "Access level: read, plan, write, admin, or custom (required)")
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

	// Update team access
	options := tfe.TeamAccessUpdateOptions{
		Access: &accessLevel,
	}

	teamAccess, err := c.teamAccessService(client).Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating team access: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Team access updated successfully"))

	// Show team access details
	data := map[string]interface{}{
		"ID":          teamAccess.ID,
		"TeamID":      teamAccess.Team.ID,
		"WorkspaceID": teamAccess.Workspace.ID,
		"Access":      string(teamAccess.Access),
	}

	// Add custom permissions if access is custom
	if teamAccess.Access == "custom" {
		data["Runs"] = string(teamAccess.Runs)
		data["Variables"] = string(teamAccess.Variables)
		data["StateVersions"] = string(teamAccess.StateVersions)
		data["SentinelMocks"] = string(teamAccess.SentinelMocks)
		data["WorkspaceLocking"] = teamAccess.WorkspaceLocking
		data["RunTasks"] = teamAccess.RunTasks
	}

	formatter.KeyValue(data)
	return 0
}

func (c *TeamAccessUpdateCommand) teamAccessService(client *client.Client) teamAccessUpdater {
	if c.teamAccessSvc != nil {
		return c.teamAccessSvc
	}
	return client.TeamAccess
}

// Help returns help text for the team access update command
func (c *TeamAccessUpdateCommand) Help() string {
	helpText := `
Usage: hcptf teamaccess update [options]

  Update team workspace permissions.

Options:

  -id=<id>         Team access ID (required)
  -access=<level>  Access level: read, plan, write, admin, or custom (required)
  -output=<format> Output format: table (default) or json

Access Levels:
  read   - Read-only access to workspace
  plan   - Can create plans
  write  - Can create plans and apply runs
  admin  - Full administrative access
  custom - Custom permissions

Example:

  hcptf teamaccess update -id=tws-123abc -access=write
  hcptf teamaccess update -id=tws-123abc -access=admin -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team access update command
func (c *TeamAccessUpdateCommand) Synopsis() string {
	return "Update team workspace permissions"
}
