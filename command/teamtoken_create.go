package command

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
)

// TeamTokenCreateCommand is a command to create a team token
type TeamTokenCreateCommand struct {
	Meta
	teamID      string
	description string
	expiredAt   string
	format      string
}

// Run executes the team token create command
func (c *TeamTokenCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamtoken create")
	flags.StringVar(&c.teamID, "team-id", "", "Team ID (required)")
	flags.StringVar(&c.description, "description", "", "Token description (required)")
	flags.StringVar(&c.expiredAt, "expired-at", "", "Expiration date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.teamID == "" {
		c.Ui.Error("Error: -team-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.description == "" {
		c.Ui.Error("Error: -description flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.TeamTokenCreateOptions{
		Description: tfe.String(c.description),
	}

	// Parse expiration date if provided
	if c.expiredAt != "" {
		t, err := time.Parse(time.RFC3339, c.expiredAt)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error parsing expired-at date: %s", err))
			c.Ui.Error("Use ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
			return 1
		}
		options.ExpiredAt = &t
	}

	// Create team token
	token, err := client.TeamTokens.CreateWithOptions(client.Context(), c.teamID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating team token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output("Team token created successfully")

	// Show token details including the secret value
	data := map[string]interface{}{
		"ID":          token.ID,
		"Description": token.Description,
		"Token":       token.Token,
		"CreatedAt":   token.CreatedAt,
	}

	if token.Team != nil {
		data["TeamID"] = token.Team.ID
	}

	if token.CreatedBy != nil && token.CreatedBy.User != nil {
		data["CreatedBy"] = token.CreatedBy.User.Username
	}

	if !token.ExpiredAt.IsZero() {
		data["ExpiredAt"] = token.ExpiredAt.Format("2006-01-02 15:04:05")
	} else {
		data["ExpiredAt"] = "Never"
	}

	formatter.KeyValue(data)

	// Warning about token visibility
	c.Ui.Warn("\nWARNING: This is the only time the token will be displayed. Save it securely.")

	return 0
}

// Help returns help text for the team token create command
func (c *TeamTokenCreateCommand) Help() string {
	helpText := `
Usage: hcptf teamtoken create [options]

  Create a team API token.

  Team tokens authenticate as the team and have access to all workspaces
  the team has permissions to access. Teams can have multiple tokens with
  different descriptions to track usage across different applications.

Options:

  -team-id=<id>         Team ID (required)
  -description=<text>   Token description (required, must be unique per team)
  -expired-at=<date>    Token expiration date in ISO 8601 format
                        (e.g., 2024-12-31T23:59:59Z). If omitted, token never expires.
  -output=<format>      Output format: table (default) or json

Example:

  # Create team token that never expires
  hcptf teamtoken create -team-id=team-abc123 -description="CI/CD Pipeline"

  # Create team token with expiration
  hcptf teamtoken create -team-id=team-abc123 \
    -description="Temporary Access" -expired-at=2024-12-31T23:59:59Z

Security Note:

  The token value is only displayed once upon creation and cannot be retrieved
  later. Store it securely. Each description must be unique within the team.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team token create command
func (c *TeamTokenCreateCommand) Synopsis() string {
	return "Create a team API token"
}
