package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetRemoveCommand removes variable set associations from resources.
type VariableSetRemoveCommand struct {
	Meta
	id         string
	workspaces string
	projects   string
	stacks     string
	varSetSvc  variableSetRemover
}

// Run executes the variableset remove command.
func (c *VariableSetRemoveCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset remove")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
	flags.StringVar(&c.workspaces, "workspaces", "", "Comma-separated workspace IDs to remove")
	flags.StringVar(&c.projects, "projects", "", "Comma-separated project IDs to remove")
	flags.StringVar(&c.stacks, "stacks", "", "Comma-separated stack IDs to remove")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}
	if c.workspaces == "" && c.projects == "" && c.stacks == "" {
		c.Ui.Error("Error: at least one of -workspaces, -projects, or -stacks is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	if c.workspaces != "" {
		workspaceIDs := splitCommaList(c.workspaces)
		workspaces := make([]*tfe.Workspace, 0, len(workspaceIDs))
		for _, workspaceID := range workspaceIDs {
			workspaces = append(workspaces, &tfe.Workspace{ID: strings.TrimSpace(workspaceID)})
		}
		options := tfe.VariableSetRemoveFromWorkspacesOptions{
			Workspaces: workspaces,
		}
		if err := c.varSetService(client).RemoveFromWorkspaces(client.Context(), c.id, &options); err != nil {
			c.Ui.Error(fmt.Sprintf("Error removing variable set from workspaces: %s", err))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Variable set '%s' removed from %d workspace(s)", c.id, len(workspaceIDs)))
	}

	if c.projects != "" {
		projectIDs := splitCommaList(c.projects)
		projects := make([]*tfe.Project, 0, len(projectIDs))
		for _, projectID := range projectIDs {
			projects = append(projects, &tfe.Project{ID: strings.TrimSpace(projectID)})
		}
		options := tfe.VariableSetRemoveFromProjectsOptions{
			Projects: projects,
		}
		if err := c.varSetService(client).RemoveFromProjects(client.Context(), c.id, options); err != nil {
			c.Ui.Error(fmt.Sprintf("Error removing variable set from projects: %s", err))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Variable set '%s' removed from %d project(s)", c.id, len(projectIDs)))
	}

	if c.stacks != "" {
		stackIDs := splitCommaList(c.stacks)
		stacks := make([]*tfe.Stack, 0, len(stackIDs))
		for _, stackID := range stackIDs {
			stacks = append(stacks, &tfe.Stack{ID: strings.TrimSpace(stackID)})
		}
		options := tfe.VariableSetRemoveFromStacksOptions{
			Stacks: stacks,
		}
		if err := c.varSetService(client).RemoveFromStacks(client.Context(), c.id, &options); err != nil {
			c.Ui.Error(fmt.Sprintf("Error removing variable set from stacks: %s", err))
			return 1
		}
		c.Ui.Output(fmt.Sprintf("Variable set '%s' removed from %d stack(s)", c.id, len(stackIDs)))
	}

	return 0
}

func (c *VariableSetRemoveCommand) varSetService(client *client.Client) variableSetRemover {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Help returns help text for the variableset remove command.
func (c *VariableSetRemoveCommand) Help() string {
	helpText := `
Usage: hcptf variableset remove [options]

  Remove a variable set from workspaces, projects, or stacks.

Options:

  -id=<id>            Variable set ID (required)
  -workspaces=<ids>   Comma-separated workspace IDs
  -projects=<ids>     Comma-separated project IDs
  -stacks=<ids>       Comma-separated stack IDs

Example:

  hcptf variableset remove -id=varset-12345 -workspaces=ws-abc123
  hcptf variableset remove -id=varset-12345 -projects=prj-123,prj-456
  hcptf variableset remove -id=varset-12345 -stacks=stack-123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variableset remove command.
func (c *VariableSetRemoveCommand) Synopsis() string {
	return "Remove a variable set from resources"
}
