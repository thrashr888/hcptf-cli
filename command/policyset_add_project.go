package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// PolicySetAddProjectCommand adds projects to a policy set.
type PolicySetAddProjectCommand struct {
	Meta
	policySetID  string
	projectIDs   string
	policySetSvc policySetProjectAdder
}

// Run executes the policyset add-project command.
func (c *PolicySetAddProjectCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset add-project")
	flags.StringVar(&c.policySetID, "policyset-id", "", "Policy set ID (required)")
	flags.StringVar(&c.projectIDs, "project-ids", "", "Comma-separated project IDs to add (required)")

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
	options := tfe.PolicySetAddProjectsOptions{
		Projects: projects,
	}

	if err := c.policySetService(client).AddProjects(client.Context(), c.policySetID, options); err != nil {
		c.Ui.Error(fmt.Sprintf("Error adding projects to policy set: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Added %d project(s) to policy set '%s'", len(ids), c.policySetID))
	return 0
}

func (c *PolicySetAddProjectCommand) policySetService(client *client.Client) policySetProjectAdder {
	if c.policySetSvc != nil {
		return c.policySetSvc
	}
	return client.PolicySets
}

// Help returns help text for policyset add-project.
func (c *PolicySetAddProjectCommand) Help() string {
	helpText := `
Usage: hcptf policyset add-project [options]

  Add projects to a policy set.

Options:

  -policyset-id=<id>   Policy set ID (required)
  -project-ids=<ids>   Comma-separated project IDs to add (required)

Example:

  hcptf policyset add-project -policyset-id=polset-12345 -project-ids=prj-abc123,prj-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for policyset add-project.
func (c *PolicySetAddProjectCommand) Synopsis() string {
	return "Add projects to a policy set"
}
