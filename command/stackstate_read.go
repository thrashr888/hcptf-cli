package command

import (
	"fmt"
	"strings"
)

// StackStateReadCommand is a command to read stack state details
type StackStateReadCommand struct {
	Meta
	stateID string
	format  string
}

// Run executes the stack state read command
func (c *StackStateReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackstate read")
	flags.StringVar(&c.stateID, "id", "", "Stack state ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.stateID == "" {
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

	// Read stack state
	state, err := client.StackStates.Read(client.Context(), c.stateID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading stack state: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                    state.ID,
		"Generation":            state.Generation,
		"Deployment":            state.Deployment,
		"Status":                state.Status,
		"IsCurrent":             state.IsCurrent,
		"ResourceInstanceCount": state.ResourceInstanceCount,
	}

	if state.Stack != nil {
		data["StackID"] = state.Stack.ID
	}

	if state.StackDeploymentRun != nil {
		data["DeploymentRunID"] = state.StackDeploymentRun.ID
	}

	// Show components if available
	if len(state.Components) > 0 {
		componentNames := []string{}
		for _, comp := range state.Components {
			componentNames = append(componentNames, comp.Name)
		}
		data["Components"] = strings.Join(componentNames, ", ")
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack state read command
func (c *StackStateReadCommand) Help() string {
	helpText := `
Usage: hcptf stackstate read [options]

  Read stack state details including components and resource counts.

Options:

  -id=<state-id>    Stack state ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackstate read -id=sts-abc123
  hcptf stackstate read -id=sts-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack state read command
func (c *StackStateReadCommand) Synopsis() string {
	return "Read stack state details"
}
