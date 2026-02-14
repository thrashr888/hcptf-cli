package command

import (
	"fmt"
	"strings"
)

// UserTokenReadCommand is a command to read a user token
type UserTokenReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the user token read command
func (c *UserTokenReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("usertoken read")
	flags.StringVar(&c.id, "id", "", "User token ID (required)")
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

	// Read user token
	token, err := client.UserTokens.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading user token: %s", err))
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

// Help returns help text for the user token read command
func (c *UserTokenReadCommand) Help() string {
	helpText := `
Usage: hcptf usertoken read [options]

  Show user token details.

  This command displays metadata about a user token. Note that the actual
  token value is only displayed when the token is created and cannot be
  retrieved later.

Options:

  -id=<id>          User token ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf usertoken read -id=at-abc123xyz
  hcptf usertoken read -id=at-abc123xyz -output=json

Note:

  This command can only read tokens for the currently authenticated user.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the user token read command
func (c *UserTokenReadCommand) Synopsis() string {
	return "Show user token details"
}
