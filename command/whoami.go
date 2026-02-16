package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// WhoAmICommand returns information about the authenticated user.
type WhoAmICommand struct {
	Meta
	format     string
	accountSvc accountReader
}

// Run executes the whoami command.
func (c *WhoAmICommand) Run(args []string) int {
	flags := c.Meta.FlagSet("whoami")
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

	// Get current user account details
	account, err := c.accountService(client).ReadCurrent(client.Context())
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading account: %s", err))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	data := map[string]interface{}{
		"ID":               account.ID,
		"Username":         account.Username,
		"Email":            account.Email,
		"IsServiceAccount": account.IsServiceAccount,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *WhoAmICommand) accountService(client *client.Client) accountReader {
	if c.accountSvc != nil {
		return c.accountSvc
	}
	return client.Users
}

// Help returns help text for the whoami command
func (c *WhoAmICommand) Help() string {
	helpText := `
Usage: hcptf whoami [options]

  Show the currently authenticated user.

  This command is useful for verifying that authentication is working and for
  confirming the active CLI identity.

Options:

  -output=<format>  Output format: table (default) or json

Example:

  hcptf whoami
  hcptf whoami -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the whoami command
func (c *WhoAmICommand) Synopsis() string {
	return "Show the current authenticated user"
}
