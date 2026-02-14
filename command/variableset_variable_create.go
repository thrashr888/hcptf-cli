package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// VariableSetVariableCreateCommand is a command to add a variable to a variable set
type VariableSetVariableCreateCommand struct {
	Meta
	variableSetID string
	key           string
	value         string
	category      string
	sensitive     bool
	hcl           bool
	description   string
	format        string
}

// Run executes the variable set variable create command
func (c *VariableSetVariableCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset variable create")
	flags.StringVar(&c.variableSetID, "variableset-id", "", "Variable set ID (required)")
	flags.StringVar(&c.key, "key", "", "Variable key/name (required)")
	flags.StringVar(&c.value, "value", "", "Variable value (required)")
	flags.StringVar(&c.category, "category", "terraform", "Variable category: terraform or env")
	flags.BoolVar(&c.sensitive, "sensitive", false, "Mark variable as sensitive")
	flags.BoolVar(&c.hcl, "hcl", false, "Parse variable as HCL")
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

	if c.key == "" {
		c.Ui.Error("Error: -key flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.value == "" && !c.sensitive {
		c.Ui.Error("Error: -value flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate category
	if c.category != "terraform" && c.category != "env" {
		c.Ui.Error("Error: -category must be 'terraform' or 'env'")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create variable
	options := tfe.VariableSetVariableCreateOptions{
		Key:       tfe.String(c.key),
		Value:     tfe.String(c.value),
		Category:  tfe.Category(tfe.CategoryType(c.category)),
		Sensitive: tfe.Bool(c.sensitive),
		HCL:       tfe.Bool(c.hcl),
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	variable, err := client.VariableSetVariables.Create(client.Context(), c.variableSetID, &options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating variable: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Variable '%s' created successfully in variable set", variable.Key))

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

// Help returns help text for the variable set variable create command
func (c *VariableSetVariableCreateCommand) Help() string {
	helpText := `
Usage: hcptf variableset variable create [options]

  Add a variable to a variable set.

Options:

  -variableset-id=<id>  Variable set ID (required)
  -key=<key>            Variable key/name (required)
  -value=<value>        Variable value (required unless sensitive)
  -category=<type>      Variable category: terraform (default) or env
  -sensitive            Mark variable as sensitive (default: false)
  -hcl                  Parse variable as HCL (default: false)
  -description=<text>   Variable description
  -output=<format>      Output format: table (default) or json

Example:

  hcptf variableset variable create -variableset-id=varset-12345 -key=region -value=us-east-1
  hcptf variableset variable create -variableset-id=varset-12345 -key=password -value=secret -sensitive
  hcptf variableset variable create -variableset-id=varset-12345 -key=PATH -value=/usr/bin -category=env
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable set variable create command
func (c *VariableSetVariableCreateCommand) Synopsis() string {
	return "Add a variable to a variable set"
}
