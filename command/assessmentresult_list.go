package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AssessmentResultListCommand lists assessment results for a workspace
type AssessmentResultListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

// Run executes the assessmentresult list command
func (c *AssessmentResultListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("assessmentresult list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
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

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Get workspace to verify it exists
	_, err = client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	// Since the go-tfe library doesn't have direct support for listing assessment results yet,
	// we provide guidance to the user
	c.Ui.Output("Note: Assessment result listing is not yet fully supported in the CLI.")
	c.Ui.Output("\nHealth assessments provide drift detection and continuous validation results.")
	c.Ui.Output("To view assessment results:")
	c.Ui.Output("")
	c.Ui.Output("1. View workspace health in the HCP Terraform UI")
	c.Ui.Output("2. Check drift detection results in workspace settings")
	c.Ui.Output("3. Use the API directly to list assessment results")
	c.Ui.Output("")
	c.Ui.Output("If you have an assessment result ID, you can view its details with:")
	c.Ui.Output("  hcptf assessmentresult read -id=<assessment-result-id>")
	c.Ui.Output("")
	c.Ui.Output("Assessment result IDs are typically prefixed with 'asmtres-'")
	c.Ui.Output("")
	c.Ui.Output("This feature requires HCP Terraform Plus or Enterprise with health")
	c.Ui.Output("assessments enabled in workspace settings.")

	// Return empty table to satisfy the formatter
	headers := []string{"Info"}
	rows := [][]string{{"See guidance above"}}
	formatter.Table(headers, rows)

	return 0
}

// Help returns help text for the assessmentresult list command
func (c *AssessmentResultListCommand) Help() string {
	helpText := `
Usage: hcptf assessmentresult list [options]

  List health assessment results for a workspace.

  Health assessments check if a workspace's real infrastructure matches
  its Terraform configuration. This includes drift detection and continuous
  validation results.

  Note: This feature requires HCP Terraform Plus or Enterprise, and
  health assessments must be enabled in workspace settings.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf assessmentresult list -org=my-org -workspace=my-workspace
  hcptf assessmentresult list -org=my-org -workspace=prod -output=json

Notes:

  Assessment results are generated when:
  - Drift detection runs automatically or manually
  - Continuous validation checks are performed
  - Health assessments are completed after applies

  To enable health assessments, update workspace settings via the UI
  or use: hcptf workspace update
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the assessmentresult list command
func (c *AssessmentResultListCommand) Synopsis() string {
	return "List health assessment results for a workspace"
}
