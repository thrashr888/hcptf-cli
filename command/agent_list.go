package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// AgentListCommand is a command to list agents in an agent pool
type AgentListCommand struct {
	Meta
	agentPoolID string
	format      string
}

// Run executes the agent list command
func (c *AgentListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agent list")
	flags.StringVar(&c.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
	flags.StringVar(&c.agentPoolID, "pool", "", "Agent pool ID (alias)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.agentPoolID == "" {
		c.Ui.Error("Error: -agent-pool-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List agents
	agents, err := client.Agents.List(client.Context(), c.agentPoolID, &tfe.AgentListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing agents: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(agents.Items) == 0 {
		c.Ui.Output("No agents found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Status", "IP Address", "Last Ping"}
	var rows [][]string

	for _, agent := range agents.Items {
		name := agent.Name
		if name == "" {
			name = "-"
		}

		ipAddress := agent.IP
		if ipAddress == "" {
			ipAddress = "-"
		}

		lastPing := agent.LastPingAt
		if lastPing == "" {
			lastPing = "-"
		}

		rows = append(rows, []string{
			agent.ID,
			name,
			agent.Status,
			ipAddress,
			lastPing,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the agent list command
func (c *AgentListCommand) Help() string {
	helpText := `
Usage: hcptf agent list [options]

  List agents in an agent pool.

  Agents are self-hosted Terraform Cloud/Enterprise workers that run Terraform
  operations on your own infrastructure. This command shows the status and health
  of agents in a pool, including whether they are idle, busy, or in error states.

Options:

  -agent-pool-id=<id>  Agent pool ID (required)
  -pool=<id>          Alias for -agent-pool-id
  -output=<format>    Output format: table (default) or json

Agent Status Values:

  idle     Agent is online and available to accept jobs
  busy     Agent is currently executing a job
  unknown  Agent status is unknown (may indicate connectivity issues)
  exited   Agent has gracefully shut down
  errored  Agent encountered an error

Example:

  hcptf agent list -agent-pool-id=apool-abc123
  hcptf agent list -pool=apool-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent list command
func (c *AgentListCommand) Synopsis() string {
	return "List agents in an agent pool"
}
