package command

import (
	"fmt"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// PolicyCreateCommand is a command to create a policy
type PolicyCreateCommand struct {
	Meta
	organization string
	name         string
	description  string
	enforce      string
	policyFile   string
	format       string
}

// Run executes the policy create command
func (c *PolicyCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policy create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Policy name (required)")
	flags.StringVar(&c.description, "description", "", "Policy description")
	flags.StringVar(&c.enforce, "enforce", "advisory", "Enforcement level: advisory, soft-mandatory, or hard-mandatory")
	flags.StringVar(&c.policyFile, "policy-file", "", "Path to policy file (required)")
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

	if c.policyFile == "" {
		c.Ui.Error("Error: -policy-file flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate enforcement level
	var enforcementLevel tfe.EnforcementLevel
	switch c.enforce {
	case "advisory":
		enforcementLevel = tfe.EnforcementAdvisory
	case "soft-mandatory":
		enforcementLevel = tfe.EnforcementSoft
	case "hard-mandatory":
		enforcementLevel = tfe.EnforcementHard
	default:
		c.Ui.Error("Error: -enforce must be 'advisory', 'soft-mandatory', or 'hard-mandatory'")
		return 1
	}

	// Read policy file
	policyContent, err := os.ReadFile(c.policyFile)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy file: %s", err))
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create policy
	options := tfe.PolicyCreateOptions{
		Name:             tfe.String(c.name),
		EnforcementLevel: &enforcementLevel,
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	policy, err := client.Policies.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating policy: %s", err))
		return 1
	}

	// Upload policy content
	err = client.Policies.Upload(client.Context(), policy.ID, policyContent)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error uploading policy content: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Policy '%s' created successfully", policy.Name))

	// Show policy details
	data := map[string]interface{}{
		"ID":               policy.ID,
		"Name":             policy.Name,
		"Organization":     c.organization,
		"Description":      policy.Description,
		"EnforcementLevel": string(policy.EnforcementLevel),
		"PolicySetCount":   policy.PolicySetCount,
		"UpdatedAt":        policy.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy create command
func (c *PolicyCreateCommand) Help() string {
	helpText := `
Usage: hcptf policy create [options]

  Create a new policy.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Policy name (required)
  -description=<text>  Policy description
  -enforce=<level>     Enforcement level: advisory (default), soft-mandatory, or hard-mandatory
  -policy-file=<path>  Path to policy file (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policy create -org=my-org -name=my-policy -policy-file=policy.sentinel
  hcptf policy create -org=my-org -name=prod-policy -enforce=hard-mandatory -policy-file=policy.sentinel
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy create command
func (c *PolicyCreateCommand) Synopsis() string {
	return "Create a new policy"
}
