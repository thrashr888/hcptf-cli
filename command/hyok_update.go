package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// HYOKUpdateCommand is a command to update a HYOK configuration
type HYOKUpdateCommand struct {
	Meta
	id          string
	name        string
	kekID       string
	primary     string
	keyRegion   string
	keyLocation string
	keyRingID   string
	format      string
}

// Run executes the HYOK update command
func (c *HYOKUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyok update")
	flags.StringVar(&c.id, "id", "", "HYOK configuration ID (required)")
	flags.StringVar(&c.name, "name", "", "HYOK configuration name")
	flags.StringVar(&c.kekID, "kek-id", "", "Key Encryption Key ID from your KMS")
	flags.StringVar(&c.primary, "primary", "", "Set as primary HYOK configuration (true/false)")
	flags.StringVar(&c.keyRegion, "key-region", "", "AWS KMS key region (for AWS KMS only)")
	flags.StringVar(&c.keyLocation, "key-location", "", "GCP key location (for GCP Cloud KMS only)")
	flags.StringVar(&c.keyRingID, "key-ring-id", "", "GCP key ring ID (for GCP Cloud KMS only)")
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

	// Build update options
	options := tfe.HYOKConfigurationsUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.kekID != "" {
		options.KEKID = tfe.String(c.kekID)
	}

	if c.primary != "" {
		if c.primary == "true" {
			options.Primary = tfe.Bool(true)
		} else if c.primary == "false" {
			options.Primary = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -primary must be 'true' or 'false'")
			return 1
		}
	}

	// Build KMS options if any were provided
	if c.keyRegion != "" || c.keyLocation != "" || c.keyRingID != "" {
		options.KMSOptions = &tfe.KMSOptions{
			KeyRegion:   c.keyRegion,
			KeyLocation: c.keyLocation,
			KeyRingID:   c.keyRingID,
		}
	}

	// Update HYOK configuration
	config, err := client.HYOKConfigurations.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating HYOK configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("HYOK configuration '%s' updated successfully", config.Name))

	// Show configuration details
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

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the HYOK update command
func (c *HYOKUpdateCommand) Help() string {
	helpText := `
Usage: hcptf hyok update [options]

  Update a HYOK (Hold Your Own Key) configuration.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -id=<id>                HYOK configuration ID (required)
  -name=<name>            HYOK configuration name
  -kek-id=<id>            Key Encryption Key ID from your KMS
  -primary=<bool>         Set as primary HYOK configuration (true/false)
  -key-region=<region>    AWS KMS key region (for AWS KMS only)
  -key-location=<loc>     GCP key location (for GCP Cloud KMS only)
  -key-ring-id=<id>       GCP key ring ID (for GCP Cloud KMS only)
  -output=<format>        Output format: table (default) or json

Example:

  hcptf hyok update -id=hyokc-123456 -name=updated-name
  hcptf hyok update -id=hyokc-123456 -primary=true
  hcptf hyok update -id=hyokc-123456 -key-region=us-east-1
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK update command
func (c *HYOKUpdateCommand) Synopsis() string {
	return "Update a HYOK configuration"
}
