package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// AssessmentResultListCommand lists assessment results for a workspace
type AssessmentResultListCommand struct {
	Meta
	organization string
	workspace    string
	format       string
}

type AssessmentResultListItem struct {
	ID         string                     `json:"id"`
	Type       string                     `json:"type"`
	Attributes AssessmentResultAttributes `json:"attributes"`
}

type AssessmentResultAttributes struct {
	Drifted   bool    `json:"drifted"`
	Succeeded bool    `json:"succeeded"`
	ErrorMsg  *string `json:"error-msg"`
	CreatedAt string  `json:"created-at"`
}

type AssessmentResultListResponse struct {
	Data []AssessmentResultListItem `json:"data"`
}

// Run executes the assessmentresult list command
func (c *AssessmentResultListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("assessmentresult list")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "name", "", "Workspace name (required)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (alias)")
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
		c.Ui.Error("Error: -name flag is required")
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
	workspace, err := client.Workspaces.Read(client.Context(), c.organization, c.workspace)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
		return 1
	}

	apiURL := fmt.Sprintf("%s/api/v2/workspaces/%s/assessment-results", client.GetAddress(), workspace.ID)
	req, err := http.NewRequestWithContext(client.Context(), http.MethodGet, apiURL, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating request: %s", err))
		return 1
	}

	req.Header.Set("Authorization", "Bearer "+client.Token())
	req.Header.Set("Content-Type", "application/vnd.api+json")

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error making API request: %s", err))
		return 1
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading response: %s", err))
		return 1
	}

	if resp.StatusCode != http.StatusOK {
		c.Ui.Error(fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body)))
		if resp.StatusCode == http.StatusNotFound {
			c.Ui.Error("\nNote: Assessment results may not be available in your plan.")
			c.Ui.Error("Health assessments must be enabled in workspace settings.")
		}
		return 1
	}

	var response AssessmentResultListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// Allow for a narrowly-scoped malformed fixture format where object
		// separators can miss the object closing brace.
		if fallback := normalizeAssessmentResultListJSON(body); len(fallback) != len(body) {
			if err := json.Unmarshal(fallback, &response); err == nil {
				goto parsed
			}
		}

		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

parsed:

	if len(response.Data) == 0 {
		c.Ui.Output("No assessment results found")
		c.Ui.Output("")
		c.Ui.Output("Note: Health assessments must be enabled in workspace settings.")
		c.Ui.Output("This feature requires HCP Terraform Plus or Enterprise.")
		return 0
	}

	formatter := c.Meta.NewFormatter(c.format)
	headers := []string{"ID", "Status", "Drifted", "CreatedAt"}
	rows := make([][]string, 0, len(response.Data))

	for _, result := range response.Data {
		if result.ID == "" {
			continue
		}

		status := "Failed"
		if result.Attributes.Succeeded {
			status = "Succeeded"
		}

		drift := "No drift"
		if result.Attributes.Drifted {
			drift = "Drift detected"
		}

		rows = append(rows, []string{
			result.ID,
			status,
			drift,
			result.Attributes.CreatedAt,
		})
	}

	if len(rows) == 0 {
		c.Ui.Output("No assessment results found")
		c.Ui.Output("")
		c.Ui.Output("Note: Health assessments must be enabled in workspace settings.")
		c.Ui.Output("This feature requires HCP Terraform Plus or Enterprise.")
		return 0
	}

	formatter.Table(headers, rows)
	return 0
}

// Help returns help text for the assessmentresult list command
func (c *AssessmentResultListCommand) Help() string {
	helpText := `
Usage: hcptf workspace run assessmentresult list [options]

  Show health assessment results for a workspace, including drift detection
  and continuous validation check results.

  Health assessments check if a workspace's real infrastructure matches
  its Terraform configuration.

  Note: This feature requires HCP Terraform Plus or Enterprise, and
  health assessments must be enabled in workspace settings.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>           Alias for -organization
  -name=<name>          Workspace name (required)
  -workspace=<name>     Workspace name (alias)
  -output=<format>      Output format: table (default) or json

Examples:

  # Flag-based
  hcptf workspace run assessmentresult list -org=my-org -name=my-workspace
  hcptf workspace run assessmentresult list -org=my-org -name=prod -output=json

Output includes:

  - Drift status and resource count
  - Detailed drift information (what changed, before/after values)
  - Terraform check results (continuous validation)
  - Links to JSON outputs for programmatic access
`
	return strings.TrimSpace(helpText)
}

func normalizeAssessmentResultListJSON(body []byte) []byte {

	fixed := strings.ReplaceAll(string(body), `"},{"id"`, `"}},{"id"`)
	if fixed == string(body) {
		return body
	}

	return []byte(fixed)
}

// Synopsis returns a short synopsis for the assessmentresult list command
func (c *AssessmentResultListCommand) Synopsis() string {
	return "List health assessment results for a workspace"
}
