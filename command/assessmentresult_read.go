package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/hcptf-cli/internal/client"
)

// AssessmentResultReadCommand shows details of a specific assessment result
type AssessmentResultReadCommand struct {
	Meta
	id          string
	format      string
	showDrift   bool
	summaryOnly bool
}

// AssessmentResult represents a health assessment result
type AssessmentResult struct {
	ID   string               `json:"id"`
	Type string               `json:"type"`
	Data AssessmentResultData `json:"data"`
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
		Self              string `json:"self"`
		JSONOutput        string `json:"json-output"`
		JSONSchema        string `json:"json-schema"`
		LogOutput         string `json:"log-output"`
		HealthJSONRedacted string `json:"health-json-redacted"`
	} `json:"links"`
}

// AssessmentResultResponse represents the API response
type AssessmentResultResponse struct {
	ID   string               `json:"id"`
	Type string               `json:"type"`
	Data AssessmentResultData `json:"data"`
}

// TerraformPlan represents the Terraform JSON plan structure
type TerraformPlan struct {
	ResourceDrift []struct {
		Address  string `json:"address"`
		Mode     string `json:"mode"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Provider string `json:"provider_name"`
		Change   struct {
			Actions []string               `json:"actions"`
			Before  map[string]interface{} `json:"before"`
			After   map[string]interface{} `json:"after"`
		} `json:"change"`
	} `json:"resource_drift"`
	Checks []struct {
		Address struct {
			Kind      string `json:"kind"`       // "resource", "output", "check"
			Mode      string `json:"mode"`       // "managed", "data"
			Name      string `json:"name"`       // resource name
			ToDisplay string `json:"to_display"` // full resource address like "tls_self_signed_cert.user"
			Type      string `json:"type"`       // resource type
		} `json:"address"`
		Status    string `json:"status"` // "pass", "fail", "error", "unknown"
		Instances []struct {
			Address struct {
				ToDisplay string `json:"to_display"`
			} `json:"address"`
			Status   string `json:"status"`
			Problems []struct {
				Message string `json:"message"`
			} `json:"problems"`
		} `json:"instances"`
	} `json:"checks"`
}

// Run executes the assessmentresult read command
func (c *AssessmentResultReadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("assessmentresult read")
	flags.StringVar(&c.id, "id", "", "Assessment result ID (required)")
	flags.StringVar(&c.format, "output", "table", "Output format: table or json")
	flags.BoolVar(&c.showDrift, "show-drift", true, "Show detailed drift information (default: true)")
	flags.BoolVar(&c.summaryOnly, "summary-only", false, "Show only summary without drift details")

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

	// Get token from client for authorization
	token := client.Token()
	req.Header.Set("Authorization", "Bearer "+token)
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
	formatter := c.Meta.NewFormatter(c.format)

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

	// Show drift and check details if requested
	// Use health-json-redacted if available (has checks), otherwise fall back to json-output (drift only)
	outputURL := ar.Links.HealthJSONRedacted
	if outputURL == "" {
		outputURL = ar.Links.JSONOutput
	}

	if c.showDrift && !c.summaryOnly && outputURL != "" {
		if ar.Attributes.Drifted {
			c.Ui.Output("\n" + strings.Repeat("=", 80))
			c.Ui.Output("DRIFT DETAILS")
			c.Ui.Output(strings.Repeat("=", 80))
		}

		if err := c.showDriftDetails(client, outputURL); err != nil {
			c.Ui.Warn(fmt.Sprintf("\nWarning: Could not fetch drift details: %s", err))
			c.Ui.Output("\nTo retrieve detailed assessment outputs manually, use curl with your token:")
			c.Ui.Output(fmt.Sprintf("  JSON Plan: curl -H 'Authorization: Bearer $TOKEN' '%s%s'",
				client.GetAddress(), outputURL))
		}
	} else if ar.Attributes.Succeeded && !c.summaryOnly && !ar.Attributes.Drifted {
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

// showDriftDetails fetches and displays detailed drift information
func (c *AssessmentResultReadCommand) showDriftDetails(client *client.Client, jsonOutputPath string) error {
	// Fetch the JSON plan output
	fullURL := fmt.Sprintf("%s%s", client.GetAddress(), jsonOutputPath)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+client.Token())
	req.Header.Set("Accept", "application/json")

	httpClient := newHTTPClient()
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned %d: %s", resp.StatusCode, string(body))
	}

	var plan TerraformPlan
	if err := json.NewDecoder(resp.Body).Decode(&plan); err != nil {
		return fmt.Errorf("parsing JSON plan: %w", err)
	}

	// Process drifted resources
	if len(plan.ResourceDrift) == 0 {
		c.Ui.Output("\nNo resource drift details available in the plan output.")
		return nil
	}

	// Display drift summary
	driftCount := 0
	for _, res := range plan.ResourceDrift {
		if len(res.Change.Actions) > 0 && res.Change.Actions[0] != "no-op" {
			driftCount++
		}
	}

	if driftCount == 0 {
		c.Ui.Output("\nNo actual drift detected in resources (all resources are in sync).")
		return nil
	}

	c.Ui.Output(fmt.Sprintf("\n%d resource(s) have drifted:\n", driftCount))

	// Display each drifted resource
	index := 1
	for _, res := range plan.ResourceDrift {
		if len(res.Change.Actions) == 0 || res.Change.Actions[0] == "no-op" {
			continue
		}

		c.Ui.Output(fmt.Sprintf("%d. %s (%s.%s)", index, res.Address, res.Type, res.Name))
		c.Ui.Output(fmt.Sprintf("   Provider: %s", res.Provider))
		c.Ui.Output(fmt.Sprintf("   Action: %s", strings.Join(res.Change.Actions, ", ")))

		// Find and display changed attributes
		changedAttrs := c.findChangedAttributes(res.Change.Before, res.Change.After)
		if len(changedAttrs) > 0 && len(changedAttrs) <= 10 {
			c.Ui.Output("   Changed attributes:")
			for _, attr := range changedAttrs {
				prevVal := c.formatValue(res.Change.Before[attr])
				newVal := c.formatValue(res.Change.After[attr])
				c.Ui.Output(fmt.Sprintf("     • %s:", attr))
				c.Ui.Output(fmt.Sprintf("       - Previous: %s", prevVal))
				c.Ui.Output(fmt.Sprintf("       + Current:  %s", newVal))
			}
		} else if len(changedAttrs) > 10 {
			c.Ui.Output(fmt.Sprintf("   Changed attributes: %d attributes changed (too many to display)", len(changedAttrs)))
		}
		c.Ui.Output("")
		index++
	}

	// Show Terraform check results if available
	if len(plan.Checks) > 0 {
		c.Ui.Output("\n" + strings.Repeat("=", 80))
		c.Ui.Output("TERRAFORM CHECK RESULTS")
		c.Ui.Output(strings.Repeat("=", 80))

		passed := 0
		failed := 0
		errored := 0
		unknown := 0

		for _, check := range plan.Checks {
			switch check.Status {
			case "pass":
				passed++
			case "fail":
				failed++
			case "error":
				errored++
			default:
				unknown++
			}
		}

		c.Ui.Output(fmt.Sprintf("\nCheck Summary: %d passed, %d failed, %d errored, %d unknown\n", passed, failed, errored, unknown))

		// Show failed/errored checks
		if failed > 0 || errored > 0 {
			c.Ui.Output("Issues found:")
			for _, check := range plan.Checks {
				if check.Status == "fail" || check.Status == "error" {
					c.Ui.Output(fmt.Sprintf("\n  %s (%s) - %s", check.Address.ToDisplay, check.Address.Kind, strings.ToUpper(check.Status)))
					if check.Address.Name != "" {
						c.Ui.Output(fmt.Sprintf("  Name: %s", check.Address.Name))
					}
					// Show problems from instances
					for _, instance := range check.Instances {
						if instance.Status == "fail" || instance.Status == "error" {
							for _, problem := range instance.Problems {
								c.Ui.Output(fmt.Sprintf("    • %s", problem.Message))
							}
						}
					}
				}
			}
		} else if passed > 0 {
			c.Ui.Output("All checks passed successfully!")
		}
	}

	return nil
}

// findChangedAttributes compares before and after to find changed attributes
func (c *AssessmentResultReadCommand) findChangedAttributes(before, after map[string]interface{}) []string {
	changed := make([]string, 0)
	seen := make(map[string]bool)

	// Check all attributes in 'after'
	for key, afterVal := range after {
		beforeVal, exists := before[key]
		if !exists || !c.valuesEqual(beforeVal, afterVal) {
			changed = append(changed, key)
			seen[key] = true
		}
	}

	// Check for removed attributes (in before but not in after)
	for key := range before {
		if _, exists := after[key]; !exists && !seen[key] {
			changed = append(changed, key)
		}
	}

	return changed
}

// valuesEqual compares two values for equality
func (c *AssessmentResultReadCommand) valuesEqual(a, b interface{}) bool {
	// Simple comparison using JSON marshaling
	aJSON, _ := json.Marshal(a)
	bJSON, _ := json.Marshal(b)
	return string(aJSON) == string(bJSON)
}

// formatValue formats a value for display
func (c *AssessmentResultReadCommand) formatValue(val interface{}) string {
	if val == nil {
		return "<nil>"
	}

	switch v := val.(type) {
	case string:
		if len(v) > 100 {
			return v[:97] + "..."
		}
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		return fmt.Sprintf("%v", v)
	case map[string]interface{}:
		if len(v) == 0 {
			return "{}"
		}
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	case []interface{}:
		if len(v) == 0 {
			return "[]"
		}
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	default:
		return fmt.Sprintf("%v", v)
	}
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
  -show-drift       Show detailed drift information (default: true)
  -summary-only     Show only summary without drift details

Examples:

  # Show assessment with drift details
  hcptf assessmentresult read -id=asmtres-abc123

  # Show summary only
  hcptf assessmentresult read -id=asmtres-abc123 -summary-only

  # JSON output
  hcptf assessmentresult read -id=asmtres-abc123 -output=json

Assessment Result Details:

  - DriftStatus: Indicates if infrastructure has drifted from configuration
  - AssessmentStatus: Shows if the assessment completed successfully
  - JSONOutput: Link to the JSON plan output from the assessment
  - JSONSchema: Link to the provider schema used in the assessment
  - LogOutput: Link to Terraform JSON log output

Drift Details:

  When drift is detected and -show-drift is enabled (default), the command
  will automatically fetch and display which resources have drifted and what
  attributes have changed, making it easy to understand and remediate drift.
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the assessmentresult read command
func (c *AssessmentResultReadCommand) Synopsis() string {
	return "Show details of a specific health assessment result"
}
