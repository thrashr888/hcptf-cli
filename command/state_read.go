package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// StateReadCommand is a command to read state version details
type StateReadCommand struct {
	Meta
	stateVersionID string
	format         string
}

// Run executes the state read command
func (c *StateReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("state read")
	flags.StringVar(&c.stateVersionID, "id", "", "State version ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.stateVersionID == "" {
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

	// Read state version
	stateVersion, err := client.StateVersions.ReadCurrent(client.Context(), c.stateVersionID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading state version: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                 stateVersion.ID,
		"Serial":             stateVersion.Serial,
		"CreatedAt":          stateVersion.CreatedAt,
		"DownloadURL":        stateVersion.DownloadURL,
		"Resources":          stateVersion.Resources,
		"ResourcesProcessed": stateVersion.ResourcesProcessed,
		"StateVersion":       stateVersion.StateVersion,
		"TerraformVersion":   stateVersion.TerraformVersion,
		"VCSCommitSHA":       stateVersion.VCSCommitSHA,
		"VCSCommitURL":       stateVersion.VCSCommitURL,
	}

	if stateVersion.Run != nil {
		data["RunID"] = stateVersion.Run.ID
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the state read command
func (c *StateReadCommand) Help() string {
	helpText := `
Usage: hcptf state read [options]

  Read state version details.

Options:

  -id=<state-version-id>  State version ID (required)
  -output=<format>        Output format: table (default) or json

Example:

  hcptf state read -id=sv-abc123
  hcptf state read -id=sv-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the state read command
func (c *StateReadCommand) Synopsis() string {
	return "Read state version details"
}
