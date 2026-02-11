package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// HYOKReadCommand is a command to show HYOK configuration details
type HYOKReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the HYOK read command
func (c *HYOKReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyok read")
	flags.StringVar(&c.id, "id", "", "HYOK configuration ID (required)")
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

	// Read HYOK configuration
	config, err := client.HYOKConfigurations.Read(client.Context(), c.id, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading HYOK configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":      config.ID,
		"Name":    config.Name,
		"KEK ID":  config.KEKID,
		"Primary": config.Primary,
		"Status":  string(config.Status),
	}

	if config.Error != nil {
		data["Error"] = *config.Error
	}

	if config.KMSOptions != nil {
		kmsData := make(map[string]string)
		if config.KMSOptions.KeyRegion != "" {
			kmsData["KeyRegion"] = config.KMSOptions.KeyRegion
		}
		if config.KMSOptions.KeyLocation != "" {
			kmsData["KeyLocation"] = config.KMSOptions.KeyLocation
		}
		if config.KMSOptions.KeyRingID != "" {
			kmsData["KeyRingID"] = config.KMSOptions.KeyRingID
		}
		if len(kmsData) > 0 {
			data["KMS Options"] = kmsData
		}
	}

	if config.Organization != nil {
		data["Organization"] = config.Organization.Name
	}

	if config.AgentPool != nil {
		data["Agent Pool ID"] = config.AgentPool.ID
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the HYOK read command
func (c *HYOKReadCommand) Help() string {
	helpText := `
Usage: hcptf hyok read [options]

  Show HYOK (Hold Your Own Key) configuration details.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -id=<id>          HYOK configuration ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf hyok read -id=hyokc-123456
  hcptf hyok read -id=hyokc-123456 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK read command
func (c *HYOKReadCommand) Synopsis() string {
	return "Show HYOK configuration details"
}
