package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// NotificationUpdateCommand is a command to update a notification configuration
type NotificationUpdateCommand struct {
	Meta
	id             string
	name           string
	enabled        string
	url            string
	token          string
	triggers       string
	emailAddresses string
	format         string
}

// Run executes the notification update command
func (c *NotificationUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("notification update")
	flags.StringVar(&c.id, "id", "", "Notification configuration ID (required)")
	flags.StringVar(&c.name, "name", "", "Notification configuration name")
	flags.StringVar(&c.enabled, "enabled", "", "Enable notification configuration (true/false)")
	flags.StringVar(&c.url, "url", "", "Webhook URL")
	flags.StringVar(&c.token, "token", "", "Token for authentication")
	flags.StringVar(&c.triggers, "triggers", "", "Comma-separated list of trigger types")
	flags.StringVar(&c.emailAddresses, "email-addresses", "", "Comma-separated list of email addresses (TFE only)")
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

	// Build update options
	options := tfe.NotificationConfigurationUpdateOptions{}

	if c.name != "" {
		options.Name = tfe.String(c.name)
	}

	if c.enabled != "" {
		if c.enabled == "true" {
			options.Enabled = tfe.Bool(true)
		} else if c.enabled == "false" {
			options.Enabled = tfe.Bool(false)
		} else {
			c.Ui.Error("Error: -enabled must be 'true' or 'false'")
			return 1
		}
	}

	if c.url != "" {
		options.URL = tfe.String(c.url)
	}

	if c.token != "" {
		options.Token = tfe.String(c.token)
	}

	// Parse triggers
	if c.triggers != "" {
		triggerList := strings.Split(c.triggers, ",")
		var triggers []tfe.NotificationTriggerType
		for _, t := range triggerList {
			triggers = append(triggers, tfe.NotificationTriggerType(strings.TrimSpace(t)))
		}
		options.Triggers = triggers
	}

	// Parse email addresses (TFE only)
	if c.emailAddresses != "" {
		emailList := strings.Split(c.emailAddresses, ",")
		var emails []string
		for _, e := range emailList {
			emails = append(emails, strings.TrimSpace(e))
		}
		options.EmailAddresses = emails
	}

	// Update notification configuration
	notification, err := client.NotificationConfigurations.Update(client.Context(), c.id, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating notification configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Notification configuration '%s' updated successfully", notification.Name))

	// Show notification details
	data := map[string]interface{}{
		"ID":              notification.ID,
		"Name":            notification.Name,
		"DestinationType": notification.DestinationType,
		"Enabled":         notification.Enabled,
		"URL":             notification.URL,
		"Triggers":        notification.Triggers,
		"UpdatedAt":       notification.UpdatedAt,
	}

	if len(notification.EmailAddresses) > 0 {
		data["EmailAddresses"] = notification.EmailAddresses
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the notification update command
func (c *NotificationUpdateCommand) Help() string {
	helpText := `
Usage: hcptf notification update [options]

  Update notification configuration settings.

Options:

  -id=<id>                  Notification configuration ID (required)
  -name=<name>              Notification configuration name
  -enabled=<bool>           Enable notification configuration (true/false)
  -url=<url>                Webhook URL
  -token=<token>            Token for authentication
  -triggers=<list>          Comma-separated list of trigger types
                            Options: run:created, run:planning, run:needs_attention,
                                    run:applying, run:completed, run:errored,
                                    assessment:drifted, assessment:failed,
                                    assessment:check_failure, workspace:auto_destroy_reminder,
                                    workspace:auto_destroy_run_results, change_request:created
  -email-addresses=<list>   Comma-separated list of email addresses (TFE only)
  -output=<format>          Output format: table (default) or json

Example:

  hcptf notification update -id=nc-ABC123 -enabled=false
  hcptf notification update -id=nc-ABC123 -name="Updated Slack Notifications" \
    -triggers=run:completed,run:errored,run:needs_attention
  hcptf notification update -id=nc-ABC123 -url=https://hooks.slack.com/services/NEW
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the notification update command
func (c *NotificationUpdateCommand) Synopsis() string {
	return "Update notification configuration settings"
}
