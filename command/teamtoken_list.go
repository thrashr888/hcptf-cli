package command

import (
	"fmt"
	"strings"
)

// TeamTokenListCommand is a command to list team tokens
type TeamTokenListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the team token list command
func (c *TeamTokenListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("teamtoken list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List team tokens
	tokens, err := client.TeamTokens.List(client.Context(), c.organization, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing team tokens: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(tokens.Items) == 0 {
		c.Ui.Output("No team tokens found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Team ID", "Description", "Created At", "Last Used At", "Expires At"}
	var rows [][]string

	for _, token := range tokens.Items {
		lastUsed := "Never"
		if !token.LastUsedAt.IsZero() {
			lastUsed = token.LastUsedAt.Format("2006-01-02 15:04:05")
		}

		expiresAt := "Never"
		if !token.ExpiredAt.IsZero() {
			expiresAt = token.ExpiredAt.Format("2006-01-02 15:04:05")
		}

		description := "-"
		if token.Description != nil && *token.Description != "" {
			description = *token.Description
		}

		teamID := ""
		if token.Team != nil {
			teamID = token.Team.ID
		}

		rows = append(rows, []string{
			token.ID,
			teamID,
			description,
			token.CreatedAt.Format("2006-01-02 15:04:05"),
			lastUsed,
			expiresAt,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the team token list command
func (c *TeamTokenListCommand) Help() string {
	helpText := `
Usage: hcptf teamtoken list [options]

  List team API tokens for an organization.

  This command displays all team tokens in the specified organization. Note
  that the actual token values are only displayed when tokens are created
  and cannot be retrieved later.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf teamtoken list -org=my-org
  hcptf teamtoken list -organization=my-org -output=json

Note:

  Team tokens authenticate as the team and have access to all workspaces
  the team has permissions to access.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the team token list command
func (c *TeamTokenListCommand) Synopsis() string {
	return "List team API tokens"
}
