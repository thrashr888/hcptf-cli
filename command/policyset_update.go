package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// PolicySetUpdateCommand is a command to update a policy set
type PolicySetUpdateCommand struct {
	Meta
	id           string
	name         string
	description  string
	global       string
	overridable  string
	agentEnabled string
	toolVersion  string
	policiesPath string
	format       string
}

// Run executes the policy set update command
func (c *PolicySetUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policyset update")
	flags.StringVar(&c.id, "id", "", "Policy set ID (required)")
	flags.StringVar(&c.name, "name", "", "Policy set name")
	flags.StringVar(&c.description, "description", "", "Policy set description")
	flags.StringVar(&c.global, "global", "", "Apply to all workspaces (true/false)")
	flags.StringVar(&c.overridable, "overridable", "", "Allow failed policy overrides (true/false)")
	flags.StringVar(&c.agentEnabled, "agent-enabled", "", "Run policy evaluations in an agent (true/false)")
	flags.StringVar(&c.toolVersion, "policy-tool-version", "", "Policy tool version")
	flags.StringVar(&c.policiesPath, "policies-path", "", "Subdirectory path for policy files in VCS")
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
	options := tfe.PolicySetUpdateOptions{}

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
			return 1
		}
	}
	if c.overridable != "" {
		overridable, parseErr := parseBoolFlag(c.overridable, "overridable")
		if parseErr != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
			return 1
		}
		options.Overridable = overridable
	}
	if c.agentEnabled != "" {
		agentEnabled, parseErr := parseBoolFlag(c.agentEnabled, "agent-enabled")
		if parseErr != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
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

	// Update policy set
	policySet, err := client.PolicySets.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating policy set: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Policy set '%s' updated successfully", policySet.Name))

	// Show policy set details
	data := map[string]interface{}{
		"ID":             policySet.ID,
		"Name":           policySet.Name,
		"Description":    policySet.Description,
		"Global":         policySet.Global,
		"PolicyCount":    policySet.PolicyCount,
		"WorkspaceCount": policySet.WorkspaceCount,
		"UpdatedAt":      policySet.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy set update command
func (c *PolicySetUpdateCommand) Help() string {
	helpText := `
Usage: hcptf policyset update [options]

  Update policy set settings.

Options:

  -id=<id>              Policy set ID (required)
  -name=<name>          Policy set name
  -description=<text>   Policy set description
  -global=<bool>        Apply to all workspaces (true/false)
  -overridable=<bool>   Allow failed policy overrides (true/false)
  -agent-enabled=<bool> Run policy evaluations in an agent (true/false)
  -policy-tool-version=<ver> Policy tool version
  -policies-path=<path> Subdirectory path for policy files in VCS
  -output=<format>      Output format: table (default) or json

Example:

  hcptf policyset update -id=polset-12345 -name=new-name
  hcptf policyset update -id=polset-12345 -overridable=false -agent-enabled=true
  hcptf policyset update -id=polset-12345 -global=true
  hcptf policyset update -id=polset-12345 -description="Updated description"
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy set update command
func (c *PolicySetUpdateCommand) Synopsis() string {
	return "Update policy set settings"
}
