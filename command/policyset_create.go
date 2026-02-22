package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicySetCreateCommand is a command to create a policy set
type PolicySetCreateCommand struct {
	Meta
	organization string
	name         string
	description  string
	global       bool
	kind         string
	overridable  string
	agentEnabled string
	toolVersion  string
	policiesPath string
	format       string
}

// Run executes the policy set create command
func (c *PolicySetCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Policy set name (required)")
	flags.StringVar(&c.description, "description", "", "Policy set description")
	flags.BoolVar(&c.global, "global", false, "Apply to all workspaces")
	flags.StringVar(&c.kind, "kind", "sentinel", "Policy set kind: sentinel or opa")
	flags.StringVar(&c.overridable, "overridable", "", "Allow failed policy overrides (true/false)")
	flags.StringVar(&c.agentEnabled, "agent-enabled", "", "Run policy evaluations in an agent (true/false)")
	flags.StringVar(&c.toolVersion, "policy-tool-version", "", "Policy tool version")
	flags.StringVar(&c.policiesPath, "policies-path", "", "Subdirectory path for policy files in VCS")
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

	kind, parseErr := parsePolicyKind(c.kind)
	if parseErr != nil {
		c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
		return 1
	}

	// Create policy set
	options := tfe.PolicySetCreateOptions{
		Name:   tfe.String(c.name),
		Global: tfe.Bool(c.global),
		Kind:   kind,
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}
	if c.overridable != "" {
		overridable, err := parseBoolFlag(c.overridable, "overridable")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", err))
			return 1
		}
		options.Overridable = overridable
	}
	if c.agentEnabled != "" {
		agentEnabled, err := parseBoolFlag(c.agentEnabled, "agent-enabled")
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", err))
			return 1
		}
		options.AgentEnabled = agentEnabled
	}
	if c.toolVersion != "" {
		options.PolicyToolVersion = tfe.String(c.toolVersion)
	}
	if c.policiesPath != "" {
		options.PoliciesPath = tfe.String(c.policiesPath)
	}

	policySet, err := client.PolicySets.Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating policy set: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Policy set '%s' created successfully", policySet.Name))

	// Show policy set details
	data := map[string]interface{}{
		"ID":             policySet.ID,
		"Name":           policySet.Name,
		"Description":    policySet.Description,
		"Global":         policySet.Global,
		"PolicyCount":    policySet.PolicyCount,
		"WorkspaceCount": policySet.WorkspaceCount,
		"CreatedAt":      policySet.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy set create command
func (c *PolicySetCreateCommand) Help() string {
	helpText := `
Usage: hcptf policyset create [options]

  Create a new policy set in an organization.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Policy set name (required)
  -description=<text>  Policy set description
  -global              Apply to all workspaces (default: false)
  -kind=<kind>         Policy set kind: sentinel (default) or opa
  -overridable=<bool>  Allow failed policy overrides (true/false)
  -agent-enabled=<bool> Run policy evaluations in an agent (true/false)
  -policy-tool-version=<ver> Policy tool version
  -policies-path=<path> Subdirectory path for policies in VCS
  -output=<format>     Output format: table (default) or json

Example:

  hcptf policyset create -org=my-org -name=security-policies -description="Security policies"
  hcptf policyset create -org=my-org -name=opa-policies -kind=opa -overridable=false
  hcptf policyset create -org=my-org -name=global-policies -global
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set create command
func (c *PolicySetCreateCommand) Synopsis() string {
	return "Create a new policy set"
}
