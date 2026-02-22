package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetUpdateCommand is a command to update a variable set
type VariableSetUpdateCommand struct {
	Meta
	id          string
	name        string
	description string
	global      string
	priority    string
	format      string
}

// Run executes the variable set update command
func (c *VariableSetUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset update")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
	flags.StringVar(&c.name, "name", "", "Variable set name")
	flags.StringVar(&c.description, "description", "", "Variable set description")
	flags.StringVar(&c.global, "global", "", "Apply to all workspaces (true or false)")
	flags.StringVar(&c.priority, "priority", "", "Variable set priority override (true/false)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build update options
	options := tfe.VariableSetUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	if c.global != "" {
		if c.global == "true" {
			options.Global = tfe.Bool(true)
		} else if c.global == "false" {
			options.Global = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -global must be 'true' or 'false'")
			c.Ui.Error(c.Help())
			return 1
		}
	}

	if c.priority != "" {
		priority, parseErr := parseBoolFlag(c.priority, "priority")
		if parseErr != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
			c.Ui.Error(c.Help())
			return 1
		}
		options.Priority = priority
	}

	// Update variable set
	variableSet, err := client.VariableSets.Update(client.Context(), c.id, &options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating variable set: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Variable set '%s' updated successfully", variableSet.Name))

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

// Help returns help text for the variable set update command
func (c *VariableSetUpdateCommand) Help() string {
	helpText := `
Usage: hcptf variableset update [options]

  Update a variable set's settings.

Options:

  -id=<id>             Variable set ID (required)
  -name=<name>         New variable set name
  -description=<text>  New variable set description
  -global=<bool>       Apply to all workspaces (true or false)
  -priority=<bool>     Override lower-scope variable values (true/false)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variableset update -id=varset-12345 -name=new-name
  hcptf variableset update -id=varset-12345 -priority=true
  hcptf variableset update -id=varset-12345 -global=true -description="Global variables"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set update command
func (c *VariableSetUpdateCommand) Synopsis() string {
	return "Update a variable set's settings"
}
