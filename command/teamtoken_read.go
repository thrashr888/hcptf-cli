package command

import (
	"fmt"
	"strings"
)

// TeamTokenReadCommand is a command to read a team token
type TeamTokenReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the team token read command
func (c *TeamTokenReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamtoken read")
	flags.StringVar(&c.id, "id", "", "Team token ID (required)")
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

	// Read team token
	token, err := client.TeamTokens.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading team token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	// Show token details (without the secret)
	data := map[string]interface{}{
		"ID":          token.ID,
		"Description": token.Description,
		"CreatedAt":   token.CreatedAt,
	}

	if token.Team != nil {
		data["TeamID"] = token.Team.ID
	}

	if token.CreatedBy != nil && token.CreatedBy.User != nil {
		data["CreatedBy"] = token.CreatedBy.User.Username
	}

	if !token.LastUsedAt.IsZero() {
		data["LastUsedAt"] = token.LastUsedAt.Format("2006-01-02 15:04:05")
	} else {
		data["LastUsedAt"] = "Never"
	}

	if !token.ExpiredAt.IsZero() {
		data["ExpiredAt"] = token.ExpiredAt.Format("2006-01-02 15:04:05")
	} else {
		data["ExpiredAt"] = "Never"
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the team token read command
func (c *TeamTokenReadCommand) Help() string {
	helpText := `
Usage: hcptf teamtoken read [options]

  Show team token details.

  This command displays metadata about a team token. Note that the actual
  token value is only displayed when the token is created and cannot be
  retrieved later.

Options:

  -id=<id>          Team token ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf teamtoken read -id=at-abc123xyz
  hcptf teamtoken read -id=at-abc123xyz -output=json

Note:

  Team tokens authenticate as the team and have access to all workspaces
  the team has permissions to access.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team token read command
func (c *TeamTokenReadCommand) Synopsis() string {
	return "Show team token details"
}
