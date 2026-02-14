package command

import (
	"fmt"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

type GPGKeyCreateCommand struct {
	Meta
	namespace  string
	asciiArmor string
	file       string
	format     string
	gpgKeySvc  gpgKeyCreator
}

// Run executes the GPG key create command
func (c *GPGKeyCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gpgkey create")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (organization name) (required)")
	flags.StringVar(&c.asciiArmor, "ascii-armor", "", "GPG public key in ASCII armor format")
	flags.StringVar(&c.file, "file", "", "Path to file containing GPG public key")
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

	if c.asciiArmor == "" && c.file == "" {
		c.Ui.Error("Error: either -ascii-armor or -file flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.asciiArmor != "" && c.file != "" {
		c.Ui.Error("Error: cannot specify both -ascii-armor and -file flags")
		c.Ui.Error(c.Help())
		return 1
	}

	// Read from file if specified
	if c.file != "" {
		content, err := os.ReadFile(c.file)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading file: %s", err))
			return 1
		}
		c.asciiArmor = string(content)
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create GPG key
	options := tfe.GPGKeyCreateOptions{
		Namespace:  c.namespace,
		AsciiArmor: c.asciiArmor,
	}

	key, err := c.gpgService(client).Create(client.Context(), tfe.PrivateRegistry, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating GPG key: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("GPG key created successfully with Key ID: %s", key.KeyID))

	// Show key details
	data := map[string]interface{}{
		"ID":        key.ID,
		"KeyID":     key.KeyID,
		"Namespace": key.Namespace,
		"CreatedAt": key.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the GPG key create command
func (c *GPGKeyCreateCommand) Help() string {
	helpText := `
Usage: hcptf gpgkey create [options]

  Upload a GPG public key for signing private providers.

  The GPG key must be in ASCII armor format. You can export your public
  key using: gpg --armor --export your@email.com

Options:

  -namespace=<name>    Namespace (organization name) (required)
  -ascii-armor=<key>   GPG public key in ASCII armor format
  -file=<path>         Path to file containing GPG public key
  -output=<format>     Output format: table (default) or json

Example:

  hcptf gpgkey create -namespace=my-org -file=public-key.asc
  hcptf gpgkey create -namespace=my-org -ascii-armor="-----BEGIN PGP PUBLIC KEY BLOCK-----..."
`
	return strings.TrimSpace(helpText)
}

func (c *GPGKeyCreateCommand) gpgService(client *client.Client) gpgKeyCreator {
	if c.gpgKeySvc != nil {
		return c.gpgKeySvc
	}
	return client.GPGKeys
}

// Synopsis returns a short synopsis for the GPG key create command
func (c *GPGKeyCreateCommand) Synopsis() string {
	return "Upload a GPG public key for signing private providers"
}
