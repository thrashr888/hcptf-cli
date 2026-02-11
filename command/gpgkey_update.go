package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type GPGKeyUpdateCommand struct {
	Meta
	namespace    string
	keyID        string
	newNamespace string
	format       string
	gpgKeySvc    gpgKeyUpdater
}

// Run executes the GPG key update command
func (c *GPGKeyUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gpgkey update")
	flags.StringVar(&c.namespace, "namespace", "", "Current namespace (organization name) (required)")
	flags.StringVar(&c.keyID, "key-id", "", "GPG key ID (required)")
	flags.StringVar(&c.newNamespace, "new-namespace", "", "New namespace (organization name) (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.namespace == "" {
		c.Ui.Error("Error: -namespace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.keyID == "" {
		c.Ui.Error("Error: -key-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.newNamespace == "" {
		c.Ui.Error("Error: -new-namespace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Update GPG key
	keyID := tfe.GPGKeyID{
		RegistryName: tfe.PrivateRegistry,
		Namespace:    c.namespace,
		KeyID:        c.keyID,
	}

	options := tfe.GPGKeyUpdateOptions{
		Namespace: c.newNamespace,
	}

	key, err := c.gpgService(client).Update(client.Context(), keyID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating GPG key: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("GPG key namespace updated from '%s' to '%s'", c.namespace, c.newNamespace))

	// Show key details
	data := map[string]interface{}{
		"ID":        key.ID,
		"KeyID":     key.KeyID,
		"Namespace": key.Namespace,
		"UpdatedAt": key.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the GPG key update command
func (c *GPGKeyUpdateCommand) Help() string {
	helpText := `
Usage: hcptf gpgkey update [options]

  Update a GPG key's namespace (organization).

  Only the namespace can be updated. The namespace must be an
  organization you have permission to access.

Options:

  -namespace=<name>      Current namespace (organization name) (required)
  -key-id=<id>           GPG key ID (required)
  -new-namespace=<name>  New namespace (organization name) (required)
  -output=<format>       Output format: table (default) or json

Example:

  hcptf gpgkey update -namespace=old-org -key-id=32966F3FB5AC1129 -new-namespace=new-org
`
	return strings.TrimSpace(helpText)
}

func (c *GPGKeyUpdateCommand) gpgService(client *client.Client) gpgKeyUpdater {
	if c.gpgKeySvc != nil {
		return c.gpgKeySvc
	}
	return client.GPGKeys
}

// Synopsis returns a short synopsis for the GPG key update command
func (c *GPGKeyUpdateCommand) Synopsis() string {
	return "Update a GPG key's namespace"
}
