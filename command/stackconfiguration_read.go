package command

import (
	"fmt"
	"strings"
)

// StackConfigurationReadCommand is a command to read stack configuration details
type StackConfigurationReadCommand struct {
	Meta
	configID string
	format   string
}

// Run executes the stack configuration read command
func (c *StackConfigurationReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stackconfiguration read")
	flags.StringVar(&c.configID, "id", "", "Stack configuration ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.configID == "" {
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

	// Read stack configuration
	config, err := client.StackConfigurations.Read(client.Context(), c.configID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading stack configuration: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":             config.ID,
		"SequenceNumber": config.SequenceNumber,
		"Status":         config.Status,
		"Speculative":    config.Speculative,
		"CreatedAt":      config.CreatedAt,
		"UpdatedAt":      config.UpdatedAt,
	}

	if config.Stack != nil {
		data["StackID"] = config.Stack.ID
	}

	// Show components if available
	if len(config.Components) > 0 {
		componentNames := []string{}
		for _, comp := range config.Components {
			componentNames = append(componentNames, comp.Name)
		}
		data["Components"] = strings.Join(componentNames, ", ")
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack configuration read command
func (c *StackConfigurationReadCommand) Help() string {
	helpText := `
Usage: hcptf stackconfiguration read [options]

  Read stack configuration details including status and components.

Options:

  -id=<config-id>   Stack configuration ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stackconfiguration read -id=stc-abc123
  hcptf stackconfiguration read -id=stc-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack configuration read command
func (c *StackConfigurationReadCommand) Synopsis() string {
	return "Read stack configuration details"
}
