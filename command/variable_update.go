package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VariableUpdateCommand is a command to update a variable
type VariableUpdateCommand struct {
	Meta
	organization string
	workspace    string
	id           string
	key          string
	value        string
	sensitive    string
	hcl          string
	description  string
	format       string
	workspaceSvc workspaceReader
	variableSvc  variableUpdater
}

// Run executes the variable update command
func (c *VariableUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variable update")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.id, "id", "", "Variable ID (required)")
	flags.StringVar(&c.key, "key", "", "Variable key")
	flags.StringVar(&c.value, "value", "", "Variable value")
	flags.StringVar(&c.sensitive, "sensitive", "", "Mark variable as sensitive (true/false)")
	flags.StringVar(&c.hcl, "hcl", "", "Parse value as HCL (true/false)")
	flags.StringVar(&c.description, "description", "", "Variable description")
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

	// Get workspace first
	ws, err := c.workspaceService(client).Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Build update options
	options := tfe.VariableUpdateOptions{}

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
			return 1
		}
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	// Update variable
	variable, err := c.variableService(client).Update(client.Context(), ws.ID, c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating variable: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Variable '%s' updated successfully", variable.Key))

	// Show variable details
	value := variable.Value
	if variable.Sensitive {
		value = "(sensitive)"
	}

	data := map[string]interface{}{
		"ID":          variable.ID,
		"Key":         variable.Key,
		"Value":       value,
		"Category":    variable.Category,
		"Sensitive":   variable.Sensitive,
		"HCL":         variable.HCL,
		"Description": variable.Description,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *VariableUpdateCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *VariableUpdateCommand) variableService(client *client.Client) variableUpdater {
	if c.variableSvc != nil {
		return c.variableSvc
	}
	return client.Variables
}

// Help returns help text for the variable update command
func (c *VariableUpdateCommand) Help() string {
	helpText := `
Usage: hcptf variable update [options]

  Update a variable.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -id=<id>             Variable ID (required)
  -key=<name>          Variable key
  -value=<value>       Variable value
  -sensitive=<bool>    Mark variable as sensitive (true/false)
  -hcl=<bool>          Parse value as HCL (true/false)
  -description=<text>  Variable description
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variable update -org=my-org -workspace=prod -id=var-123 -value=us-west-2
  hcptf variable update -org=my-org -workspace=prod -id=var-456 -sensitive=true
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable update command
func (c *VariableUpdateCommand) Synopsis() string {
	return "Update a variable"
}
