package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// HYOKCreateCommand is a command to create a HYOK configuration
type HYOKCreateCommand struct {
	Meta
	organization      string
	name              string
	kekID             string
	agentPoolID       string
	oidcConfigID      string
	oidcType          string
	keyRegion         string
	keyLocation       string
	keyRingID         string
	format            string
}

// Run executes the HYOK create command
func (c *HYOKCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyok create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.name, "name", "", "HYOK configuration name (required)")
	flags.StringVar(&c.kekID, "kek-id", "", "Key Encryption Key ID from your KMS (required)")
	flags.StringVar(&c.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
	flags.StringVar(&c.oidcConfigID, "oidc-config-id", "", "OIDC configuration ID (required)")
	flags.StringVar(&c.oidcType, "oidc-type", "", "OIDC type: aws, azure, gcp, or vault (required)")
	flags.StringVar(&c.keyRegion, "key-region", "", "AWS KMS key region (for AWS KMS only)")
	flags.StringVar(&c.keyLocation, "key-location", "", "GCP key location (for GCP Cloud KMS only)")
	flags.StringVar(&c.keyRingID, "key-ring-id", "", "GCP key ring ID (for GCP Cloud KMS only)")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.kekID == "" {
		c.Ui.Error("Error: -kek-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.agentPoolID == "" {
		c.Ui.Error("Error: -agent-pool-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.oidcConfigID == "" {
		c.Ui.Error("Error: -oidc-config-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.oidcType == "" {
		c.Ui.Error("Error: -oidc-type flag is required (must be: aws, azure, gcp, or vault)")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build KMS options if any were provided
	var kmsOptions *tfe.KMSOptions
	if c.keyRegion != "" || c.keyLocation != "" || c.keyRingID != "" {
		kmsOptions = &tfe.KMSOptions{
			KeyRegion:   c.keyRegion,
			KeyLocation: c.keyLocation,
			KeyRingID:   c.keyRingID,
		}
	}

	// Build OIDC configuration based on type
	oidcConfig := &tfe.OIDCConfigurationTypeChoice{}
	switch strings.ToLower(c.oidcType) {
	case "aws":
		oidcConfig.AWSOIDCConfiguration = &tfe.AWSOIDCConfiguration{ID: c.oidcConfigID}
	case "azure":
		oidcConfig.AzureOIDCConfiguration = &tfe.AzureOIDCConfiguration{ID: c.oidcConfigID}
	case "gcp":
		oidcConfig.GCPOIDCConfiguration = &tfe.GCPOIDCConfiguration{ID: c.oidcConfigID}
	case "vault":
		oidcConfig.VaultOIDCConfiguration = &tfe.VaultOIDCConfiguration{ID: c.oidcConfigID}
	default:
		c.Ui.Error("Error: -oidc-type must be one of: aws, azure, gcp, vault")
		return 1
	}

	// Create HYOK configuration
	options := tfe.HYOKConfigurationsCreateOptions{
		Name:       c.name,
		KEKID:      c.kekID,
		KMSOptions: kmsOptions,
		AgentPool: &tfe.AgentPool{
			ID: c.agentPoolID,
		},
		OIDCConfiguration: oidcConfig,
	}

	config, err := client.HYOKConfigurations.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating HYOK configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("HYOK configuration '%s' created successfully", config.Name))

	// Show configuration details
	data := map[string]interface{}{
		"ID":             config.ID,
		"Name":           config.Name,
		"KEK ID":         config.KEKID,
		"Primary":        config.Primary,
		"Status":         string(config.Status),
		"Agent Pool ID":  c.agentPoolID,
		"OIDC Config ID": c.oidcConfigID,
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

// Help returns help text for the HYOK create command
func (c *HYOKCreateCommand) Help() string {
	helpText := `
Usage: hcptf hyok create [options]

  Create a HYOK (Hold Your Own Key) configuration.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system:
  - AWS KMS: Use -oidc-type=aws and -key-region to specify the AWS region
  - Azure Key Vault: Use -oidc-type=azure (no additional options required)
  - GCP Cloud KMS: Use -oidc-type=gcp with -key-location and -key-ring-id
  - Vault: Use -oidc-type=vault (no additional options required)

Options:

  -organization=<name>    Organization name (required)
  -name=<name>            HYOK configuration name (required)
  -kek-id=<id>            Key Encryption Key ID from your KMS (required)
  -agent-pool-id=<id>     Agent pool ID (required)
  -oidc-config-id=<id>    OIDC configuration ID (required)
  -oidc-type=<type>       OIDC type: aws, azure, gcp, or vault (required)
  -key-region=<region>    AWS KMS key region (for AWS KMS only)
  -key-location=<loc>     GCP key location (for GCP Cloud KMS only)
  -key-ring-id=<id>       GCP key ring ID (for GCP Cloud KMS only)
  -output=<format>        Output format: table (default) or json

Example:

  # AWS KMS
  hcptf hyok create \
    -organization=my-org \
    -name=my-aws-key \
    -kek-id=arn:aws:kms:us-west-2:123456789012:key/12345678-1234-1234-1234-123456789012 \
    -agent-pool-id=apool-123 \
    -oidc-config-id=awsoidc-456 \
    -oidc-type=aws \
    -key-region=us-west-2

  # GCP Cloud KMS
  hcptf hyok create \
    -organization=my-org \
    -name=my-gcp-key \
    -kek-id=my-key \
    -agent-pool-id=apool-123 \
    -oidc-config-id=gcpoidc-456 \
    -oidc-type=gcp \
    -key-location=us-central1 \
    -key-ring-id=my-keyring
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK create command
func (c *HYOKCreateCommand) Synopsis() string {
	return "Create a HYOK configuration"
}
