package command

import (
	"fmt"
	"strings"

)

// AgentReadCommand is a command to read agent details
type AgentReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the agent read command
func (c *AgentReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agent read")
	flags.StringVar(&c.id, "id", "", "Agent ID (required)")
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

	// Read agent
	agent, err := client.Agents.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading agent: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":     agent.ID,
		"Status": agent.Status,
	}

	// Add optional fields
	if agent.Name != "" {
		data["Name"] = agent.Name
	}

	if agent.IP != "" {
		data["IPAddress"] = agent.IP
	}

	if agent.LastPingAt != "" {
		data["LastPingAt"] = agent.LastPingAt
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the agent read command
func (c *AgentReadCommand) Help() string {
	helpText := `
Usage: hcptf agent read [options]

  Read agent details and status.

  This command retrieves detailed information about a specific agent, including
  its current status, IP address, and last communication time. Use this to
  monitor the health of individual self-hosted Terraform agents.

Options:

  -id=<id>             Agent ID (required)
  -output=<format>     Output format: table (default) or json

Agent Status Values:

  idle     Agent is online and available to accept jobs
  busy     Agent is currently executing a job
  unknown  Agent status is unknown (may indicate connectivity issues)
  exited   Agent has gracefully shut down
  errored  Agent encountered an error

Example:

  hcptf agent read -id=agent-abc123
  hcptf agent read -id=agent-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent read command
func (c *AgentReadCommand) Synopsis() string {
	return "Read agent details and status"
}
