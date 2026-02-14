package command

import (
	"fmt"
	"strings"
)

// AgentPoolReadCommand is a command to read agent pool details
type AgentPoolReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the agent pool read command
func (c *AgentPoolReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool read")
	flags.StringVar(&c.id, "id", "", "Agent pool ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Read agent pool
	agentPool, err := client.AgentPools.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading agent pool: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

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

// Help returns help text for the agent pool read command
func (c *AgentPoolReadCommand) Help() string {
	helpText := `
Usage: hcptf agentpool read [options]

  Read agent pool details.

Options:

  -id=<id>             Agent pool ID (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf agentpool read -id=apool-abc123
  hcptf agentpool read -id=apool-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool read command
func (c *AgentPoolReadCommand) Synopsis() string {
	return "Read agent pool details"
}
