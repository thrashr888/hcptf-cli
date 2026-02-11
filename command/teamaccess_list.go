package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type TeamAccessListCommand struct {
	Meta
	workspaceID   string
	format        string
	teamAccessSvc teamAccessLister
}

// Run executes the team access list command
func (c *TeamAccessListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamaccess list")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List team access
	teamAccessList, err := c.teamAccessService(client).List(client.Context(), &tfe.TeamAccessListOptions{
		WorkspaceID: c.workspaceID,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing team access: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(teamAccessList.Items) == 0 {
		c.Ui.Output("No team access found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Team ID", "Access Level"}
	var rows [][]string

	for _, ta := range teamAccessList.Items {
		accessLevel := string(ta.Access)
		if ta.Access == tfe.AccessCustom {
			// Show that it's custom
			accessLevel = "custom"
		}

		rows = append(rows, []string{
			ta.ID,
			ta.Team.ID,
			accessLevel,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *TeamAccessListCommand) teamAccessService(client *client.Client) teamAccessLister {
	if c.teamAccessSvc != nil {
		return c.teamAccessSvc
	}
	return client.TeamAccess
}

// Help returns help text for the team access list command
func (c *TeamAccessListCommand) Help() string {
	helpText := `
Usage: hcptf teamaccess list [options]

  List team access for a workspace.

Options:

  -workspace-id=<id>  Workspace ID (required)
  -output=<format>    Output format: table (default) or json

Example:

  hcptf teamaccess list -workspace-id=ws-123abc
  hcptf teamaccess list -workspace-id=ws-123abc -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team access list command
func (c *TeamAccessListCommand) Synopsis() string {
	return "List team access for a workspace"
}
