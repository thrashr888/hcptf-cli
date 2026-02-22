package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetAddWorkspaceCommand adds workspaces to a policy set.
type PolicySetAddWorkspaceCommand struct {
	Meta
	policySetID  string
	workspaceIDs string
	policySetSvc policySetWorkspaceAdder
}

// Run executes the policyset add-workspace command.
func (c *PolicySetAddWorkspaceCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset add-workspace")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.workspaceIDs, "workspace-ids", "", "Comma-separated workspace IDs to add (required)")

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
	options := tfe.PolicySetAddWorkspacesOptions{
		Workspaces: workspaces,
	}

	if err := c.policySetService(client).AddWorkspaces(client.Context(), c.policySetID, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Error adding workspaces to policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Added %d workspace(s) to policy set '%s'", len(ids), c.policySetID))
	return 0
}

func (c *PolicySetAddWorkspaceCommand) policySetService(client *client.Client) policySetWorkspaceAdder {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Help returns help text for policyset add-workspace.
func (c *PolicySetAddWorkspaceCommand) Help() string {
	helpText := `
Usage: hcptf policyset add-workspace [options]

  Add workspaces to a policy set.

Options:

  -policyset-id=<id>      Policy set ID (required)
  -workspace-ids=<ids>    Comma-separated workspace IDs to add (required)

Example:

  hcptf policyset add-workspace -policyset-id=polset-12345 -workspace-ids=ws-abc123,ws-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for policyset add-workspace.
func (c *PolicySetAddWorkspaceCommand) Synopsis() string {
	return "Add workspaces to a policy set"
}
