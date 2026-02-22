package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetRemoveWorkspaceCommand removes workspaces from a policy set.
type PolicySetRemoveWorkspaceCommand struct {
	Meta
	policySetID  string
	workspaceIDs string
	policySetSvc policySetWorkspaceRemover
}

// Run executes the policyset remove-workspace command.
func (c *PolicySetRemoveWorkspaceCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset remove-workspace")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.workspaceIDs, "workspace-ids", "", "Comma-separated workspace IDs to remove (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.policySetID == "" {
		c.Ui.Error("Error: -policyset-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.workspaceIDs == "" {
		c.Ui.Error("Error: -workspace-ids flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	ids := splitCommaList(c.workspaceIDs)
	workspaces := make([]*tfe.Workspace, 0, len(ids))
	for _, id := range ids {
		workspaces = append(workspaces, &tfe.Workspace{ID: strings.TrimSpace(id)})
	}
	options := tfe.PolicySetRemoveWorkspacesOptions{
		Workspaces: workspaces,
	}

	if err := c.policySetService(client).RemoveWorkspaces(client.Context(), c.policySetID, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing workspaces from policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Removed %d workspace(s) from policy set '%s'", len(ids), c.policySetID))
	return 0
}

func (c *PolicySetRemoveWorkspaceCommand) policySetService(client *client.Client) policySetWorkspaceRemover {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Help returns help text for policyset remove-workspace.
func (c *PolicySetRemoveWorkspaceCommand) Help() string {
	helpText := `
Usage: hcptf policyset remove-workspace [options]

  Remove workspaces from a policy set.

Options:

  -policyset-id=<id>      Policy set ID (required)
  -workspace-ids=<ids>    Comma-separated workspace IDs to remove (required)

Example:

  hcptf policyset remove-workspace -policyset-id=polset-12345 -workspace-ids=ws-abc123,ws-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for policyset remove-workspace.
func (c *PolicySetRemoveWorkspaceCommand) Synopsis() string {
	return "Remove workspaces from a policy set"
}
