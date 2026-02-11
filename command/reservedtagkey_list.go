package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// ReservedTagKeyListCommand is a command to list reserved tag keys
type ReservedTagKeyListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the reservedtagkey list command
func (c *ReservedTagKeyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("reservedtagkey list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// List reserved tag keys
	keys, err := client.ReservedTagKeys.List(client.Context(), c.organization, &tfe.ReservedTagKeyListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing reserved tag keys: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(keys.Items) == 0 {
		c.Ui.Output("No reserved tag keys found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Key", "Disable Overrides", "Created At"}
	var rows [][]string

	for _, key := range keys.Items {
		disableOverrides := "false"
		if key.DisableOverrides {
			disableOverrides = "true"
		}

		rows = append(rows, []string{
			key.ID,
			key.Key,
			disableOverrides,
			key.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the reservedtagkey list command
func (c *ReservedTagKeyListCommand) Help() string {
	helpText := `
Usage: hcptf reservedtagkey list [options]

  List reserved tag keys for an organization.
  Reserved tag keys enable consistent tagging strategies and can
  prevent workspaces from overriding inherited project tags.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf reservedtagkey list -org=my-org
  hcptf reservedtagkey list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the reservedtagkey list command
func (c *ReservedTagKeyListCommand) Synopsis() string {
	return "List reserved tag keys"
}
