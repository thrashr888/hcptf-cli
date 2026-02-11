package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicySetParameterListCommand is a command to list policy set parameters
type PolicySetParameterListCommand struct {
	Meta
	policySetID string
	format      string
}

// Run executes the policy set parameter list command
func (c *PolicySetParameterListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policysetparameter list")
	flags.StringVar(&c.policySetID, "policy-set-id", "", "Policy Set ID (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List policy set parameters
	parameters, err := client.PolicySetParameters.List(client.Context(), c.policySetID, &tfe.PolicySetParameterListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing policy set parameters: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(parameters.Items) == 0 {
		c.Ui.Output("No policy set parameters found")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Key", "Value", "Category", "Sensitive"}
	var rows [][]string

	for _, param := range parameters.Items {
		value := param.Value
		if param.Sensitive {
			value = "(sensitive)"
		}

		sensitive := "No"
		if param.Sensitive {
			sensitive = "Yes"
		}

		rows = append(rows, []string{
			param.ID,
			param.Key,
			value,
			string(param.Category),
			sensitive,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the policy set parameter list command
func (c *PolicySetParameterListCommand) Help() string {
	helpText := `
Usage: hcptf policysetparameter list [options]

  List parameters for a policy set. Parameters are key/value pairs that
  Sentinel uses during policy checks.

Options:

  -policy-set-id=<id>  Policy Set ID (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policysetparameter list -policy-set-id=polset-abc123
  hcptf policysetparameter list -policy-set-id=polset-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set parameter list command
func (c *PolicySetParameterListCommand) Synopsis() string {
	return "List parameters for a policy set"
}
