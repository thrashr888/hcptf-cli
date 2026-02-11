package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// NotificationReadCommand is a command to read notification configuration details
type NotificationReadCommand struct {
	Meta
	id     string
	format string
}

// Run executes the notification read command
func (c *NotificationReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("notification read")
	flags.StringVar(&c.id, "id", "", "Notification configuration ID (required)")
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

	// Read notification configuration
	notification, err := client.NotificationConfigurations.Read(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading notification configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	data := map[string]interface{}{
		"ID":              notification.ID,
		"Name":            notification.Name,
		"DestinationType": notification.DestinationType,
		"Enabled":         notification.Enabled,
		"URL":             notification.URL,
		"Triggers":        notification.Triggers,
		"CreatedAt":       notification.CreatedAt,
		"UpdatedAt":       notification.UpdatedAt,
	}

	// Add email addresses if available (TFE only)
	if len(notification.EmailAddresses) > 0 {
		data["EmailAddresses"] = notification.EmailAddresses
	}

	// Add delivery responses if available
	if len(notification.DeliveryResponses) > 0 {
		var responses []string
		for _, dr := range notification.DeliveryResponses {
			responses = append(responses, fmt.Sprintf("%s: %s (Code: %s)", dr.SentAt, dr.Body, dr.Code))
		}
		data["DeliveryResponses"] = responses
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the notification read command
func (c *NotificationReadCommand) Help() string {
	helpText := `
Usage: hcptf notification read [options]

  Read notification configuration details.

Options:

  -id=<id>          Notification configuration ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf notification read -id=nc-ABC123
  hcptf notification read -id=nc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the notification read command
func (c *NotificationReadCommand) Synopsis() string {
	return "Read notification configuration details"
}
