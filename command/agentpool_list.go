package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AgentPoolListCommand is a command to list agent pools
type AgentPoolListCommand struct {
	Meta
	organization string
	format       string
}

// Run executes the agent pool list command
func (c *AgentPoolListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool list")
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

	// List agent pools
	agentPools, err := client.AgentPools.List(client.Context(), c.organization, &tfe.AgentPoolListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing agent pools: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(agentPools.Items) == 0 {
		c.Ui.Output("No agent pools found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Agent Count", "Organization Scoped"}
	var rows [][]string

	for _, pool := range agentPools.Items {
		orgScoped := "false"
		if pool.OrganizationScoped {
			orgScoped = "true"
		}

		rows = append(rows, []string{
			pool.ID,
			pool.Name,
			fmt.Sprintf("%d", pool.AgentCount),
			orgScoped,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the agent pool list command
func (c *AgentPoolListCommand) Help() string {
	helpText := `
Usage: hcptf agentpool list [options]

  List agent pools in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf agentpool list -organization=my-org
  hcptf agentpool list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool list command
func (c *AgentPoolListCommand) Synopsis() string {
	return "List agent pools in an organization"
}
