package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// GCPoidcCreateCommand is a command to create a GCP OIDC configuration
type GCPoidcCreateCommand struct {
	Meta
	organization         string
	serviceAccountEmail  string
	workloadProviderName string
	projectNumber        string
	format               string
}

// Run executes the gcpoidc create command
func (c *GCPoidcCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("gcpoidc create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.serviceAccountEmail, "service-account-email", "", "GCP service account email (required)")
	flags.StringVar(&c.workloadProviderName, "workload-provider-name", "", "Workload provider path (required)")
	flags.StringVar(&c.projectNumber, "project-number", "", "GCP project number (required)")
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

	if c.serviceAccountEmail == "" {
		c.Ui.Error("Error: -service-account-email flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workloadProviderName == "" {
		c.Ui.Error("Error: -workload-provider-name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.projectNumber == "" {
		c.Ui.Error("Error: -project-number flag is required")
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
	options := tfe.GCPOIDCConfigurationCreateOptions{
		ServiceAccountEmail:  c.serviceAccountEmail,
		WorkloadProviderName: c.workloadProviderName,
		ProjectNumber:        c.projectNumber,
	}

	// Create GCP OIDC configuration
	config, err := client.GCPOIDCConfigurations.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating GCP OIDC configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("GCP OIDC configuration created successfully with ID: %s", config.ID))

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

// Help returns help text for the gcpoidc create command
func (c *GCPoidcCreateCommand) Help() string {
	helpText := `
Usage: hcptf gcpoidc create [options]

  Create a new GCP OIDC configuration for dynamic GCP credentials.

  GCP OIDC configurations enable HCP Terraform to dynamically generate
  GCP credentials using OpenID Connect. This eliminates the need to store
  static GCP service account keys in HCP Terraform.

  Prerequisites:
  - GCP Workload Identity Pool configured for HCP Terraform
  - Service account with appropriate permissions for required GCP operations
  - Service account bound to the workload identity pool

Options:

  -organization=<name>           Organization name (required)
  -org=<name>                   Alias for -organization
  -service-account-email=<email> GCP service account email (required)
                                Format: sa-name@project-id.iam.gserviceaccount.com
  -workload-provider-name=<path> Fully qualified workload provider path (required)
                                Format: projects/PROJECT_NUMBER/locations/global/
                                        workloadIdentityPools/POOL_ID/providers/PROVIDER_ID
  -project-number=<number>      GCP project number (required)
  -output=<format>              Output format: table (default) or json

Example:

  hcptf gcpoidc create -org=my-org \
    -service-account-email=terraform@my-project.iam.gserviceaccount.com \
    -workload-provider-name=projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider \
    -project-number=123456789
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the gcpoidc create command
func (c *GCPoidcCreateCommand) Synopsis() string {
	return "Create a GCP OIDC configuration for dynamic credentials"
}
