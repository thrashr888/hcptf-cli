package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// SSHKeyCreateCommand is a command to create an SSH key
type SSHKeyCreateCommand struct {
	Meta
	organization string
	name         string
	value        string
	format       string
	sshKeySvc    sshKeyCreator
}

// Run executes the SSH key create command
func (c *SSHKeyCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("sshkey create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "SSH key name (required)")
	flags.StringVar(&c.value, "value", "", "SSH private key content (required)")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.value == "" {
		c.Ui.Error("Error: -value flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create SSH key
	options := tfe.SSHKeyCreateOptions{
		Name:  tfe.String(c.name),
		Value: tfe.String(c.value),
	}

	sshKey, err := c.sshKeyService(client).Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating SSH key: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("SSH key '%s' created successfully", sshKey.Name))

	// Show SSH key details
	data := map[string]interface{}{
		"ID":   sshKey.ID,
		"Name": sshKey.Name,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *SSHKeyCreateCommand) sshKeyService(client *client.Client) sshKeyCreator {
	if c.sshKeySvc != nil {
		return c.sshKeySvc
	}
	return client.SSHKeys
}

// Help returns help text for the SSH key create command
func (c *SSHKeyCreateCommand) Help() string {
	helpText := `
Usage: hcptf sshkey create [options]

  Create a new SSH key for an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         SSH key name (required)
  -value=<key>         SSH private key content (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf sshkey create -org=my-org -name=my-ssh-key -value="$(cat ~/.ssh/id_rsa)"
  hcptf sshkey create -org=my-org -name=deploy-key -value="-----BEGIN RSA PRIVATE KEY-----\n..."
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the SSH key create command
func (c *SSHKeyCreateCommand) Synopsis() string {
	return "Create a new SSH key"
}
