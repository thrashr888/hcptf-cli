package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// StateOutputsCommand is a command to display state outputs
type StateOutputsCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

// Run executes the state outputs command
func (c *StateOutputsCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("state outputs")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace first
	ws, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Get current state version
	currentStateVersion, err := client.StateVersions.ReadCurrent(client.Context(), ws.ID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading current state version: %s", err))
		return 1
	}

	if currentStateVersion == nil {
		c.Ui.Output("No current state version found for workspace")
		return 0
	}

	// Read state outputs
	outputsList, err := client.StateVersionOutputs.ReadCurrent(client.Context(), ws.ID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading state outputs: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(outputsList.Items) == 0 {
		c.Ui.Output("No outputs found in current state")
		return 0
	}

	// Format based on output type
	if c.format == "json" {
		// For JSON output, create a structured map
		jsonOutputs := make(map[string]interface{})
		for _, out := range outputsList.Items {
			outputData := map[string]interface{}{
				"sensitive": out.Sensitive,
				"type":      out.Type,
			}

			if out.Sensitive {
				outputData["value"] = "<sensitive>"
			} else {
				outputData["value"] = out.Value
			}

			if out.DetailedType != nil {
				outputData["detailed_type"] = out.DetailedType
			}

			jsonOutputs[out.Name] = outputData
		}
		formatter.JSON(jsonOutputs)
	} else {
		// For table output, display as key-value pairs
		headers := []string{"Name", "Value", "Sensitive", "Type"}
		var rows [][]string

		for _, out := range outputsList.Items {
			value := ""
			if out.Sensitive {
				value = "<sensitive>"
			} else {
				value = fmt.Sprintf("%v", out.Value)
			}

			rows = append(rows, []string{
				out.Name,
				value,
				fmt.Sprintf("%t", out.Sensitive),
				out.Type,
			})
		}

		formatter.Table(headers, rows)
	}

	return 0
}

// Help returns help text for the state outputs command
func (c *StateOutputsCommand) Help() string {
	helpText := `
Usage: hcptf state outputs [options]

  Display the outputs from the current state version of a workspace.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf state outputs -org=my-org -workspace=prod
  hcptf state outputs -org=my-org -workspace=prod -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the state outputs command
func (c *StateOutputsCommand) Synopsis() string {
	return "Display outputs from the current state version"
}
