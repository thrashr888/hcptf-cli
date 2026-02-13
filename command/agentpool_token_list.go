package command

import (
	"fmt"
	"strings"

)

// AgentPoolTokenListCommand is a command to list agent tokens
type AgentPoolTokenListCommand struct {
	Meta
	agentPoolID string
	format      string
}

// Run executes the agent pool token list command
func (c *AgentPoolTokenListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool token-list")
	flags.StringVar(&c.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
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

	// List agent tokens
	agentTokens, err := client.AgentTokens.List(client.Context(), c.agentPoolID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing agent tokens: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(agentTokens.Items) == 0 {
		c.Ui.Output("No agent tokens found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Description", "Created At", "Last Used At"}
	var rows [][]string

	for _, token := range agentTokens.Items {
		lastUsed := "Never"
		if !token.LastUsedAt.IsZero() {
			lastUsed = token.LastUsedAt.Format("2006-01-02 15:04:05")
		}

		rows = append(rows, []string{
			token.ID,
			token.Description,
			token.CreatedAt.Format("2006-01-02 15:04:05"),
			lastUsed,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the agent pool token list command
func (c *AgentPoolTokenListCommand) Help() string {
	helpText := `
Usage: hcptf agentpool token-list [options]

  List agent tokens in an agent pool.

Options:

  -agent-pool-id=<id>  Agent pool ID (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf agentpool token-list -agent-pool-id=apool-abc123
  hcptf agentpool token-list -agent-pool-id=apool-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool token list command
func (c *AgentPoolTokenListCommand) Synopsis() string {
	return "List agent tokens in an agent pool"
}
