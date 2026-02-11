package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AgentPoolCreateCommand is a command to create an agent pool
type AgentPoolCreateCommand struct {
	Meta
	organization       string
	name               string
	organizationScoped bool
	format             string
}

// Run executes the agent pool create command
func (c *AgentPoolCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Agent pool name (required)")
	flags.BoolVar(&c.organizationScoped, "organization-scoped", false, "Make agent pool organization scoped")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create agent pool
	options := tfe.AgentPoolCreateOptions{
		Name:               tfe.String(c.name),
		OrganizationScoped: tfe.Bool(c.organizationScoped),
	}

	agentPool, err := client.AgentPools.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating agent pool: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Agent pool '%s' created successfully", agentPool.Name))

	// Show agent pool details
	data := map[string]interface{}{
		"ID":                 agentPool.ID,
		"Name":               agentPool.Name,
		"Organization":       c.organization,
		"AgentCount":         agentPool.AgentCount,
		"OrganizationScoped": agentPool.OrganizationScoped,
		"CreatedAt":          agentPool.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the agent pool create command
func (c *AgentPoolCreateCommand) Help() string {
	helpText := `
Usage: hcptf agentpool create [options]

  Create a new agent pool.

Options:

  -organization=<name>      Organization name (required)
  -org=<name>              Alias for -organization
  -name=<name>             Agent pool name (required)
  -organization-scoped     Make agent pool organization scoped (default: false)
  -output=<format>         Output format: table (default) or json

Example:

  hcptf agentpool create -org=my-org -name=my-pool
  hcptf agentpool create -org=my-org -name=shared-pool -organization-scoped
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool create command
func (c *AgentPoolCreateCommand) Synopsis() string {
	return "Create a new agent pool"
}
