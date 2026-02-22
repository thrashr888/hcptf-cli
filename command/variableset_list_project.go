package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetListProjectCommand lists variable sets for a project.
type VariableSetListProjectCommand struct {
	Meta
	projectID string
	query     string
	include   string
	format    string
	varSetSvc variableSetProjectLister
}

// Run executes the variableset list-project command.
func (c *VariableSetListProjectCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset list-project")
	flags.StringVar(&c.projectID, "project-id", "", "Project ID (required)")
	flags.StringVar(&c.query, "query", "", "Filter variable sets by name query")
	flags.StringVar(&c.include, "include", "", "Include related resources (comma-separated)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.projectID == "" {
		c.Ui.Error("Error: -project-id flag is required")
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

	variableSets, err := c.varSetService(client).ListForProject(client.Context(), c.projectID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing variable sets for project: %s", err))
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

func (c *VariableSetListProjectCommand) varSetService(client *client.Client) variableSetProjectLister {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Help returns help text for variableset list-project.
func (c *VariableSetListProjectCommand) Help() string {
	helpText := `
Usage: hcptf variableset list-project [options]

  List variable sets associated with a project.

Options:

  -project-id=<id>     Project ID (required)
  -query=<query>       Filter by variable set name query
  -include=<values>    Include related resources (comma-separated)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variableset list-project -project-id=prj-abc123
  hcptf variableset list-project -project-id=prj-abc123 -include=vars,workspaces
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for variableset list-project.
func (c *VariableSetListProjectCommand) Synopsis() string {
	return "List variable sets for a project"
}
