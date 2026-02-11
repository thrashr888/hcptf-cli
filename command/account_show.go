package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AccountShowCommand is a command to show current account details
type AccountShowCommand struct {
	Meta
	format string
}

// Run executes the account show command
func (c *AccountShowCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("account show")
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

	// Get current account details (using Users API)
	account, err := client.Users.ReadCurrent(client.Context())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading account: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                    account.ID,
		"Email":                 account.Email,
		"Username":              account.Username,
		"AvatarURL":             account.AvatarURL,
		"TwoFactorEnabled":      account.TwoFactor != nil && account.TwoFactor.Enabled,
		"IsServiceAccount":      account.IsServiceAccount,
		"UnconfirmedEmail":      account.UnconfirmedEmail,
	}

	if account.TwoFactor != nil {
		data["TwoFactorVerified"] = account.TwoFactor.Verified
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the account show command
func (c *AccountShowCommand) Help() string {
	helpText := `
Usage: hcptf account show [options]

  Show current account details.

  This command displays information about the currently authenticated user
  account, including email, username, and security settings.

Options:

  -output=<format>  Output format: table (default) or json

Example:

  # Show current account
  hcptf account show

  # Show as JSON
  hcptf account show -output=json

Note:

  This command requires authentication. Use 'hcptf login' first if needed.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the account show command
func (c *AccountShowCommand) Synopsis() string {
	return "Show current account details"
}
