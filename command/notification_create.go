package command

import (
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/output"
)

// NotificationCreateCommand is a command to create a notification configuration
type NotificationCreateCommand struct {
	Meta
	organization    string
	workspace       string
	name            string
	destinationType string
	enabled         bool
	url             string
	token           string
	triggers        string
	emailAddresses  string
	format          string
}

// Run executes the notification create command
func (c *NotificationCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("notification create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.name, "name", "", "Notification configuration name (required)")
	flags.StringVar(&c.destinationType, "destination-type", "", "Destination type: email, slack, generic, microsoft-teams (required)")
	flags.BoolVar(&c.enabled, "enabled", true, "Enable notification configuration")
	flags.StringVar(&c.url, "url", "", "Webhook URL (required for slack, generic, microsoft-teams)")
	flags.StringVar(&c.token, "token", "", "Token for authentication (optional for generic)")
	flags.StringVar(&c.triggers, "triggers", "", "Comma-separated list of trigger types (e.g., run:created,run:completed)")
	flags.StringVar(&c.emailAddresses, "email-addresses", "", "Comma-separated list of email addresses (for email type, TFE only)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.organization == "" {
		c.Ui.Error("Error: -organization flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.workspace == "" {
		c.Ui.Error("Error: -workspace flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.destinationType == "" {
		c.Ui.Error("Error: -destination-type flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Validate destination type
	var destType tfe.NotificationDestinationType
	switch c.destinationType {
	case "email":
		destType = tfe.NotificationDestinationTypeEmail
	case "slack":
		destType = tfe.NotificationDestinationTypeSlack
	case "generic":
		destType = tfe.NotificationDestinationTypeGeneric
	case "microsoft-teams":
		destType = tfe.NotificationDestinationTypeMicrosoftTeams
	default:
		c.Ui.Error(fmt.Sprintf("Error: invalid destination-type '%s'. Must be one of: email, slack, generic, microsoft-teams", c.destinationType))
		return 1
	}

	// Validate URL for non-email types
	if destType != tfe.NotificationDestinationTypeEmail && c.url == "" {
		c.Ui.Error("Error: -url flag is required for slack, generic, and microsoft-teams destination types")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace to obtain its ID
	workspace, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Build create options
	options := tfe.NotificationConfigurationCreateOptions{
		Name:            tfe.String(c.name),
		DestinationType: &destType,
		Enabled:         tfe.Bool(c.enabled),
		SubscribableChoice: &tfe.NotificationConfigurationSubscribableChoice{
			Workspace: workspace,
		},
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
	if c.emailAddresses != "" && destType == tfe.NotificationDestinationTypeEmail {
		emailList := strings.Split(c.emailAddresses, ",")
		var emails []string
		for _, e := range emailList {
			emails = append(emails, strings.TrimSpace(e))
		}
		options.EmailAddresses = emails
	}

	// Create notification configuration
	notification, err := client.NotificationConfigurations.Create(client.Context(), workspace.ID, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating notification configuration: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Notification configuration '%s' created successfully", notification.Name))

	// Show notification details
	data := map[string]interface{}{
		"ID":              notification.ID,
		"Name":            notification.Name,
		"DestinationType": notification.DestinationType,
		"Enabled":         notification.Enabled,
		"URL":             notification.URL,
		"Triggers":        notification.Triggers,
		"CreatedAt":       notification.CreatedAt,
	}

	if destType == tfe.NotificationDestinationTypeEmail && len(notification.EmailAddresses) > 0 {
		data["EmailAddresses"] = notification.EmailAddresses
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the notification create command
func (c *NotificationCreateCommand) Help() string {
	helpText := `
Usage: hcptf notification create [options]

  Create a new notification configuration for a workspace.

Options:

  -organization=<name>       Organization name (required)
  -org=<name>               Alias for -organization
  -workspace=<name>         Workspace name (required)
  -name=<name>              Notification configuration name (required)
  -destination-type=<type>  Destination type (required)
                            Options: email, slack, generic, microsoft-teams
  -enabled                  Enable notification configuration (default: true)
  -url=<url>                Webhook URL (required for slack, generic, microsoft-teams)
  -token=<token>            Token for authentication (optional for generic)
  -triggers=<list>          Comma-separated list of trigger types
                            Options: run:created, run:planning, run:needs_attention,
                                    run:applying, run:completed, run:errored,
                                    assessment:drifted, assessment:failed,
                                    assessment:check_failure, workspace:auto_destroy_reminder,
                                    workspace:auto_destroy_run_results, change_request:created
  -email-addresses=<list>   Comma-separated list of email addresses (for email type, TFE only)
  -output=<format>          Output format: table (default) or json

Example:

  hcptf notification create -org=my-org -workspace=prod \
    -name="Slack Notifications" -destination-type=slack \
    -url=https://hooks.slack.com/services/XXX \
    -triggers=run:completed,run:errored

  hcptf notification create -org=my-org -workspace=dev \
    -name="Email Notifications" -destination-type=email \
    -email-addresses=admin@example.com,team@example.com \
    -triggers=run:needs_attention
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the notification create command
func (c *NotificationCreateCommand) Synopsis() string {
	return "Create a new notification configuration"
}
