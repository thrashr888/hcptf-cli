package command

import (
	"fmt"
	"strings"

)

// StackReadCommand is a command to read stack details
type StackReadCommand struct {
	Meta
	stackID string
	format  string
}

// Run executes the stack read command
func (c *StackReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("stack read")
	flags.StringVar(&c.stackID, "id", "", "Stack ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
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

	// Read stack
	stack, err := client.Stacks.Read(client.Context(), c.stackID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading stack: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":                 stack.ID,
		"Name":               stack.Name,
		"Description":        stack.Description,
		"SpeculativeEnabled": stack.SpeculativeEnabled,
		"CreatedAt":          stack.CreatedAt,
		"UpdatedAt":          stack.UpdatedAt,
	}

	if stack.Project != nil {
		data["ProjectID"] = stack.Project.ID
		data["ProjectName"] = stack.Project.Name
	}

	if stack.VCSRepo != nil {
		if stack.VCSRepo.Identifier != "" {
			data["VCSIdentifier"] = stack.VCSRepo.Identifier
		}
		if stack.VCSRepo.Branch != "" {
			data["VCSBranch"] = stack.VCSRepo.Branch
		}
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the stack read command
func (c *StackReadCommand) Help() string {
	helpText := `
Usage: hcptf stack read [options]

  Read stack details including configuration and status.

Options:

  -id=<stack-id>    Stack ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf stack read -id=st-abc123
  hcptf stack read -id=st-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the stack read command
func (c *StackReadCommand) Synopsis() string {
	return "Read stack details"
}
