package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// CommentListCommand is a command to list comments for a run
type CommentListCommand struct {
	Meta
	runID  string
	format string
}

// Run executes the comment list command
func (c *CommentListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("comment list")
	flags.StringVar(&c.runID, "run-id", "", "Run ID (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// List comments for the run
	comments, err := client.Comments.List(client.Context(), c.runID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error listing comments: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	if len(comments.Items) == 0 {
		c.Ui.Output("No comments found for this run")
		return 0
	}

	// Prepare table data
	headers := []string{"ID", "Body"}
	var rows [][]string

	for _, comment := range comments.Items {
		// Truncate body if too long for table display
		body := comment.Body
		if len(body) > 100 {
			body = body[:97] + "..."
		}
		// Replace newlines with spaces for table display
		body = strings.ReplaceAll(body, "\n", " ")

		rows = append(rows, []string{
			comment.ID,
			body,
		})
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the comment list command
func (c *CommentListCommand) Help() string {
	helpText := `
Usage: hcptf comment list [options]

  List comments for a run. Comments allow team members to leave
  feedback or record decisions about a run on a workspace.

Options:

  -run-id=<id>      Run ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf comment list -run-id=run-ABC123
  hcptf comment list -run-id=run-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the comment list command
func (c *CommentListCommand) Synopsis() string {
	return "List comments for a run"
}
