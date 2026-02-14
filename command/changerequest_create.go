package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ChangeRequestCreateCommand creates a new change request
type ChangeRequestCreateCommand struct {
	Meta
	organization string
	workspace    string
	subject      string
	message      string
	format       string
}

// ChangeRequestCreatePayload represents the creation request
type ChangeRequestCreatePayload struct {
	Data struct {
		Type       string `json:"type"`
		Attributes struct {
			ActionType   string `json:"action_type"`
			ActionInputs struct {
				Subject string `json:"subject"`
				Message string `json:"message"`
			} `json:"action_inputs"`
			TargetIDs []string `json:"target_ids"`
		} `json:"attributes"`
	} `json:"data"`
}

// BulkActionResponse represents the bulk action response
type BulkActionResponse struct {
	Data struct {
		ID         string `json:"id"`
		Type       string `json:"type"`
		Attributes struct {
			OrganizationID string `json:"organization_id"`
			ActionType     string `json:"action_type"`
			ActionInputs   struct {
				Subject string `json:"subject"`
				Message string `json:"message"`
			} `json:"action_inputs"`
		} `json:"attributes"`
	} `json:"data"`
}

// Run executes the changerequest create command
func (c *ChangeRequestCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("changerequest create")
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name (required)")
	flags.StringVar(&c.subject, "subject", "", "Change request subject (required)")
	flags.StringVar(&c.message, "message", "", "Change request message (required)")
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

	if c.subject == "" {
		c.Ui.Error("Error: -subject flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	if c.message == "" {
		c.Ui.Error("Error: -message flag is required")
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

	// Build the request payload using the explorer bulk actions endpoint
	var payload ChangeRequestCreatePayload
	payload.Data.Type = "bulk_actions"
	payload.Data.Attributes.ActionType = "change_requests"
	payload.Data.Attributes.ActionInputs.Subject = c.subject
	payload.Data.Attributes.ActionInputs.Message = c.message
	payload.Data.Attributes.TargetIDs = []string{workspace.ID}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error marshaling payload: %s", err))
		return 1
	}

	// Make direct API call to create change request via bulk actions
	apiURL := fmt.Sprintf("%s/api/v2/organizations/%s/explorer/bulk-actions", client.GetAddress(), c.organization)

	req, err := http.NewRequestWithContext(client.Context(), "POST", apiURL, bytes.NewBuffer(payloadBytes))
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.Ui.Error(fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body)))
		if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusForbidden {
			c.Ui.Error("\nNote: Change Requests may not be available in your HCP Terraform plan.")
			c.Ui.Error("This feature requires HCP Terraform Plus or Enterprise.")
		}
		return 1
	}

	var response BulkActionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing response: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	c.Ui.Output(fmt.Sprintf("Change request created successfully via bulk action '%s'", response.Data.ID))

	data := map[string]interface{}{
		"BulkActionID":   response.Data.ID,
		"OrganizationID": response.Data.Attributes.OrganizationID,
		"Subject":        response.Data.Attributes.ActionInputs.Subject,
		"Message":        response.Data.Attributes.ActionInputs.Message,
		"WorkspaceID":    workspace.ID,
		"WorkspaceName":  workspace.Name,
	}

	formatter.KeyValue(data)
	return 0
}

// Help returns help text for the changerequest create command
func (c *ChangeRequestCreateCommand) Help() string {
	helpText := `
Usage: hcptf changerequest create [options]

  Create a new change request for a workspace.

  Change requests are used to track workspace to-dos, helping teams
  manage compliance and best practices. Creating a change request
  can trigger team notifications.

  Note: This feature requires HCP Terraform Plus or Enterprise.

Options:

  -organization=<name>  Organization name (required)
  -org=<name>          Alias for -organization
  -workspace=<name>    Workspace name (required)
  -subject=<text>      Change request subject line (required)
  -message=<text>      Change request message, supports markdown (required)
  -output=<format>     Output format: table (default) or json

Example:

  hcptf changerequest create -org=my-org -workspace=prod \
    -subject="[Action Required] Update Terraform version" \
    -message="Please update workspace to use Terraform 1.6.0"

  hcptf changerequest create -org=my-org -workspace=staging \
    -subject="Security: Enable GitHub Actions pinning" \
    -message="Pin all GitHub Actions to specific commit SHAs" \
    -output=json
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the changerequest create command
func (c *ChangeRequestCreateCommand) Synopsis() string {
	return "Create a new change request for a workspace"
}
