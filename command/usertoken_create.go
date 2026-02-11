package command

import (
	"fmt"
	"strings"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// UserTokenCreateCommand is a command to create a user token
type UserTokenCreateCommand struct {
	Meta
	description string
	expiredAt   string
	format      string
}

// Run executes the user token create command
func (c *UserTokenCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("usertoken create")
	flags.StringVar(&c.description, "description", "", "Token description (required)")
	flags.StringVar(&c.expiredAt, "expired-at", "", "Expiration date in ISO 8601 format (e.g., 2024-12-31T23:59:59Z)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
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

	// Get current user
	user, err := client.Users.ReadCurrent(client.Context())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading current user: %s", err))
		return 1
	}

	// Build create options
	options := tfe.UserTokenCreateOptions{
		Description: c.description,
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

	// Create user token
	token, err := client.UserTokens.Create(client.Context(), user.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating user token: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output("User token created successfully")

	// Show token details including the secret value
	data := map[string]interface{}{
		"ID":          token.ID,
		"Description": token.Description,
		"Token":       token.Token,
		"CreatedAt":   token.CreatedAt,
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

// Help returns help text for the user token create command
func (c *UserTokenCreateCommand) Help() string {
	helpText := `
Usage: hcptf usertoken create [options]

  Create a user API token.

  User tokens authenticate as the user account and have the same permissions
  as the user. You can create multiple user tokens with different descriptions
  to track usage across different applications.

Options:

  -description=<text>   Token description (required)
  -expired-at=<date>    Token expiration date in ISO 8601 format
                        (e.g., 2024-12-31T23:59:59Z). If omitted, token never expires.
  -output=<format>      Output format: table (default) or json

Example:

  # Create user token that never expires
  hcptf usertoken create -description="CI/CD Pipeline"

  # Create user token with expiration
  hcptf usertoken create -description="Temporary Access" -expired-at=2024-12-31T23:59:59Z

Security Note:

  The token value is only displayed once upon creation and cannot be retrieved
  later. Store it securely. This command creates tokens for the currently
  authenticated user only.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the user token create command
func (c *UserTokenCreateCommand) Synopsis() string {
	return "Create a user API token"
}
