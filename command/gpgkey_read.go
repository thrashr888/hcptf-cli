package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type GPGKeyReadCommand struct {
	Meta
	namespace string
	keyID     string
	format    string
	gpgKeySvc gpgKeyReader
}

// Run executes the GPG key read command
func (c *GPGKeyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gpgkey read")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (organization name) (required)")
	flags.StringVar(&c.keyID, "key-id", "", "GPG key ID (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read GPG key
	keyID := tfe.GPGKeyID{
		RegistryName: tfe.PrivateRegistry,
		Namespace:    c.namespace,
		KeyID:        c.keyID,
	}

	key, err := c.gpgService(client).Read(client.Context(), keyID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading GPG key: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	// Show key details
	data := map[string]interface{}{
		"ID":         key.ID,
		"KeyID":      key.KeyID,
		"Namespace":  key.Namespace,
		"CreatedAt":  key.CreatedAt,
		"UpdatedAt":  key.UpdatedAt,
		"AsciiArmor": key.AsciiArmor,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the GPG key read command
func (c *GPGKeyReadCommand) Help() string {
	helpText := `
Usage: hcptf gpgkey read [options]

  Show details of a GPG key including its public key content.

Options:

  -namespace=<name>  Namespace (organization name) (required)
  -key-id=<id>       GPG key ID (required)
  -output=<format>   Output format: table (default) or json

Example:

  hcptf gpgkey read -namespace=my-org -key-id=32966F3FB5AC1129
  hcptf gpgkey read -namespace=my-org -key-id=32966F3FB5AC1129 -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *GPGKeyReadCommand) gpgService(client *client.Client) gpgKeyReader {
	if c.gpgKeySvc != nil {
		return c.gpgKeySvc
	}
	return client.GPGKeys
}

// Synopsis returns a short synopsis for the GPG key read command
func (c *GPGKeyReadCommand) Synopsis() string {
	return "Show details of a GPG key"
}
