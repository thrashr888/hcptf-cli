package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// CommentReadCommand is a command to show comment details
type CommentReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the comment read command
func (c *CommentReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("comment read")
	flags.StringVar(&c.id, "id", "", "Comment ID (required)")
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

	// Read comment
	comment, err := client.Comments.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading comment: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	// Show comment details
	data := map[string]interface{}{
		"ID":   comment.ID,
		"Body": comment.Body,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the comment read command
func (c *CommentReadCommand) Help() string {
	helpText := `
Usage: hcptf comment read [options]

  Show details of a comment.

Options:

  -id=<id>          Comment ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf comment read -id=wsc-ABC123
  hcptf comment read -id=wsc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the comment read command
func (c *CommentReadCommand) Synopsis() string {
	return "Show comment details"
}
