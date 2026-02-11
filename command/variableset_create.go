package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VariableSetCreateCommand is a command to create a variable set
type VariableSetCreateCommand struct {
	Meta
	organization string
	name         string
	description  string
	global       bool
	format       string
}

// Run executes the variable set create command
func (c *VariableSetCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Variable set name (required)")
	flags.StringVar(&c.description, "description", "", "Variable set description")
	flags.BoolVar(&c.global, "global", false, "Apply to all workspaces in the organization")
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create variable set
	options := tfe.VariableSetCreateOptions{
		Name:   tfe.String(c.name),
		Global: tfe.Bool(c.global),
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	variableSet, err := client.VariableSets.Create(client.Context(), c.organization, &options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating variable set: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Variable set '%s' created successfully", variableSet.Name))

	// Show variable set details
	data := map[string]interface{}{
		"ID":          variableSet.ID,
		"Name":        variableSet.Name,
		"Description": variableSet.Description,
		"Global":      variableSet.Global,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the variable set create command
func (c *VariableSetCreateCommand) Help() string {
	helpText := `
Usage: hcptf variableset create [options]

  Create a new variable set.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Variable set name (required)
  -description=<text>  Variable set description
  -global              Apply to all workspaces (default: false)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variableset create -org=my-org -name=my-varset
  hcptf variableset create -org=my-org -name=prod-vars -global -description="Production variables"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set create command
func (c *VariableSetCreateCommand) Synopsis() string {
	return "Create a new variable set"
}
