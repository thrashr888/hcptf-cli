package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableListCommand is a command to list variables
type VariableListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

// Run executes the variable list command
func (c *VariableListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variable list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace first
	ws, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// List variables
	variables, err := client.Variables.List(client.Context(), ws.ID, &tfe.VariableListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing variables: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(variables.Items) == 0 {
		c.Ui.Output("No variables found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Key", "Value", "Category", "Sensitive", "HCL"}
	var rows [][]string
	var fullRows [][]string

	for _, v := range variables.Items {
		value := v.Value
		if v.Sensitive {
			value = "(sensitive)"
		}

		hcl := "false"
		if v.HCL {
			hcl = "true"
		}

		sensitive := "false"
		if v.Sensitive {
			sensitive = "true"
		}

		fullRows = append(fullRows, []string{
			v.ID,
			v.Key,
			value,
			string(v.Category),
			sensitive,
			hcl,
		})

		// Truncate display value for table output only
		displayValue := value
		if len(displayValue) > 50 {
			displayValue = displayValue[:47] + "..."
		}

		rows = append(rows, []string{
			v.ID,
			v.Key,
			displayValue,
			string(v.Category),
			sensitive,
			hcl,
		})
	}

	formatter.TableWithFullRows(headers, rows, fullRows)
	return 0
}

// Help returns help text for the variable list command
func (c *VariableListCommand) Help() string {
	helpText := `
Usage: hcptf variable list [options]

  List variables for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variable list -org=my-org -workspace=my-workspace
  hcptf variable list -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable list command
func (c *VariableListCommand) Synopsis() string {
	return "List variables for a workspace"
}
