package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
	"os"
)

// PolicyCheckReadCommand is a command to read policy check details
type PolicyCheckReadCommand struct {
	Meta
	policyCheckID string
	format        string
}

// Run executes the policy check read command
func (c *PolicyCheckReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("policycheck read")
	flags.StringVar(&c.policyCheckID, "id", "", "Policy Check ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Validate required flags
	if c.policyCheckID == "" {
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

	// Read policy check
	policyCheck, err := client.PolicyChecks.Read(client.Context(), c.policyCheckID)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading policy check: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)
	if c.Meta.OutputWriter == nil && c.Meta.ErrorWriter == nil {
		formatter = output.NewFormatterWithWriters(c.format, os.Stdout, os.Stderr)
	}

	data := map[string]interface{}{
		"ID":            policyCheck.ID,
		"Status":        string(policyCheck.Status),
		"Scope":         string(policyCheck.Scope),
		"IsOverridable": policyCheck.Actions.IsOverridable,
		"CanOverride":   policyCheck.Permissions.CanOverride,
	}

	if policyCheck.Result != nil {
		data["Passed"] = policyCheck.Result.Passed
		data["TotalFailed"] = policyCheck.Result.TotalFailed
		data["HardFailed"] = policyCheck.Result.HardFailed
		data["SoftFailed"] = policyCheck.Result.SoftFailed
		data["AdvisoryFailed"] = policyCheck.Result.AdvisoryFailed
		data["DurationMs"] = policyCheck.Result.Duration
	}

	if policyCheck.StatusTimestamps != nil {
		if !policyCheck.StatusTimestamps.QueuedAt.IsZero() {
			data["QueuedAt"] = policyCheck.StatusTimestamps.QueuedAt
		}
		if !policyCheck.StatusTimestamps.PassedAt.IsZero() {
			data["PassedAt"] = policyCheck.StatusTimestamps.PassedAt
		}
		if !policyCheck.StatusTimestamps.HardFailedAt.IsZero() {
			data["HardFailedAt"] = policyCheck.StatusTimestamps.HardFailedAt
		}
		if !policyCheck.StatusTimestamps.SoftFailedAt.IsZero() {
			data["SoftFailedAt"] = policyCheck.StatusTimestamps.SoftFailedAt
		}
		if !policyCheck.StatusTimestamps.ErroredAt.IsZero() {
			data["ErroredAt"] = policyCheck.StatusTimestamps.ErroredAt
		}
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the policy check read command
func (c *PolicyCheckReadCommand) Help() string {
	helpText := `
Usage: hcptf policy check read [options]

  Read policy check details and results.

Options:

  -id=<policy-check-id>  Policy Check ID (required)
  -output=<format>       Output format: table (default) or json

Example:

  hcptf policy check read -id=polchk-abc123
  hcptf policy check read -id=polchk-abc123 -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the policy check read command
func (c *PolicyCheckReadCommand) Synopsis() string {
	return "Read policy check details"
}
