package command

import (
	"fmt"
	"strings"

)

// SSHKeyReadCommand is a command to read SSH key details
type SSHKeyReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the SSH key read command
func (c *SSHKeyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("sshkey read")
	flags.StringVar(&c.id, "id", "", "SSH key ID (required)")
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

	// Read SSH key
	sshKey, err := client.SSHKeys.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading SSH key: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":   sshKey.ID,
		"Name": sshKey.Name,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the SSH key read command
func (c *SSHKeyReadCommand) Help() string {
	helpText := `
Usage: hcptf sshkey read [options]

  Read SSH key details.

Options:

  -id=<id>             SSH key ID (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf sshkey read -id=sshkey-123abc
  hcptf sshkey read -id=sshkey-456def -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the SSH key read command
func (c *SSHKeyReadCommand) Synopsis() string {
	return "Read SSH key details"
}
