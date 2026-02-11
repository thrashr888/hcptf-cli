package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/output"
)

// AssessmentResultReadCommand shows details of a specific assessment result
type AssessmentResultReadCommand struct {
	Meta
	id     string
	format string
}

// AssessmentResult represents a health assessment result
type AssessmentResult struct {
	ID   string                      `json:"id"`
	Type string                      `json:"type"`
	Data AssessmentResultData        `json:"data"`
}

// AssessmentResultData contains assessment result details
type AssessmentResultData struct {
	Attributes struct {
		Drifted   bool    `json:"drifted"`
		Succeeded bool    `json:"succeeded"`
		ErrorMsg  *string `json:"error-msg"`
		CreatedAt string  `json:"created-at"`
	} `json:"attributes"`
	Links struct {
		Self       string `json:"self"`
		JSONOutput string `json:"json-output"`
		JSONSchema string `json:"json-schema"`
		LogOutput  string `json:"log-output"`
	} `json:"links"`
}

// AssessmentResultResponse represents the API response
type AssessmentResultResponse struct {
	ID   string               `json:"id"`
	Type string               `json:"type"`
	Data AssessmentResultData `json:"data"`
}

// Run executes the assessmentresult read command
func (c *AssessmentResultReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("assessmentresult read")
	flags.StringVar(&c.id, "id", "", "Assessment result ID (required)")
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

	// Make direct API call to read assessment result
	apiURL := fmt.Sprintf("%s/api/v2/assessment-results/%s", client.GetAddress(), c.id)

	req, err := http.NewRequestWithContext(client.Context(), "GET", apiURL, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating request: %s", err))
		return 1
	}

	// Get token from config for authorization
	u := client.BaseURL()
	cfg, err := c.Meta.Config()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error loading config: %s", err))
		return 1
	}
	token := cfg.GetToken(u.Hostname())
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	httpClient := &http.Client{}
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
			c.Ui.Error("\nNote: Assessment result not found or Health Assessments feature is not available.")
			c.Ui.Error("This feature requires HCP Terraform Plus or Enterprise.")
		} else if resp.StatusCode == http.StatusForbidden {
			c.Ui.Error("\nNote: You may not have permission to view this assessment result.")
			c.Ui.Error("You need at least read access to the workspace.")
		}
		return 1
	}

	var assessmentResult AssessmentResultResponse
	if err := json.Unmarshal(body, &assessmentResult); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := output.NewFormatter(c.format)

	ar := assessmentResult.Data
	driftStatus := "No drift detected"
	if ar.Attributes.Drifted {
		driftStatus = "Drift detected"
	}

	assessmentStatus := "Succeeded"
	if !ar.Attributes.Succeeded {
		assessmentStatus = "Failed"
	}

	data := map[string]interface{}{
		"ID":               assessmentResult.ID,
		"DriftStatus":      driftStatus,
		"AssessmentStatus": assessmentStatus,
		"CreatedAt":        ar.Attributes.CreatedAt,
	}

	if ar.Attributes.ErrorMsg != nil && *ar.Attributes.ErrorMsg != "" {
		data["ErrorMessage"] = *ar.Attributes.ErrorMsg
	}

	// Add links to detailed outputs
	if ar.Links.JSONOutput != "" {
		data["JSONOutputURL"] = ar.Links.JSONOutput
	}
	if ar.Links.JSONSchema != "" {
		data["JSONSchemaURL"] = ar.Links.JSONSchema
	}
	if ar.Links.LogOutput != "" {
		data["LogOutputURL"] = ar.Links.LogOutput
	}

	formatter.KeyValue(data)

	// Provide helpful information about accessing detailed outputs
	if ar.Attributes.Succeeded {
		c.Ui.Output("\nTo retrieve detailed assessment outputs, use curl with your token:")
		if ar.Links.JSONOutput != "" {
			c.Ui.Output(fmt.Sprintf("  JSON Plan: curl -H 'Authorization: Bearer $TOKEN' '%s%s'",
				client.GetAddress(), ar.Links.JSONOutput))
		}
		if ar.Links.LogOutput != "" {
			c.Ui.Output(fmt.Sprintf("  Log Output: curl -H 'Authorization: Bearer $TOKEN' '%s%s'",
				client.GetAddress(), ar.Links.LogOutput))
		}
	}

	return 0
}

// Help returns help text for the assessmentresult read command
func (c *AssessmentResultReadCommand) Help() string {
	helpText := `
Usage: hcptf assessmentresult read [options]

  Show details of a specific health assessment result.

  Health assessments provide information about drift detection and
  continuous validation results, indicating whether infrastructure
  matches the Terraform configuration.

  Note: This feature requires HCP Terraform Plus or Enterprise, and
  you must have at least read access to the workspace.

Options:

  -id=<id>          Assessment result ID (required, format: asmtres-*)
  -output=<format>  Output format: table (default) or json

Example:

  hcptf assessmentresult read -id=asmtres-abc123
  hcptf assessmentresult read -id=asmtres-abc123 -output=json

Assessment Result Details:

  - DriftStatus: Indicates if infrastructure has drifted from configuration
  - AssessmentStatus: Shows if the assessment completed successfully
  - JSONOutput: Link to the JSON plan output from the assessment
  - JSONSchema: Link to the provider schema used in the assessment
  - LogOutput: Link to Terraform JSON log output

Retrieving Detailed Outputs:

  The assessment result includes links to detailed outputs (JSON plan,
  schema, and logs). These require admin-level workspace access and
  must be retrieved using direct API calls with a user or team token.

  Use curl or similar tools to fetch these outputs:
    curl -H "Authorization: Bearer $TOKEN" <output-url>
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the assessmentresult read command
func (c *AssessmentResultReadCommand) Synopsis() string {
	return "Show details of a specific health assessment result"
}
