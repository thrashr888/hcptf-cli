package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// SSHKeyListCommand is a command to list SSH keys
type SSHKeyListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the SSH key list command
func (c *SSHKeyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("sshkey list")
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

	// List SSH keys
	sshKeys, err := client.SSHKeys.List(client.Context(), c.organization, &tfe.SSHKeyListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing SSH keys: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(sshKeys.Items) == 0 {
		c.Ui.Output("No SSH keys found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name"}
	var rows [][]string

	for _, key := range sshKeys.Items {
		rows = append(rows, []string{
			key.ID,
			key.Name,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the SSH key list command
func (c *SSHKeyListCommand) Help() string {
	helpText := `
Usage: hcptf sshkey list [options]

  List SSH keys for an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf sshkey list -org=my-org
  hcptf sshkey list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the SSH key list command
func (c *SSHKeyListCommand) Synopsis() string {
	return "List SSH keys for an organization"
}
