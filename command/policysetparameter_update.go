package command

import (
	"flag"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicySetParameterUpdateCommand is a command to update a policy set parameter
type PolicySetParameterUpdateCommand struct {
	Meta
	policySetID string
	parameterID string
	key         string
	value       string
	sensitive   *bool
	format      string
}

// Run executes the policy set parameter update command
func (c *PolicySetParameterUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetparameter update")
	flags.StringVar(&c.policySetID, "policy-set-id", "", "Policy Set ID (required)")
	flags.StringVar(&c.parameterID, "id", "", "Parameter ID (required)")
	flags.StringVar(&c.key, "key", "", "Parameter key")
	flags.StringVar(&c.value, "value", "", "Parameter value")
	sensitiveFlag := flags.Bool("sensitive", false, "Mark parameter as sensitive")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Check if sensitive flag was actually set
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "sensitive" {
			c.sensitive = sensitiveFlag
		}
	})

	// Validate required flags
	if c.policySetID == "" {
		c.Ui.Error("Error: -policy-set-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.parameterID == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// At least one field must be provided for update
	if c.key == "" && c.value == "" && c.sensitive == nil {
		c.Ui.Error("Error: At least one of -key, -value, or -sensitive must be provided")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Update policy set parameter
	options := tfe.PolicySetParameterUpdateOptions{}

	if c.key != "" {
		options.Key = tfe.String(c.key)
	}

	if c.value != "" {
		options.Value = tfe.String(c.value)
	}

	if c.sensitive != nil {
		options.Sensitive = tfe.Bool(*c.sensitive)
	}

	parameter, err := client.PolicySetParameters.Update(client.Context(), c.policySetID, c.parameterID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating policy set parameter: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Policy set parameter '%s' updated successfully", parameter.Key))

	// Show parameter details
	value := parameter.Value
	if parameter.Sensitive {
		value = "(sensitive)"
	}

	data := map[string]interface{}{
		"ID":        parameter.ID,
		"Key":       parameter.Key,
		"Value":     value,
		"Category":  string(parameter.Category),
		"Sensitive": parameter.Sensitive,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy set parameter update command
func (c *PolicySetParameterUpdateCommand) Help() string {
	helpText := `
Usage: hcptf policysetparameter update [options]

  Update a policy set parameter. You can update the key, value, or
  sensitive flag. At least one field must be provided.

Options:

  -policy-set-id=<id>  Policy Set ID (required)
  -id=<parameter-id>   Parameter ID (required)
  -key=<name>          New parameter key
  -value=<value>       New parameter value
  -sensitive           Mark parameter as sensitive
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policysetparameter update -policy-set-id=polset-abc123 -id=var-xyz789 -value=2000
  hcptf policysetparameter update -policy-set-id=polset-abc123 -id=var-xyz789 -key=new_key -value=new_value
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set parameter update command
func (c *PolicySetParameterUpdateCommand) Synopsis() string {
	return "Update a policy set parameter"
}
