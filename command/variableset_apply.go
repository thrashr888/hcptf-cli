package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetApplyCommand is a command to apply a variable set to workspaces or projects
type VariableSetApplyCommand struct {
	Meta
	id         string
	workspaces string
	projects   string
	stacks     string
}

type variableSetApplyPayload struct {
	Workspaces []string `json:"workspaces"`
	Projects   []string `json:"projects"`
	Stacks     []string `json:"stacks"`
}

// Run executes the variable set apply command
func (c *VariableSetApplyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset apply")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
	flags.StringVar(&c.workspaces, "workspaces", "", "Comma-separated list of workspace IDs to apply to")
	flags.StringVar(&c.projects, "projects", "", "Comma-separated list of project IDs to apply to")
	flags.StringVar(&c.stacks, "stacks", "", "Comma-separated list of stack IDs to apply to")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspaces == "" && c.projects == "" && c.stacks == "" {
		if c.Meta.JSONInput == "" {
			c.Ui.Error("Error: at least one of -workspaces, -projects, or -stacks is required")
			c.Ui.Error(c.Help())
			return 1
		}
	}

	if !c.Meta.ValidateID(c.id, "-id") {
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	payload := variableSetApplyPayload{
		Workspaces: splitCommaList(c.workspaces),
		Projects:   splitCommaList(c.projects),
		Stacks:     splitCommaList(c.stacks),
	}
	if c.Meta.JSONInput != "" {
		if err := c.Meta.ParseJSONInput(&payload); err != nil {
			c.Ui.Error(fmt.Sprintf("Error parsing JSON input: %s", err))
			return 1
		}
	}

	if len(payload.Workspaces) == 0 && len(payload.Projects) == 0 && len(payload.Stacks) == 0 {
		c.Ui.Error("Error: at least one of workspaces, projects, or stacks is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.Meta.DryRun {
		formatter := c.Meta.NewFormatter("json")
		formatter.JSON(map[string]interface{}{
			"action":   "apply",
			"resource": "variableset",
			"id":       c.id,
			"options":  payload,
		})
		return 0
	}

	// Apply to workspaces
	if len(payload.Workspaces) > 0 {
		workspaceIDs := payload.Workspaces
		workspaces := make([]*tfe.Workspace, 0, len(workspaceIDs))
		for _, wsID := range workspaceIDs {
			workspaces = append(workspaces, &tfe.Workspace{ID: strings.TrimSpace(wsID)})
		}

		options := tfe.VariableSetApplyToWorkspacesOptions{
			Workspaces: workspaces,
		}

		err = client.VariableSets.ApplyToWorkspaces(client.Context(), c.id, &options)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error applying variable set to workspaces: %s", err))
			return 1
		}

		c.Ui.Output(fmt.Sprintf("Variable set '%s' applied to %d workspace(s) successfully", c.id, len(workspaceIDs)))
	}

	// Apply to projects
	if len(payload.Projects) > 0 {
		projectIDs := payload.Projects
		projects := make([]*tfe.Project, 0, len(projectIDs))
		for _, projID := range projectIDs {
			projects = append(projects, &tfe.Project{ID: strings.TrimSpace(projID)})
		}

		options := tfe.VariableSetApplyToProjectsOptions{
			Projects: projects,
		}

		err = client.VariableSets.ApplyToProjects(client.Context(), c.id, options)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error applying variable set to projects: %s", err))
			return 1
		}

		c.Ui.Output(fmt.Sprintf("Variable set '%s' applied to %d project(s) successfully", c.id, len(projectIDs)))
	}

	// Apply to stacks
	if len(payload.Stacks) > 0 {
		stackIDs := payload.Stacks
		stacks := make([]*tfe.Stack, 0, len(stackIDs))
		for _, stackID := range stackIDs {
			stacks = append(stacks, &tfe.Stack{ID: strings.TrimSpace(stackID)})
		}

		options := tfe.VariableSetApplyToStacksOptions{
			Stacks: stacks,
		}

		err = client.VariableSets.ApplyToStacks(client.Context(), c.id, &options)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error applying variable set to stacks: %s", err))
			return 1
		}

		c.Ui.Output(fmt.Sprintf("Variable set '%s' applied to %d stack(s) successfully", c.id, len(stackIDs)))
	}

	return 0
}

// Help returns help text for the variable set apply command
func (c *VariableSetApplyCommand) Help() string {
	helpText := `
Usage: hcptf variableset apply [options]

  Apply a variable set to specific workspaces or projects.

Options:

  -id=<id>                  Variable set ID (required)
  -workspaces=<ids>         Comma-separated list of workspace IDs
  -projects=<ids>           Comma-separated list of project IDs
  -stacks=<ids>             Comma-separated list of stack IDs

Example:

  hcptf variableset apply -id=varset-12345 -workspaces=ws-abc123,ws-def456
  hcptf variableset apply -id=varset-12345 -projects=prj-abc123
  hcptf variableset apply -id=varset-12345 -stacks=stack-abc123
  hcptf variableset apply -id=varset-12345 -workspaces=ws-abc123 -projects=prj-abc123 -stacks=stack-abc123
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set apply command
func (c *VariableSetApplyCommand) Synopsis() string {
	return "Apply a variable set to workspaces, projects, or stacks"
}
