package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicySetListCommand is a command to list policy sets
type PolicySetListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the policy set list command
func (c *PolicySetListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset list")
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

	// List policy sets
	policySets, err := client.PolicySets.List(client.Context(), c.organization, &tfe.PolicySetListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policy sets: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(policySets.Items) == 0 {
		c.Ui.Output("No policy sets found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Description", "Global", "Policy Count", "Workspace Count"}
	var rows [][]string

	for _, ps := range policySets.Items {
		description := ps.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}

		rows = append(rows, []string{
			ps.ID,
			ps.Name,
			description,
			fmt.Sprintf("%t", ps.Global),
			fmt.Sprintf("%d", ps.PolicyCount),
			fmt.Sprintf("%d", ps.WorkspaceCount),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the policy set list command
func (c *PolicySetListCommand) Help() string {
	helpText := `
Usage: hcptf policyset list [options]

  List policy sets in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policyset list -org=my-org
  hcptf policyset list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set list command
func (c *PolicySetListCommand) Synopsis() string {
	return "List policy sets"
}
