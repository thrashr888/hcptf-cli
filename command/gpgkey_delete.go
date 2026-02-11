package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type GPGKeyDeleteCommand struct {
	Meta
	namespace string
	keyID     string
	gpgKeySvc gpgKeyDeleter
}

// Run executes the GPG key delete command
func (c *GPGKeyDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gpgkey delete")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (organization name) (required)")
	flags.StringVar(&c.keyID, "key-id", "", "GPG key ID (required)")

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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Delete GPG key
	keyID := tfe.GPGKeyID{
		RegistryName: tfe.PrivateRegistry,
		Namespace:    c.namespace,
		KeyID:        c.keyID,
	}

	err = c.gpgService(client).Delete(client.Context(), keyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting GPG key: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("GPG key '%s' deleted successfully from namespace '%s'", c.keyID, c.namespace))
	return 0
}

// Help returns help text for the GPG key delete command
func (c *GPGKeyDeleteCommand) Help() string {
	helpText := `
Usage: hcptf gpgkey delete [options]

  Delete a GPG key from the private registry.

Options:

  -namespace=<name>  Namespace (organization name) (required)
  -key-id=<id>       GPG key ID (required)

Example:

  hcptf gpgkey delete -namespace=my-org -key-id=32966F3FB5AC1129
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the GPG key delete command
func (c *GPGKeyDeleteCommand) Synopsis() string {
	return "Delete a GPG key from the private registry"
}

func (c *GPGKeyDeleteCommand) gpgService(client *client.Client) gpgKeyDeleter {
	if c.gpgKeySvc != nil {
		return c.gpgKeySvc
	}
	return client.GPGKeys
}
