package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetUpdateWorkspacesCommand synchronizes workspace associations.
type VariableSetUpdateWorkspacesCommand struct {
	Meta
	id         string
	workspaces string
	format     string
	varSetSvc  variableSetWorkspaceUpdater
}

// Run executes the variableset update-workspaces command.
func (c *VariableSetUpdateWorkspacesCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset update-workspaces")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
	flags.StringVar(&c.workspaces, "workspaces", "", "Comma-separated workspace IDs (empty clears all)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	workspaces := make([]*tfe.Workspace, 0)
	for _, workspaceID := range splitCommaList(c.workspaces) {
		if strings.TrimSpace(workspaceID) == "" {
			continue
		}
		workspaces = append(workspaces, &tfe.Workspace{ID: strings.TrimSpace(workspaceID)})
	}

	options := &tfe.VariableSetUpdateWorkspacesOptions{
		Workspaces: workspaces,
	}
	updated, err := c.varSetService(client).UpdateWorkspaces(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating variable set workspaces: %s", err))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Variable set '%s' workspace associations updated", updated.Name))
	}
	formatter.KeyValue(map[string]interface{}{
		"ID":             updated.ID,
		"Name":           updated.Name,
		"WorkspaceCount": len(updated.Workspaces),
	})
	return 0
}

func (c *VariableSetUpdateWorkspacesCommand) varSetService(client *client.Client) variableSetWorkspaceUpdater {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Help returns help text for variableset update-workspaces.
func (c *VariableSetUpdateWorkspacesCommand) Help() string {
	helpText := `
Usage: hcptf variableset update-workspaces [options]

  Replace workspace associations for a variable set.

Options:

  -id=<id>            Variable set ID (required)
  -workspaces=<ids>   Comma-separated workspace IDs (empty clears all)
  -output=<format>    Output format: table (default) or json

Example:

  hcptf variableset update-workspaces -id=varset-12345 -workspaces=ws-abc123,ws-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for variableset update-workspaces.
func (c *VariableSetUpdateWorkspacesCommand) Synopsis() string {
	return "Replace variable set workspace associations"
}
