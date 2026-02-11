package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// NotificationVerifyCommand is a command to verify a notification configuration
type NotificationVerifyCommand struct {
	Meta
	id     string
	format string
}

// Run executes the notification verify command
func (c *NotificationVerifyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("notification verify")
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

	// Verify notification configuration
	notification, err := client.NotificationConfigurations.Verify(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error verifying notification configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Notification configuration '%s' verified successfully", notification.Name))
	c.Ui.Output("A test notification has been sent to the configured destination.")

	// Show notification details including delivery responses
	data := map[string]interface{}{
		"ID":              notification.ID,
		"Name":            notification.Name,
		"DestinationType": notification.DestinationType,
		"Enabled":         notification.Enabled,
		"URL":             notification.URL,
		"UpdatedAt":       notification.UpdatedAt,
	}

	// Add the most recent delivery response if available
	if len(notification.DeliveryResponses) > 0 {
		lastResponse := notification.DeliveryResponses[len(notification.DeliveryResponses)-1]
		data["LastDeliveryAt"] = lastResponse.SentAt
		data["LastDeliveryCode"] = lastResponse.Code
		data["LastDeliveryBody"] = lastResponse.Body
		data["LastDeliverySuccessful"] = lastResponse.Successful

		if lastResponse.Successful == "true" {
			c.Ui.Output("Verification successful: Test notification delivered successfully")
		} else {
			c.Ui.Warn(fmt.Sprintf("Verification warning: Delivery returned code %s", lastResponse.Code))
		}
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the notification verify command
func (c *NotificationVerifyCommand) Help() string {
	helpText := `
Usage: hcptf notification verify [options]

  Verify a notification configuration by sending a test notification.
  This command sends a test notification to the configured destination
  to verify that the configuration is working correctly.

Options:

  -id=<id>          Notification configuration ID (required)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf notification verify -id=nc-ABC123
  hcptf notification verify -id=nc-ABC123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the notification verify command
func (c *NotificationVerifyCommand) Synopsis() string {
	return "Verify a notification configuration"
}
