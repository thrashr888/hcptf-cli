package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// AgentPoolDeleteCommand is a command to delete an agent pool
type AgentPoolDeleteCommand struct {
	Meta
	id           string
	force        bool
	agentPoolSvc agentPoolDeleter
}

// Run executes the agent pool delete command
func (c *AgentPoolDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("agentpool delete")
	flags.StringVar(&c.id, "id", "", "Agent pool ID (required)")
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
		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete agent pool '%s'? (yes/no): ", c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete agent pool
	err = c.agentPoolService(client).Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting agent pool: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Agent pool '%s' deleted successfully", c.id))
	return 0
}

func (c *AgentPoolDeleteCommand) agentPoolService(client *client.Client) agentPoolDeleter {
	if c.agentPoolSvc != nil {
		return c.agentPoolSvc
	}
	return client.AgentPools
}

// Help returns help text for the agent pool delete command
func (c *AgentPoolDeleteCommand) Help() string {
	helpText := `
Usage: hcptf agentpool delete [options]

  Delete an agent pool.

Options:

  -id=<id>     Agent pool ID (required)
  -force       Force delete without confirmation

Example:

  hcptf agentpool delete -id=apool-abc123
  hcptf agentpool delete -id=apool-abc123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the agent pool delete command
func (c *AgentPoolDeleteCommand) Synopsis() string {
	return "Delete an agent pool"
}
