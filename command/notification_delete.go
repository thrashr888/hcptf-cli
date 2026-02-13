package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// NotificationDeleteCommand is a command to delete a notification configuration
type NotificationDeleteCommand struct {
	Meta
	id       string
	force    bool
	notifSvc notificationDeleter
}

// Run executes the notification delete command
func (c *NotificationDeleteCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("notification delete")
	flags.StringVar(&c.id, "id", "", "Notification configuration ID (required)")
	flags.BoolVar(&c.force, "force", false, "Force delete without confirmation")

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

	// Confirm deletion unless force flag is set
	if !c.force {
		// Read notification configuration to get its name for confirmation
		notification, err := client.NotificationConfigurations.Read(client.Context(), c.id)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading notification configuration: %s", err))
			return 1
		}

		confirmation, err := c.Ui.Ask(fmt.Sprintf("Are you sure you want to delete notification configuration '%s' (%s)? (yes/no): ", notification.Name, c.id))
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading confirmation: %s", err))
			return 1
		}

		if strings.ToLower(confirmation) != "yes" {
			c.Ui.Output("Deletion cancelled")
			return 0
		}
	}

	// Delete notification configuration
	err = c.notificationService(client).Delete(client.Context(), c.id)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error deleting notification configuration: %s", err))
		return 1
	}

	c.Ui.Output(fmt.Sprintf("Notification configuration '%s' deleted successfully", c.id))
	return 0
}

func (c *NotificationDeleteCommand) notificationService(client *client.Client) notificationDeleter {
	if c.notifSvc != nil {
		return c.notifSvc
	}
	return client.NotificationConfigurations
}

// Help returns help text for the notification delete command
func (c *NotificationDeleteCommand) Help() string {
	helpText := `
Usage: hcptf notification delete [options]

  Delete a notification configuration.

Options:

  -id=<id>  Notification configuration ID (required)
  -force    Force delete without confirmation

Example:

  hcptf notification delete -id=nc-ABC123
  hcptf notification delete -id=nc-ABC123 -force
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the notification delete command
func (c *NotificationDeleteCommand) Synopsis() string {
	return "Delete a notification configuration"
}
