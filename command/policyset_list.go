package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetListCommand is a command to list policy sets
type PolicySetListCommand struct {
	Meta
	organization string
	search       string
	kind         string
	include      string
	format       string
	policySetSvc policySetLister
}

// Run executes the policy set list command
func (c *PolicySetListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.search, "search", "", "Search policy set names by substring")
	flags.StringVar(&c.kind, "kind", "", "Filter by policy set kind: sentinel or opa")
	flags.StringVar(&c.include, "include", "", "Include related resources (comma-separated)")
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

	options := &tfe.PolicySetListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
		Search: c.search,
	}
	if c.kind != "" {
		kind, parseErr := parsePolicyKind(c.kind)
		if parseErr != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
			return 1
		}
		options.Kind = kind
	}
	if c.include != "" {
		for _, include := range splitCommaList(c.include) {
			if include == "" {
				continue
			}
			options.Include = append(options.Include, tfe.PolicySetIncludeOpt(include))
		}
	}

	// List policy sets
	policySets, err := c.policySetService(client).List(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policy sets: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

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
  -search=<query>      Search policy set names by substring
  -kind=<kind>         Filter by kind: sentinel or opa
  -include=<values>    Include related resources (comma-separated)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policyset list -org=my-org
  hcptf policyset list -org=my-org -search=security -kind=opa -include=projects,policies
  hcptf policyset list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *PolicySetListCommand) policySetService(client *client.Client) policySetLister {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Synopsis returns a short synopsis for the policy set list command
func (c *PolicySetListCommand) Synopsis() string {
	return "List policy sets"
}
