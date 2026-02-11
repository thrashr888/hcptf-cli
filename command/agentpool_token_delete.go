package command

import (
	"fmt"
	"strings"
)

// AgentPoolTokenDeleteCommand is a command to delete an agent token
type AgentPoolTokenDeleteCommand struct {
	Meta
	id    string
	force bool
}

// Run executes the agent pool token delete command
func (c *AgentPoolTokenDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool token-delete")
	flags.StringVar(&c.id, "id", "", "Agent token ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

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

	// Confirm deletion unless force flag is set
	if !c.force {
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete agent token '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete agent token
	err = client.AgentTokens.Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting agent token: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Agent token '%s' deleted successfully", c.id))
	return 0
}

// Help returns help text for the agent pool token delete command
func (c *AgentPoolTokenDeleteCommand) Help() string {
	helpText := `
Usage: hcptf agentpool token-delete [options]

  Delete an agent token.

Options:

  -id=<id>     Agent token ID (required)
  -force       Force delete without confirmation

Example:

  hcptf agentpool token-delete -id=at-abc123
  hcptf agentpool token-delete -id=at-abc123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool token delete command
func (c *AgentPoolTokenDeleteCommand) Synopsis() string {
	return "Delete an agent token"
}
