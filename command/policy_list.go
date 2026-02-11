package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicyListCommand is a command to list policies
type PolicyListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the policy list command
func (c *PolicyListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy list")
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

	// List policies
	policies, err := client.Policies.List(client.Context(), c.organization, &tfe.PolicyListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policies: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(policies.Items) == 0 {
		c.Ui.Output("No policies found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Enforce Level", "Policy Sets", "Updated At"}
	var rows [][]string

	for _, policy := range policies.Items {
		rows = append(rows, []string{
			policy.ID,
			policy.Name,
			string(policy.EnforcementLevel),
			fmt.Sprintf("%d", policy.PolicySetCount),
			policy.UpdatedAt.String(),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the policy list command
func (c *PolicyListCommand) Help() string {
	helpText := `
Usage: hcptf policy list [options]

  List policies in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policy list -organization=my-org
  hcptf policy list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy list command
func (c *PolicyListCommand) Synopsis() string {
	return "List policies in an organization"
}
