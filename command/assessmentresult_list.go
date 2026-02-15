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
	Data     interface{}               `json:"data"`
	Included []AssessmentResultListItem `json:"included"`
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

	// Make API call to get workspace with assessment result included
	apiURL := fmt.Sprintf("%s/api/v2/workspaces/%s?assessment_meta=true&include=current_assessment_result", client.GetAddress(), workspace.ID)
	req, err := http.NewRequestWithContext(client.Context(), "GET", apiURL, nil)
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
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Filter included items for assessment-results
	var assessmentResults []AssessmentResultListItem
	for _, item := range response.Included {
		if item.Type == "assessment-results" {
			assessmentResults = append(assessmentResults, item)
		}
	}

	if len(assessmentResults) == 0 {
		c.Ui.Output("No assessment results found")
		c.Ui.Output("\nNote: Health assessments must be enabled in workspace settings.")
		c.Ui.Output("This feature requires HCP Terraform Plus or Enterprise.")
		return 0
	}

	// Since there's only ever the latest result, automatically show details
	// instead of just listing (which would be a pointless extra step)
	ar := assessmentResults[0]

	// Show details by running the read command
	readCmd := &AssessmentResultReadCommand{
		Meta:   c.Meta,
		format: c.format,
	}
	return readCmd.Run([]string{"-id", ar.ID})
}

// Help returns help text for the assessmentresult list command
func (c *AssessmentResultListCommand) Help() string {
	helpText := `
Usage: hcptf workspace run assessmentresult list [options]

  List health assessment results for a workspace.

  Health assessments check if a workspace's real infrastructure matches
  its Terraform configuration. This includes drift detection and continuous
  validation results.

  Note: This feature requires HCP Terraform Plus or Enterprise, and
  health assessments must be enabled in workspace settings.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -name=<name>         Workspace name (required)
  -workspace=<name>    Workspace name (alias)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf workspace run assessmentresult list -org=my-org -name=my-workspace
  hcptf workspace run assessmentresult list -org=my-org -name=prod -output=json

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
