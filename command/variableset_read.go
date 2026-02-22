package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetReadCommand is a command to read variable set details
type VariableSetReadCommand struct {
	Meta
	id        string
	format    string
	varSetSvc variableSetReader
}

// Run executes the variable set read command
func (c *VariableSetReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset read")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
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

	// Read variable set
	variableSet, err := c.varSetService(client).Read(client.Context(), c.id, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading variable set: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":            variableSet.ID,
		"Name":          variableSet.Name,
		"Description":   variableSet.Description,
		"Global":        variableSet.Global,
		"VariableCount": len(variableSet.Variables),
	}

	// Add workspace information if not global
	if !variableSet.Global && len(variableSet.Workspaces) > 0 {
		workspaceNames := make([]string, 0, len(variableSet.Workspaces))
		for _, ws := range variableSet.Workspaces {
			workspaceNames = append(workspaceNames, ws.Name)
		}
		data["Workspaces"] = strings.Join(workspaceNames, ", ")
	}

	// Add project information if available
	if len(variableSet.Projects) > 0 {
		projectNames := make([]string, 0, len(variableSet.Projects))
		for _, proj := range variableSet.Projects {
			projectNames = append(projectNames, proj.Name)
		}
		data["Projects"] = strings.Join(projectNames, ", ")
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the variable set read command
func (c *VariableSetReadCommand) Help() string {
	helpText := `
Usage: hcptf variableset read [options]

  Read variable set details.

Options:

  -id=<id>          Variable set ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf variableset read -id=varset-12345
  hcptf variableset read -id=varset-12345 -output=json
`
	return strings.TrimSpace(helpText)
}

func (c *VariableSetReadCommand) varSetService(client *client.Client) variableSetReader {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Synopsis returns a short synopsis for the variable set read command
func (c *VariableSetReadCommand) Synopsis() string {
	return "Read variable set details"
}
