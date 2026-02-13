package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// GCPoidcUpdateCommand is a command to update a GCP OIDC configuration
type GCPoidcUpdateCommand struct {
	Meta
	id                   string
	serviceAccountEmail  string
	workloadProviderName string
	projectNumber        string
	format               string
}

// Run executes the gcpoidc update command
func (c *GCPoidcUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gcpoidc update")
	flags.StringVar(&c.id, "id", "", "GCP OIDC configuration ID (required)")
	flags.StringVar(&c.serviceAccountEmail, "service-account-email", "", "GCP service account email")
	flags.StringVar(&c.workloadProviderName, "workload-provider-name", "", "Workload provider path")
	flags.StringVar(&c.projectNumber, "project-number", "", "GCP project number")
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
	options := tfe.GCPOIDCConfigurationUpdateOptions{}

	if c.serviceAccountEmail != "" {
		options.ServiceAccountEmail = &c.serviceAccountEmail
	}

	if c.workloadProviderName != "" {
		options.WorkloadProviderName = &c.workloadProviderName
	}

	if c.projectNumber != "" {
		options.ProjectNumber = &c.projectNumber
	}

	// Update GCP OIDC configuration
	config, err := client.GCPOIDCConfigurations.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating GCP OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("GCP OIDC configuration '%s' updated successfully", config.ID))

	// Show configuration details
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

// Help returns help text for the gcpoidc update command
func (c *GCPoidcUpdateCommand) Help() string {
	helpText := `
Usage: hcptf gcpoidc update [options]

  Update GCP OIDC configuration settings.

  Updates the service account email, workload provider name, project number,
  or audience settings for an existing GCP OIDC configuration.

Options:

  -id=<id>                      GCP OIDC configuration ID (required)
  -service-account-email=<email> GCP service account email
                                Format: sa-name@project-id.iam.gserviceaccount.com
  -workload-provider-name=<path> Fully qualified workload provider path
                                Format: projects/PROJECT_NUMBER/locations/global/
                                        workloadIdentityPools/POOL_ID/providers/PROVIDER_ID
  -project-number=<number>      GCP project number
  -output=<format>              Output format: table (default) or json

Example:

  hcptf gcpoidc update -id=gcpoidc-ABC123 \
    -service-account-email=terraform@my-project.iam.gserviceaccount.com

  hcptf gcpoidc update -id=gcpoidc-ABC123 \
    -workload-provider-name=projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider \
    -project-number=123456789
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the gcpoidc update command
func (c *GCPoidcUpdateCommand) Synopsis() string {
	return "Update GCP OIDC configuration settings"
}
