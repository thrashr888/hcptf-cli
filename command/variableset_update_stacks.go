package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// VariableSetUpdateStacksCommand synchronizes stack associations.
type VariableSetUpdateStacksCommand struct {
	Meta
	id        string
	stacks    string
	format    string
	varSetSvc variableSetStackUpdater
}

// Run executes the variableset update-stacks command.
func (c *VariableSetUpdateStacksCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("variableset update-stacks")
	flags.StringVar(&c.id, "id", "", "Variable set ID (required)")
	flags.StringVar(&c.stacks, "stacks", "", "Comma-separated stack IDs (empty clears all)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}
	if c.id == "" {
		c.Ui.Error("Error: -id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	stacks := make([]*tfe.Stack, 0)
	for _, stackID := range splitCommaList(c.stacks) {
		if strings.TrimSpace(stackID) == "" {
			continue
		}
		stacks = append(stacks, &tfe.Stack{ID: strings.TrimSpace(stackID)})
	}

	options := &tfe.VariableSetUpdateStacksOptions{
		Stacks: stacks,
	}
	updated, err := c.varSetService(client).UpdateStacks(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating variable set stacks: %s", err))
		return 1
	}

	formatter := c.Meta.NewFormatter(c.format)
	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Variable set '%s' stack associations updated", updated.Name))
	}
	formatter.KeyValue(map[string]interface{}{
		"ID":         updated.ID,
		"Name":       updated.Name,
		"StackCount": len(updated.Stacks),
	})
	return 0
}

func (c *VariableSetUpdateStacksCommand) varSetService(client *client.Client) variableSetStackUpdater {
	if c.varSetSvc != nil {
		return c.varSetSvc
	}
	return client.VariableSets
}

// Help returns help text for variableset update-stacks.
func (c *VariableSetUpdateStacksCommand) Help() string {
	helpText := `
Usage: hcptf variableset update-stacks [options]

  Replace stack associations for a variable set.

Options:

  -id=<id>          Variable set ID (required)
  -stacks=<ids>     Comma-separated stack IDs (empty clears all)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf variableset update-stacks -id=varset-12345 -stacks=stack-abc123,stack-def456
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for variableset update-stacks.
func (c *VariableSetUpdateStacksCommand) Synopsis() string {
	return "Replace variable set stack associations"
}
