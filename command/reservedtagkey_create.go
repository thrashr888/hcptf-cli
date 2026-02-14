package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ReservedTagKeyCreateCommand is a command to create a reserved tag key
type ReservedTagKeyCreateCommand struct {
	Meta
	organization      string
	key               string
	disableOverrides  bool
	format            string
	reservedTagKeySvc reservedTagKeyCreator
}

// Run executes the reservedtagkey create command
func (c *ReservedTagKeyCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("reservedtagkey create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.key, "key", "", "Tag key to reserve (required)")
	flags.BoolVar(&c.disableOverrides, "disable-overrides", false, "Disable overriding inherited tags at workspace level")
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

	if c.key == "" {
		c.Ui.Error("Error: -key flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.ReservedTagKeyCreateOptions{
		Key:              c.key,
		DisableOverrides: tfe.Bool(c.disableOverrides),
	}

	// Create reserved tag key
	reservedKey, err := c.reservedTagKeyService(client).Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating reserved tag key: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Reserved tag key '%s' created successfully", reservedKey.Key))

	// Show reserved tag key details
	data := map[string]interface{}{
		"ID":               reservedKey.ID,
		"Key":              reservedKey.Key,
		"DisableOverrides": reservedKey.DisableOverrides,
		"CreatedAt":        reservedKey.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ReservedTagKeyCreateCommand) reservedTagKeyService(client *client.Client) reservedTagKeyCreator {
	if c.reservedTagKeySvc != nil {
		return c.reservedTagKeySvc
	}
	return client.ReservedTagKeys
}

// Help returns help text for the reservedtagkey create command
func (c *ReservedTagKeyCreateCommand) Help() string {
	helpText := `
Usage: hcptf reservedtagkey create [options]

  Create a reserved tag key for an organization.
  Reserved tag keys enable consistent tagging strategies and can
  prevent workspaces from overriding inherited project tags.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -key=<key>           Tag key to reserve (required)
  -disable-overrides   Disable overriding inherited tags at workspace level
  -output=<format>     Output format: table (default) or json

Example:

  hcptf reservedtagkey create -org=my-org -key=environment
  hcptf reservedtagkey create -org=my-org -key=cost-center -disable-overrides
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the reservedtagkey create command
func (c *ReservedTagKeyCreateCommand) Synopsis() string {
	return "Create a reserved tag key"
}
