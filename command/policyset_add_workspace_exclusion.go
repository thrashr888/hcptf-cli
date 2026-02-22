package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetAddWorkspaceExclusionCommand adds workspace exclusions to a policy set.
type PolicySetAddWorkspaceExclusionCommand struct {
	Meta
	policySetID  string
	workspaceIDs string
	policySetSvc policySetWorkspaceExclusionAdder
}

// Run executes the policyset add-workspace-exclusion command.
func (c *PolicySetAddWorkspaceExclusionCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset add-workspace-exclusion")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.workspaceIDs, "workspace-ids", "", "Comma-separated workspace IDs to exclude (required)")

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
	options := tfe.PolicySetAddWorkspaceExclusionsOptions{
		WorkspaceExclusions: workspaces,
	}

	if err := c.policySetService(client).AddWorkspaceExclusions(client.Context(), c.policySetID, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Error adding workspace exclusions to policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Added %d workspace exclusion(s) to policy set '%s'", len(ids), c.policySetID))
	return 0
}

func (c *PolicySetAddWorkspaceExclusionCommand) policySetService(client *client.Client) policySetWorkspaceExclusionAdder {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Help returns help text for policyset add-workspace-exclusion.
func (c *PolicySetAddWorkspaceExclusionCommand) Help() string {
	helpText := `
Usage: hcptf policyset add-workspace-exclusion [options]

  Add workspace exclusions to a policy set.

Options:

  -policyset-id=<id>      Policy set ID (required)
  -workspace-ids=<ids>    Comma-separated workspace IDs to exclude (required)

Example:

  hcptf policyset add-workspace-exclusion -policyset-id=polset-12345 -workspace-ids=ws-abc123,ws-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for policyset add-workspace-exclusion.
func (c *PolicySetAddWorkspaceExclusionCommand) Synopsis() string {
	return "Add workspace exclusions to a policy set"
}
