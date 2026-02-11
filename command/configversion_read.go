package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ConfigVersionReadCommand is a command to read configuration version details
type ConfigVersionReadCommand struct {
	Meta
	configVersionID string
	format          string
	configVerSvc    configVersionReader
}

// Run executes the configversion read command
func (c *ConfigVersionReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("configversion read")
	flags.StringVar(&c.configVersionID, "id", "", "Configuration version ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.configVersionID == "" {
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

	// Read configuration version
	configVersion, err := c.configVersionService(client).Read(client.Context(), c.configVersionID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading configuration version: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":           configVersion.ID,
		"Status":       configVersion.Status,
		"Source":       configVersion.Source,
		"Speculative":  configVersion.Speculative,
		"Provisional":  configVersion.Provisional,
		"Error":        configVersion.Error,
		"ErrorMessage": configVersion.ErrorMessage,
	}

	if configVersion.UploadURL != "" {
		data["UploadURL"] = configVersion.UploadURL
	}

	formatter.KeyValue(data)
	return 0
}

func (c *ConfigVersionReadCommand) configVersionService(client *client.Client) configVersionReader {
	if c.configVerSvc != nil {
		return c.configVerSvc
	}
	return client.ConfigurationVersions
}

// Help returns help text for the configversion read command
func (c *ConfigVersionReadCommand) Help() string {
	helpText := `
Usage: hcptf configversion read [options]

  Show configuration version details.

Options:

  -id=<config-id>   Configuration version ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf configversion read -id=cv-abc123
  hcptf configversion read -id=cv-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the configversion read command
func (c *ConfigVersionReadCommand) Synopsis() string {
	return "Show configuration version details"
}
