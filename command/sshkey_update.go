package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// SSHKeyUpdateCommand is a command to update an SSH key
type SSHKeyUpdateCommand struct {
	Meta
	id     string
	name   string
	format string
}

// Run executes the SSH key update command
func (c *SSHKeyUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("sshkey update")
	flags.StringVar(&c.id, "id", "", "SSH key ID (required)")
	flags.StringVar(&c.name, "name", "", "SSH key name")
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

	// Build update options
	options := tfe.SSHKeyUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	// Update SSH key
	sshKey, err := client.SSHKeys.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating SSH key: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("SSH key '%s' updated successfully", sshKey.Name))

	// Show SSH key details
	data := map[string]interface{}{
		"ID":   sshKey.ID,
		"Name": sshKey.Name,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the SSH key update command
func (c *SSHKeyUpdateCommand) Help() string {
	helpText := `
Usage: hcptf sshkey update [options]

  Update an SSH key.

Options:

  -id=<id>             SSH key ID (required)
  -name=<name>         SSH key name
  -output=<format>     Output format: table (default) or json

Example:

  hcptf sshkey update -id=sshkey-123abc -name=new-ssh-key-name
  hcptf sshkey update -id=sshkey-456def -name=updated-key -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the SSH key update command
func (c *SSHKeyUpdateCommand) Synopsis() string {
	return "Update an SSH key"
}
