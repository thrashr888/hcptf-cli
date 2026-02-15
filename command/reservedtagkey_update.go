package command

import (
	"fmt"
	"strconv"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// ReservedTagKeyUpdateCommand is a command to update a reserved tag key.
type ReservedTagKeyUpdateCommand struct {
	Meta
	id                string
	key               string
	disableOverrides  string
	format            string
	reservedTagKeySvc reservedTagKeyUpdater
}

// Run executes the reservedtagkey update command.
func (c *ReservedTagKeyUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("reservedtagkey update")
	flags.StringVar(&c.id, "id", "", "Reserved tag key ID (required)")
	flags.StringVar(&c.key, "key", "", "Updated tag key")
	flags.StringVar(&c.disableOverrides, "disable-overrides", "", "Set disable-overrides (true/false)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.key == "" && c.disableOverrides == "" {
		c.Ui.Error("Error: at least one of -key or -disable-overrides must be provided")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	options := tfe.ReservedTagKeyUpdateOptions{}
	if c.key != "" {
		options.Key = tfe.String(c.key)
	}
	if c.disableOverrides != "" {
		parsed, err := strconv.ParseBool(c.disableOverrides)
		if err != nil {
			c.Ui.Error("Error: -disable-overrides must be true or false")
			return 1
		}
		options.DisableOverrides = tfe.Bool(parsed)
	}

	reservedKey, err := c.reservedTagKeyService(client).Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating reserved tag key: %s", err))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	c.Ui.Output(fmt.Sprintf("Reserved tag key '%s' updated successfully", reservedKey.ID))
	formatter.KeyValue(map[string]interface{}{
		"ID":               reservedKey.ID,
		"Key":              reservedKey.Key,
		"DisableOverrides": reservedKey.DisableOverrides,
		"CreatedAt":        reservedKey.CreatedAt,
		"UpdatedAt":        reservedKey.UpdatedAt,
	})
	return 0
}

func (c *ReservedTagKeyUpdateCommand) reservedTagKeyService(client *client.Client) reservedTagKeyUpdater {
	if c.reservedTagKeySvc != nil {
		return c.reservedTagKeySvc
	}
	return client.ReservedTagKeys
}

// Help returns help text for the reservedtagkey update command.
func (c *ReservedTagKeyUpdateCommand) Help() string {
	helpText := `
Usage: hcptf reservedtagkey update [options]

  Update a reserved tag key.

Options:

  -id=<id>                 Reserved tag key ID (required)
  -key=<key>               Updated key name
  -disable-overrides=<v>   Set disable-overrides (true/false)
  -output=<format>         Output format: table (default) or json

Example:

  hcptf reservedtagkey update -id=rtk-abc123 -key=environment
  hcptf reservedtagkey update -id=rtk-abc123 -disable-overrides=true
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the reservedtagkey update command.
func (c *ReservedTagKeyUpdateCommand) Synopsis() string {
	return "Update a reserved tag key"
}
