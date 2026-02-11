package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VariableSetVariableListCommand is a command to list variables in a variable set
type VariableSetVariableListCommand struct {
	Meta
	variableSetID string
	format        string
}

// Run executes the variable set variable list command
func (c *VariableSetVariableListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset variable list")
	flags.StringVar(&c.variableSetID, "variableset-id", "", "Variable set ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.variableSetID == "" {
		c.Ui.Error("Error: -variableset-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List variables in variable set
	variables, err := client.VariableSetVariables.List(client.Context(), c.variableSetID, &tfe.VariableSetVariableListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing variables: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(variables.Items) == 0 {
		c.Ui.Output("No variables found in variable set")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Key", "Value", "Category", "Sensitive", "HCL"}
	var rows [][]string

	for _, v := range variables.Items {
		value := v.Value
		if v.Sensitive {
			value = "(sensitive)"
		}
		if len(value) > 50 {
			value = value[:47] + "..."
		}

		hcl := "false"
		if v.HCL {
			hcl = "true"
		}

		sensitive := "false"
		if v.Sensitive {
			sensitive = "true"
		}

		rows = append(rows, []string{
			v.ID,
			v.Key,
			value,
			string(v.Category),
			sensitive,
			hcl,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the variable set variable list command
func (c *VariableSetVariableListCommand) Help() string {
	helpText := `
Usage: hcptf variableset variable list [options]

  List variables in a variable set.

Options:

  -variableset-id=<id>  Variable set ID (required)
  -output=<format>      Output format: table (default) or json

Example:

  hcptf variableset variable list -variableset-id=varset-12345
  hcptf variableset variable list -variableset-id=varset-12345 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set variable list command
func (c *VariableSetVariableListCommand) Synopsis() string {
	return "List variables in a variable set"
}
