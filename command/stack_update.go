package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// StackUpdateCommand is a command to update a stack
type StackUpdateCommand struct {
	Meta
	stackID            string
	name               string
	description        string
	speculativeEnabled *bool
	format             string
}

// Run executes the stack update command
func (c *StackUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stack update")
	flags.StringVar(&c.stackID, "id", "", "Stack ID (required)")
	flags.StringVar(&c.name, "name", "", "New stack name")
	flags.StringVar(&c.description, "description", "", "New stack description")

	// Use a string flag for boolean to distinguish between set and unset
	var speculativeEnabledStr string
	flags.StringVar(&speculativeEnabledStr, "speculative-enabled", "", "Enable/disable speculative plans (true/false)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Parse speculative-enabled flag if provided
	if speculativeEnabledStr != "" {
		if speculativeEnabledStr == "true" {
			val := true
			c.speculativeEnabled = &val
		} else if speculativeEnabledStr == "false" {
			val := false
			c.speculativeEnabled = &val
		} else {
			c.Ui.Error("Error: -speculative-enabled must be 'true' or 'false'")
			c.Ui.Error(c.Help())
			return 1
		}
	}

	// Validate required flags
	if c.stackID == "" {
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
	options := tfe.StackUpdateOptions{}

	if c.name != "" {
		options.Name = &c.name
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	if c.speculativeEnabled != nil {
		options.SpeculativeEnabled = tfe.Bool(*c.speculativeEnabled)
	}

	// Update stack
	stack, err := client.Stacks.Update(client.Context(), c.stackID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating stack: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Stack '%s' updated successfully", stack.Name))

	// Show stack details
	data := map[string]interface{}{
		"ID":                  stack.ID,
		"Name":                stack.Name,
		"Description":         stack.Description,
		"SpeculativeEnabled":  stack.SpeculativeEnabled,
		"UpdatedAt":           stack.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack update command
func (c *StackUpdateCommand) Help() string {
	helpText := `
Usage: hcptf stack update [options]

  Update stack settings.

Options:

  -id=<stack-id>                Stack ID (required)
  -name=<name>                  New stack name
  -description=<text>           New stack description
  -speculative-enabled=<bool>   Enable/disable speculative plans (true/false)
  -output=<format>              Output format: table (default) or json

Example:

  hcptf stack update -id=st-abc123 -name="New Name"
  hcptf stack update -id=st-abc123 -description="Updated description"
  hcptf stack update -id=st-abc123 -speculative-enabled=true
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack update command
func (c *StackUpdateCommand) Synopsis() string {
	return "Update stack settings"
}
