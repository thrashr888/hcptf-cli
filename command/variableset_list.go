package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetListCommand is a command to list variable sets
type VariableSetListCommand struct {
	Meta
	organization string
	format       string
	varSetSvc    variableSetLister
}

// Run executes the variable set list command
func (c *VariableSetListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List variable sets
	variableSets, err := c.varSetService(client).List(client.Context(), c.organization, &tfe.VariableSetListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing variable sets: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if len(variableSets.Items) == 0 {
		c.Ui.Output("No variable sets found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Name", "Description", "Global", "Variables"}
	var rows [][]string

	for _, vs := range variableSets.Items {
		global := "false"
		if vs.Global {
			global = "true"
		}

		description := vs.Description
		if len(description) > 50 {
			description = description[:47] + "..."
		}

		variableCount := "0"
		if len(vs.Variables) > 0 {
			variableCount = fmt.Sprintf("%d", len(vs.Variables))
		}

		rows = append(rows, []string{
			vs.ID,
			vs.Name,
			description,
			global,
			variableCount,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the variable set list command
func (c *VariableSetListCommand) Help() string {
	helpText := `
Usage: hcptf variableset list [options]

  List variable sets in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variableset list -org=my-org
  hcptf variableset list -org=my-org -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *VariableSetListCommand) varSetService(client *client.Client) variableSetLister {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Synopsis returns a short synopsis for the variable set list command
func (c *VariableSetListCommand) Synopsis() string {
	return "List variable sets in an organization"
}
