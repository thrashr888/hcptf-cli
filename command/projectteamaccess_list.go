package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type ProjectTeamAccessListCommand struct {
	Meta
	projectID              string
	format                 string
	projectTeamAccessSvc   projectTeamAccessLister
}

// Run executes the project team access list command
func (c *ProjectTeamAccessListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("projectteamaccess list")
	flags.StringVar(&c.projectID, "project-id", "", "Project ID (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List project team access
	projectTeamAccessList, err := c.projectTeamAccessService(client).List(client.Context(), tfe.TeamProjectAccessListOptions{
		ProjectID: c.projectID,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing project team access: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(projectTeamAccessList.Items) == 0 {
		c.Ui.Output("No project team access found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Team ID", "Access Level"}
	var rows [][]string

	for _, pta := range projectTeamAccessList.Items {
		accessLevel := string(pta.Access)
		if pta.Access == tfe.TeamProjectAccessCustom {
			accessLevel = "custom"
		}

		rows = append(rows, []string{
			pta.ID,
			pta.Team.ID,
			accessLevel,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *ProjectTeamAccessListCommand) projectTeamAccessService(client *client.Client) projectTeamAccessLister {
	if c.projectTeamAccessSvc != nil {
		return c.projectTeamAccessSvc
	}
	return client.TeamProjectAccess
}

// Help returns help text for the project team access list command
func (c *ProjectTeamAccessListCommand) Help() string {
	helpText := `
Usage: hcptf projectteamaccess list [options]

  List team access for a project.

Options:

  -project-id=<id>  Project ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf projectteamaccess list -project-id=prj-123abc
  hcptf projectteamaccess list -project-id=prj-123abc -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the project team access list command
func (c *ProjectTeamAccessListCommand) Synopsis() string {
	return "List team access for a project"
}
