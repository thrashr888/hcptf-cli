package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AgentPoolUpdateCommand is a command to update an agent pool
type AgentPoolUpdateCommand struct {
	Meta
	id                 string
	name               string
	organizationScoped *bool
	format             string
}

// Run executes the agent pool update command
func (c *AgentPoolUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool update")
	flags.StringVar(&c.id, "id", "", "Agent pool ID (required)")
	flags.StringVar(&c.name, "name", "", "Agent pool name")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	// Custom bool flag handling for optional organization-scoped
	var orgScopedFlag string
	flags.StringVar(&orgScopedFlag, "organization-scoped", "", "Make agent pool organization scoped (true/false)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Parse organization-scoped flag if provided
	if orgScopedFlag != "" {
		if orgScopedFlag == "true" {
			val := true
			c.organizationScoped = &val
		} else if orgScopedFlag == "false" {
			val := false
			c.organizationScoped = &val
		} else {
			c.Ui.Error("Error: -organization-scoped must be 'true' or 'false'")
			c.Ui.Error(c.Help())
			return 1
		}
	}

	// Check if any update flags are provided
	if c.name == "" && c.organizationScoped == nil {
		c.Ui.Error("Error: At least one update flag must be provided")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Update agent pool
	options := tfe.AgentPoolUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.organizationScoped != nil {
		options.OrganizationScoped = tfe.Bool(*c.organizationScoped)
	}

	agentPool, err := client.AgentPools.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating agent pool: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Agent pool '%s' updated successfully", agentPool.Name))

	// Show agent pool details
	data := map[string]interface{}{
		"ID":                 agentPool.ID,
		"Name":               agentPool.Name,
		"AgentCount":         agentPool.AgentCount,
		"OrganizationScoped": agentPool.OrganizationScoped,
		"CreatedAt":          agentPool.CreatedAt,
	}

	if agentPool.Organization != nil {
		data["Organization"] = agentPool.Organization.Name
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the agent pool update command
func (c *AgentPoolUpdateCommand) Help() string {
	helpText := `
Usage: hcptf agentpool update [options]

  Update an agent pool.

Options:

  -id=<id>                         Agent pool ID (required)
  -name=<name>                     Agent pool name
  -organization-scoped=<true|false> Make agent pool organization scoped
  -output=<format>                 Output format: table (default) or json

Example:

  hcptf agentpool update -id=apool-abc123 -name=new-name
  hcptf agentpool update -id=apool-abc123 -organization-scoped=true
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool update command
func (c *AgentPoolUpdateCommand) Synopsis() string {
	return "Update an agent pool"
}
