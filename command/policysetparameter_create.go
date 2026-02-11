package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicySetParameterCreateCommand is a command to create a policy set parameter
type PolicySetParameterCreateCommand struct {
	Meta
	policySetID string
	key         string
	value       string
	sensitive   bool
	format      string
}

// Run executes the policy set parameter create command
func (c *PolicySetParameterCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetparameter create")
	flags.StringVar(&c.policySetID, "policy-set-id", "", "Policy Set ID (required)")
	flags.StringVar(&c.key, "key", "", "Parameter key (required)")
	flags.StringVar(&c.value, "value", "", "Parameter value (required)")
	flags.BoolVar(&c.sensitive, "sensitive", false, "Mark parameter as sensitive")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policySetID == "" {
		c.Ui.Error("Error: -policy-set-id flag is required")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create policy set parameter
	category := tfe.CategoryPolicySet
	options := tfe.PolicySetParameterCreateOptions{
		Key:       tfe.String(c.key),
		Value:     tfe.String(c.value),
		Category:  &category,
		Sensitive: tfe.Bool(c.sensitive),
	}

	parameter, err := client.PolicySetParameters.Create(client.Context(), c.policySetID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating policy set parameter: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Policy set parameter '%s' created successfully", parameter.Key))

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

// Help returns help text for the policy set parameter create command
func (c *PolicySetParameterCreateCommand) Help() string {
	helpText := `
Usage: hcptf policysetparameter create [options]

  Create a parameter for a policy set. Parameters are key/value pairs that
  Sentinel uses during policy checks. Use the -sensitive flag for secret values.

Options:

  -policy-set-id=<id>  Policy Set ID (required)
  -key=<name>          Parameter key (required)
  -value=<value>       Parameter value (required)
  -sensitive           Mark parameter as sensitive (write-once, not visible thereafter)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policysetparameter create -policy-set-id=polset-abc123 -key=max_cost -value=1000
  hcptf policysetparameter create -policy-set-id=polset-abc123 -key=api_key -value=secret -sensitive
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set parameter create command
func (c *PolicySetParameterCreateCommand) Synopsis() string {
	return "Create a policy set parameter"
}
