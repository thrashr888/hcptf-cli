package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// HYOKListCommand is a command to list HYOK configurations
type HYOKListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the HYOK list command
func (c *HYOKListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("hyok list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List HYOK configurations
	configs, err := client.HYOKConfigurations.List(client.Context(), c.organization, &tfe.HYOKConfigurationsListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing HYOK configurations: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(configs.Items) == 0 {
		c.Ui.Output("No HYOK configurations found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "KEK ID", "Primary", "Status"}
	var rows [][]string

	for _, config := range configs.Items {
		primary := "false"
		if config.Primary {
			primary = "true"
		}

		rows = append(rows, []string{
			config.ID,
			config.Name,
			config.KEKID,
			primary,
			string(config.Status),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the HYOK list command
func (c *HYOKListCommand) Help() string {
	helpText := `
Usage: hcptf hyok list [options]

  List HYOK (Hold Your Own Key) configurations for an organization.

  HYOK is an Enterprise-only feature that lets you manage encryption keys
  using your own key management system (AWS KMS, Azure Key Vault, or GCP Cloud KMS).

Options:

  -organization=<name>  Organization name (required)
  -output=<format>      Output format: table (default) or json

Example:

  hcptf hyok list -organization=my-org
  hcptf hyok list -organization=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the HYOK list command
func (c *HYOKListCommand) Synopsis() string {
	return "List HYOK configurations for an organization"
}
