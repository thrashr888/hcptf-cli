package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

type GPGKeyListCommand struct {
	Meta
	namespace string
	format    string
	gpgKeySvc gpgKeyLister
}

// Run executes the GPG key list command
func (c *GPGKeyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gpgkey list")
	flags.StringVar(&c.namespace, "namespace", "", "Namespace (organization name) (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List GPG keys
	keys, err := c.gpgService(client).ListPrivate(client.Context(), tfe.GPGKeyListOptions{
		Namespaces: []string{c.namespace},
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing GPG keys: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(keys.Items) == 0 {
		c.Ui.Output("No GPG keys found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Key ID", "Namespace", "Created At"}
	var rows [][]string

	for _, key := range keys.Items {
		rows = append(rows, []string{
			key.ID,
			key.KeyID,
			key.Namespace,
			key.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the GPG key list command
func (c *GPGKeyListCommand) Help() string {
	helpText := `
Usage: hcptf gpgkey list [options]

  List GPG keys for provider signing in an organization.

  GPG keys are used to sign private provider releases to ensure
  authenticity and integrity.

Options:

  -namespace=<name>  Namespace (organization name) (required)
  -output=<format>   Output format: table (default) or json

Example:

  hcptf gpgkey list -namespace=my-org
  hcptf gpgkey list -namespace=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *GPGKeyListCommand) gpgService(client *client.Client) gpgKeyLister {
	if c.gpgKeySvc != nil {
		return c.gpgKeySvc
	}
	return client.GPGKeys
}

// Synopsis returns a short synopsis for the GPG key list command
func (c *GPGKeyListCommand) Synopsis() string {
	return "List GPG keys for provider signing"
}
