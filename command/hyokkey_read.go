package command

import (
	"fmt"
	"strings"
)

// HYOKKeyReadCommand is a command to show HYOK customer key version details
type HYOKKeyReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the HYOK key read command
func (c *HYOKKeyReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyokkey read")
	flags.StringVar(&c.id, "id", "", "HYOK customer key version ID (required)")
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

	// Read HYOK customer key version
	keyVersion, err := client.HYOKCustomerKeyVersions.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading HYOK customer key version: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                 keyVersion.ID,
		"Key Version":        keyVersion.KeyVersion,
		"Status":             string(keyVersion.Status),
		"Error":              keyVersion.Error,
		"Workspaces Secured": keyVersion.WorkspacesSecured,
		"Created At":         keyVersion.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if keyVersion.HYOKConfiguration != nil {
		data["HYOK Config ID"] = keyVersion.HYOKConfiguration.ID
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the HYOK key read command
func (c *HYOKKeyReadCommand) Help() string {
	helpText := `
Usage: hcptf hyokkey read [options]

  Show HYOK customer key version details.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -id=<id>          HYOK customer key version ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf hyokkey read -id=keyv-123456
  hcptf hyokkey read -id=keyv-123456 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK key read command
func (c *HYOKKeyReadCommand) Synopsis() string {
	return "Show HYOK customer key version details"
}
