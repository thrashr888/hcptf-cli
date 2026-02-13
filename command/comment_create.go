package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// CommentCreateCommand is a command to create a comment on a run
type CommentCreateCommand struct {
	Meta
	runID  string
	body   string
	format string
}

// Run executes the comment create command
func (c *CommentCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("comment create")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (required)")
	flags.StringVar(&c.body, "body", "", "Comment body (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.runID == "" {
		c.Ui.Error("Error: -run-id flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.body == "" {
		c.Ui.Error("Error: -body flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build create options
	options := tfe.CommentCreateOptions{
		Body: c.body,
	}

	// Create comment
	comment, err := client.Comments.Create(client.Context(), c.runID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating comment: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output("Comment created successfully")

	// Show comment details
	data := map[string]interface{}{
		"ID":   comment.ID,
		"Body": comment.Body,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the comment create command
func (c *CommentCreateCommand) Help() string {
	helpText := `
Usage: hcptf comment create [options]

  Create a comment on a run. Comments allow team members to leave
  feedback or record decisions about a run on a workspace.

Options:

  -run-id=<id>      Run ID (required)
  -body=<text>      Comment body (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf comment create -run-id=run-ABC123 -body="Approved for production"
  hcptf comment create -run-id=run-ABC123 -body="Need to review security implications" -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the comment create command
func (c *CommentCreateCommand) Synopsis() string {
	return "Create a comment on a run"
}
