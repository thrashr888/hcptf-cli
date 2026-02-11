package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// GCPoidcReadCommand is a command to read GCP OIDC configuration details
type GCPoidcReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the gcpoidc read command
func (c *GCPoidcReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gcpoidc read")
	flags.StringVar(&c.id, "id", "", "GCP OIDC configuration ID (required)")
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

	// Read GCP OIDC configuration
	config, err := client.GCPOIDCConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading GCP OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                   config.ID,
		"ServiceAccountEmail":  config.ServiceAccountEmail,
		"WorkloadProviderName": config.WorkloadProviderName,
		"ProjectNumber":        config.ProjectNumber,
	}

	if config.Organization != nil {
		data["Organization"] = config.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the gcpoidc read command
func (c *GCPoidcReadCommand) Help() string {
	helpText := `
Usage: hcptf gcpoidc read [options]

  Read GCP OIDC configuration details.

  Displays the configuration details for a GCP OIDC configuration,
  including the service account email, workload provider name,
  project number, and audience settings.

Options:

  -id=<id>          GCP OIDC configuration ID (required)
                    Format: gcpoidc-XXXXXXXXXX
  -output=<format>  Output format: table (default) or json

Example:

  hcptf gcpoidc read -id=gcpoidc-ABC123
  hcptf gcpoidc read -id=gcpoidc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the gcpoidc read command
func (c *GCPoidcReadCommand) Synopsis() string {
	return "Read GCP OIDC configuration details"
}
