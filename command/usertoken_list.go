package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// UserTokenListCommand is a command to list user tokens
type UserTokenListCommand struct {
	Meta
	format       string
	userSvc      userReader
	userTokenSvc userTokenLister
}

// Run executes the user token list command
func (c *UserTokenListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("usertoken list")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get current user
	user, err := c.userService(client).ReadCurrent(client.Context())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading current user: %s", err))
		return 1
	}

	// List user tokens
	tokens, err := c.userTokenService(client).List(client.Context(), user.ID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing user tokens: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(tokens.Items) == 0 {
		c.Ui.Output("No user tokens found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Description", "Created At", "Last Used At", "Expires At"}
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

		description := token.Description
		if description == "" {
			description = "-"
		}

		rows = append(rows, []string{
			token.ID,
			description,
			token.CreatedAt.Format("2006-01-02 15:04:05"),
			lastUsed,
			expiresAt,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *UserTokenListCommand) userService(client *client.Client) userReader {
	if c.userSvc != nil {
		return c.userSvc
	}
	return client.Users
}

func (c *UserTokenListCommand) userTokenService(client *client.Client) userTokenLister {
	if c.userTokenSvc != nil {
		return c.userTokenSvc
	}
	return client.UserTokens
}

// Help returns help text for the user token list command
func (c *UserTokenListCommand) Help() string {
	helpText := `
Usage: hcptf usertoken list [options]

  List user API tokens for the current user.

  This command displays all API tokens associated with the authenticated user.
  Note that the actual token values are only displayed when tokens are created
  and cannot be retrieved later.

Options:

  -output=<format>  Output format: table (default) or json

Example:

  hcptf usertoken list
  hcptf usertoken list -output=json

Note:

  This command lists tokens for the currently authenticated user only.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the user token list command
func (c *UserTokenListCommand) Synopsis() string {
	return "List user API tokens"
}
