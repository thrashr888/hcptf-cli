package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// AgentPoolTokenCreateCommand is a command to create an agent token
type AgentPoolTokenCreateCommand struct {
	Meta
	agentPoolID string
	description string
	format      string
}

// Run executes the agent pool token create command
func (c *AgentPoolTokenCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool token-create")
	flags.StringVar(&c.agentPoolID, "agent-pool-id", "", "Agent pool ID (required)")
	flags.StringVar(&c.description, "description", "", "Agent token description (required)")
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

	if c.description == "" {
		c.Ui.Error("Error: -description flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create agent token
	options := tfe.AgentTokenCreateOptions{
		Description: tfe.String(c.description),
	}

	agentToken, err := client.AgentTokens.Create(client.Context(), c.agentPoolID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating agent token: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Agent token created successfully"))

	// Show agent token details
	data := map[string]interface{}{
		"ID":          agentToken.ID,
		"Description": agentToken.Description,
		"Token":       agentToken.Token,
		"CreatedAt":   agentToken.CreatedAt,
	}

	if agentToken.CreatedBy != nil {
		data["CreatedBy"] = agentToken.CreatedBy.Username
	}

	formatter.KeyValue(data)

	// Warning about token visibility
	c.Ui.Warn("\nWARNING: This is the only time the token will be displayed. Save it securely.")

	return 0
}

// Help returns help text for the agent pool token create command
func (c *AgentPoolTokenCreateCommand) Help() string {
	helpText := `
Usage: hcptf agentpool token-create [options]

  Create an agent token for an agent pool.

Options:

  -agent-pool-id=<id>   Agent pool ID (required)
  -description=<text>   Agent token description (required)
  -output=<format>      Output format: table (default) or json

Example:

  hcptf agentpool token-create -agent-pool-id=apool-abc123 -description="Production agent"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool token create command
func (c *AgentPoolTokenCreateCommand) Synopsis() string {
	return "Create an agent token for an agent pool"
}
