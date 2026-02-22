package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetRemoveProjectCommand removes projects from a policy set.
type PolicySetRemoveProjectCommand struct {
	Meta
	policySetID  string
	projectIDs   string
	policySetSvc policySetProjectRemover
}

// Run executes the policyset remove-project command.
func (c *PolicySetRemoveProjectCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset remove-project")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.projectIDs, "project-ids", "", "Comma-separated project IDs to remove (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.policySetID == "" {
		c.Ui.Error("Error: -policyset-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.projectIDs == "" {
		c.Ui.Error("Error: -project-ids flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	ids := splitCommaList(c.projectIDs)
	projects := make([]*tfe.Project, 0, len(ids))
	for _, id := range ids {
		projects = append(projects, &tfe.Project{ID: strings.TrimSpace(id)})
	}
	options := tfe.PolicySetRemoveProjectsOptions{
		Projects: projects,
	}

	if err := c.policySetService(client).RemoveProjects(client.Context(), c.policySetID, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Error removing projects from policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Removed %d project(s) from policy set '%s'", len(ids), c.policySetID))
	return 0
}

func (c *PolicySetRemoveProjectCommand) policySetService(client *client.Client) policySetProjectRemover {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Help returns help text for policyset remove-project.
func (c *PolicySetRemoveProjectCommand) Help() string {
	helpText := `
Usage: hcptf policyset remove-project [options]

  Remove projects from a policy set.

Options:

  -policyset-id=<id>   Policy set ID (required)
  -project-ids=<ids>   Comma-separated project IDs to remove (required)

Example:

  hcptf policyset remove-project -policyset-id=polset-12345 -project-ids=prj-abc123,prj-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for policyset remove-project.
func (c *PolicySetRemoveProjectCommand) Synopsis() string {
	return "Remove projects from a policy set"
}
