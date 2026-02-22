package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetListWorkspaceCommand lists variable sets for a workspace.
type VariableSetListWorkspaceCommand struct {
	Meta
	workspaceID string
	query       string
	include     string
	format      string
	varSetSvc   variableSetWorkspaceLister
}

// Run executes the variableset list-workspace command.
func (c *VariableSetListWorkspaceCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset list-workspace")
	flags.StringVar(&c.workspaceID, "workspace-id", "", "Workspace ID (required)")
	flags.StringVar(&c.query, "query", "", "Filter variable sets by name query")
	flags.StringVar(&c.include, "include", "", "Include related resources (comma-separated)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.workspaceID == "" {
		c.Ui.Error("Error: -workspace-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	options := &tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{PageSize: 100},
		Query:       c.query,
	}
	if c.include != "" {
		options.Include = strings.Join(splitCommaList(c.include), ",")
	}

	variableSets, err := c.varSetService(client).ListForWorkspace(client.Context(), c.workspaceID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing variable sets for workspace: %s", err))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	if len(variableSets.Items) == 0 {
		c.Ui.Output("No variable sets found")
		return 0
	}

	headers := []string{"ID", "Name", "Description", "Global", "Priority"}
	var rows [][]string
	for _, vs := range variableSets.Items {
		rows = append(rows, []string{
			vs.ID,
			vs.Name,
			vs.Description,
			fmt.Sprintf("%t", vs.Global),
			fmt.Sprintf("%t", vs.Priority),
		})
	}

	formatter.Table(headers, rows)
	return 0
}

func (c *VariableSetListWorkspaceCommand) varSetService(client *client.Client) variableSetWorkspaceLister {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Help returns help text for variableset list-workspace.
func (c *VariableSetListWorkspaceCommand) Help() string {
	helpText := `
Usage: hcptf variableset list-workspace [options]

  List variable sets associated with a workspace.

Options:

  -workspace-id=<id>   Workspace ID (required)
  -query=<query>       Filter by variable set name query
  -include=<values>    Include related resources (comma-separated)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variableset list-workspace -workspace-id=ws-abc123
  hcptf variableset list-workspace -workspace-id=ws-abc123 -include=vars,projects
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for variableset list-workspace.
func (c *VariableSetListWorkspaceCommand) Synopsis() string {
	return "List variable sets for a workspace"
}
