package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetRemoveWorkspaceExclusionCommand removes workspace exclusions from a policy set.
type PolicySetRemoveWorkspaceExclusionCommand struct {
	Meta
	policySetID  string
	workspaceIDs string
	policySetSvc policySetWorkspaceExclusionRemover
}

// Run executes the policyset remove-workspace-exclusion command.
func (c *PolicySetRemoveWorkspaceExclusionCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset remove-workspace-exclusion")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.workspaceIDs, "workspace-ids", "", "Comma-separated workspace IDs to remove from exclusions (required)")

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
	options := tfe.PolicySetRemoveWorkspaceExclusionsOptions{
		WorkspaceExclusions: workspaces,
	}

	if err := c.policySetService(client).RemoveWorkspaceExclusions(client.Context(), c.policySetID, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing workspace exclusions from policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Removed %d workspace exclusion(s) from policy set '%s'", len(ids), c.policySetID))
	return 0
}

func (c *PolicySetRemoveWorkspaceExclusionCommand) policySetService(client *client.Client) policySetWorkspaceExclusionRemover {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Help returns help text for policyset remove-workspace-exclusion.
func (c *PolicySetRemoveWorkspaceExclusionCommand) Help() string {
	helpText := `
Usage: hcptf policyset remove-workspace-exclusion [options]

  Remove workspace exclusions from a policy set.

Options:

  -policyset-id=<id>      Policy set ID (required)
  -workspace-ids=<ids>    Comma-separated workspace IDs (required)

Example:

  hcptf policyset remove-workspace-exclusion -policyset-id=polset-12345 -workspace-ids=ws-abc123,ws-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for policyset remove-workspace-exclusion.
func (c *PolicySetRemoveWorkspaceExclusionCommand) Synopsis() string {
	return "Remove workspace exclusions from a policy set"
}
