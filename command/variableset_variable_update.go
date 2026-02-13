package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetVariableUpdateCommand is a command to update a variable in a variable set
type VariableSetVariableUpdateCommand struct {
	Meta
	variableSetID string
	variableID    string
	key           string
	value         string
	sensitive     string
	hcl           string
	description   string
	format        string
}

// Run executes the variable set variable update command
func (c *VariableSetVariableUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset variable update")
	flags.StringVar(&c.variableSetID, "variableset-id", "", "Variable set ID (required)")
	flags.StringVar(&c.variableID, "variable-id", "", "Variable ID (required)")
	flags.StringVar(&c.key, "key", "", "Variable key/name")
	flags.StringVar(&c.value, "value", "", "Variable value")
	flags.StringVar(&c.sensitive, "sensitive", "", "Mark variable as sensitive (true or false)")
	flags.StringVar(&c.hcl, "hcl", "", "Parse variable as HCL (true or false)")
	flags.StringVar(&c.description, "description", "", "Variable description")
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

	if c.variableID == "" {
		c.Ui.Error("Error: -variable-id flag is required")
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
	options := tfe.VariableSetVariableUpdateOptions{}

	if c.key != "" {
		options.Key = tfe.String(c.key)
	}

	if c.value != "" {
		options.Value = tfe.String(c.value)
	}

	if c.sensitive != "" {
		if c.sensitive == "true" {
			options.Sensitive = tfe.Bool(true)
		} else if c.sensitive == "false" {
			options.Sensitive = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -sensitive must be 'true' or 'false'")
			c.Ui.Error(c.Help())
			return 1
		}
	}

	if c.hcl != "" {
		if c.hcl == "true" {
			options.HCL = tfe.Bool(true)
		} else if c.hcl == "false" {
			options.HCL = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -hcl must be 'true' or 'false'")
			c.Ui.Error(c.Help())
			return 1
		}
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	// Update variable
	variable, err := client.VariableSetVariables.Update(client.Context(), c.variableSetID, c.variableID, &options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating variable: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Variable '%s' updated successfully", variable.Key))

	// Show variable details
	displayValue := variable.Value
	if variable.Sensitive {
		displayValue = "(sensitive)"
	}

	data := map[string]interface{}{
		"ID":          variable.ID,
		"Key":         variable.Key,
		"Value":       displayValue,
		"Category":    variable.Category,
		"Sensitive":   variable.Sensitive,
		"HCL":         variable.HCL,
		"Description": variable.Description,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the variable set variable update command
func (c *VariableSetVariableUpdateCommand) Help() string {
	helpText := `
Usage: hcptf variableset variable update [options]

  Update a variable in a variable set.

Options:

  -variableset-id=<id>  Variable set ID (required)
  -variable-id=<id>     Variable ID (required)
  -key=<key>            Variable key/name
  -value=<value>        Variable value
  -sensitive=<bool>     Mark variable as sensitive (true or false)
  -hcl=<bool>           Parse variable as HCL (true or false)
  -description=<text>   Variable description
  -output=<format>      Output format: table (default) or json

Note: Variable category (terraform/env) cannot be changed after creation.

Example:

  hcptf variableset variable update -variableset-id=varset-12345 -variable-id=var-abc123 -value=us-west-2
  hcptf variableset variable update -variableset-id=varset-12345 -variable-id=var-abc123 -sensitive=true
  hcptf variableset variable update -variableset-id=varset-12345 -variable-id=var-abc123 -key=new_name -description="Updated variable"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set variable update command
func (c *VariableSetVariableUpdateCommand) Synopsis() string {
	return "Update a variable in a variable set"
}
