package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// VariableCreateCommand is a command to create a variable
type VariableCreateCommand struct {
	Meta
	organization string
	workspace    string
	key          string
	value        string
	category     string
	sensitive    bool
	hcl          bool
	description  string
	format       string
	workspaceSvc workspaceReader
	variableSvc  variableCreator
}

// Run executes the variable create command
func (c *VariableCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variable create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.key, "key", "", "Variable key (required)")
	flags.StringVar(&c.value, "value", "", "Variable value (required)")
	flags.StringVar(&c.category, "category", "terraform", "Variable category: terraform or env (default: terraform)")
	flags.BoolVar(&c.sensitive, "sensitive", false, "Mark variable as sensitive")
	flags.BoolVar(&c.hcl, "hcl", false, "Parse value as HCL")
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

	if c.key == "" {
		c.Ui.Error("Error: -key flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.value == "" {
		c.Ui.Error("Error: -value flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate category
	var category tfe.CategoryType
	if c.category == "terraform" {
		category = tfe.CategoryTerraform
	} else if c.category == "env" {
		category = tfe.CategoryEnv
	} else {
		c.Ui.Error("Error: -category must be 'terraform' or 'env'")
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

	// Create variable
	options := tfe.VariableCreateOptions{
		Key:       tfe.String(c.key),
		Value:     tfe.String(c.value),
		Category:  &category,
		Sensitive: tfe.Bool(c.sensitive),
		HCL:       tfe.Bool(c.hcl),
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	variable, err := c.variableService(client).Create(client.Context(), ws.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating variable: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Variable '%s' created successfully", variable.Key))

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

func (c *VariableCreateCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

func (c *VariableCreateCommand) variableService(client *client.Client) variableCreator {
	if c.variableSvc != nil {
		return c.variableSvc
	}
	return client.Variables
}

// Help returns help text for the variable create command
func (c *VariableCreateCommand) Help() string {
	helpText := `
Usage: hcptf variable create [options]

  Create a new variable for a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -key=<name>          Variable key (required)
  -value=<value>       Variable value (required)
  -category=<type>     Variable category: terraform or env (default: terraform)
  -sensitive           Mark variable as sensitive
  -hcl                 Parse value as HCL
  -description=<text>  Variable description
  -output=<format>     Output format: table (default) or json

Example:

  hcptf variable create -org=my-org -workspace=prod -key=region -value=us-east-1
  hcptf variable create -org=my-org -workspace=prod -key=AWS_ACCESS_KEY_ID -value=AKIAIOSFODNN7EXAMPLE -category=env -sensitive
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the variable create command
func (c *VariableCreateCommand) Synopsis() string {
	return "Create a new variable"
}
