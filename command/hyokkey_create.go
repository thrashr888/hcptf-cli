package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// HYOKKeyCreateCommand is a command to check for new HYOK customer key versions
type HYOKKeyCreateCommand struct {
	Meta
	hyokConfigID string
	format       string
}

// Run executes the HYOK key create command
func (c *HYOKKeyCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyokkey create")
	flags.StringVar(&c.hyokConfigID, "hyok-config-id", "", "HYOK configuration ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.hyokConfigID == "" {
		c.Ui.Error("Error: -hyok-config-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Check for new key versions (this triggers creation of new versions if available in KMS)
	keyVersions, err := client.HYOKCustomerKeyVersions.List(client.Context(), c.hyokConfigID, &tfe.HYOKCustomerKeyVersionListOptions{
		Refresh: true,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error checking for new key versions: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(keyVersions.Items) == 0 {
		c.Ui.Output("No key versions found")
		return 0
	}

	c.Ui.Output(fmt.Sprintf("Found %d key version(s)", len(keyVersions.Items)))

	// Prepare table data
	headers := []string{"ID", "Key Version", "Status", "Workspaces Secured", "Created At"}
	var rows [][]string

	for _, kv := range keyVersions.Items {
		rows = append(rows, []string{
			kv.ID,
			kv.KeyVersion,
			string(kv.Status),
			fmt.Sprintf("%d", kv.WorkspacesSecured),
			kv.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the HYOK key create command
func (c *HYOKKeyCreateCommand) Help() string {
	helpText := `
Usage: hcptf hyokkey create [options]

  Check for and register new HYOK customer key versions from your KMS.

  This command queries your key management system for new versions of the
  configured key and registers them with HCP Terraform. Key versions are
  automatically created when you rotate keys in your KMS.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -hyok-config-id=<id>  HYOK configuration ID (required)
  -output=<format>      Output format: table (default) or json

Example:

  hcptf hyokkey create -hyok-config-id=hyokc-123456
  hcptf hyokkey create -hyok-config-id=hyokc-123456 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK key create command
func (c *HYOKKeyCreateCommand) Synopsis() string {
	return "Check for and register new HYOK customer key versions"
}
