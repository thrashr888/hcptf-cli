package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

type TeamAccessReadCommand struct {
	Meta
	id            string
	format        string
	teamAccessSvc teamAccessReader
}

// Run executes the team access read command
func (c *TeamAccessReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamaccess read")
	flags.StringVar(&c.id, "id", "", "Team access ID (required)")
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

	// Read team access
	teamAccess, err := c.teamAccessService(client).Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading team access: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

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

func (c *TeamAccessReadCommand) teamAccessService(client *client.Client) teamAccessReader {
	if c.teamAccessSvc != nil {
		return c.teamAccessSvc
	}
	return client.TeamAccess
}

// Help returns help text for the team access read command
func (c *TeamAccessReadCommand) Help() string {
	helpText := `
Usage: hcptf teamaccess read [options]

  Show team access details.

Options:

  -id=<id>         Team access ID (required)
  -output=<format> Output format: table (default) or json

Example:

  hcptf teamaccess read -id=tws-123abc
  hcptf teamaccess read -id=tws-123abc -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team access read command
func (c *TeamAccessReadCommand) Synopsis() string {
	return "Show team access details"
}
