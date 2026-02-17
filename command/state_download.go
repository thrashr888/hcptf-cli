package command

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// StateDownloadCommand is a command to download state file contents
type StateDownloadCommand struct {
	Meta
	organization         string
	workspace            string
	stateVersionID       string
	outputFile           string
	stateDownloadSvc     stateVersionReader
	workspaceDownloadSvc workspaceReader
}

// Run executes the state download command
func (c *StateDownloadCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("state download")
	flags.StringVar(&c.organization, "org", "", "Organization name")
	flags.StringVar(&c.workspace, "workspace", "", "Workspace name")
	flags.StringVar(&c.stateVersionID, "id", "", "State version ID (optional, defaults to current)")
	flags.StringVar(&c.outputFile, "output", "", "Output file path (required)")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Must provide either (org + workspace) or state version ID
	if c.stateVersionID == "" && (c.organization == "" || c.workspace == "") {
		c.Ui.Error("Error: must provide either -id or both -org and -workspace")
		c.Ui.Error(c.Help())
		return 1
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	var stateVersion *tfe.StateVersion

	// If state version ID not provided, get current state from workspace
	if c.stateVersionID == "" && c.organization != "" && c.workspace != "" {
		// Get workspace to find current state version
		ws, err := c.workspaceService(client).Read(client.Context(), c.organization, c.workspace)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading workspace: %s", err))
			return 1
		}

		if ws.CurrentStateVersion == nil {
			c.Ui.Error("Error: workspace has no state version")
			return 1
		}

		// Read full state version details using workspace ID
		stateVersion, err = c.stateService(client).ReadCurrent(client.Context(), ws.ID)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("Error reading state version: %s", err))
			return 1
		}
	} else if c.stateVersionID != "" {
		// For now, we need workspace context to get state version
		// The go-tfe library doesn't have a ReadByID method for state versions
		c.Ui.Error("Error: -id flag not yet supported, please use -org and -workspace")
		return 1
	}

	if stateVersion.DownloadURL == "" {
		c.Ui.Error("Error: state version has no download URL")
		return 1
	}

	// Download state file with authentication
	req, err := http.NewRequest("GET", stateVersion.DownloadURL, nil)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating download request: %s", err))
		return 1
	}

	// Add authentication token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", client.Token()))

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error downloading state file: %s", err))
		return 1
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Ui.Error(fmt.Sprintf("Error: received status code %d from download URL", resp.StatusCode))
		return 1
	}

	// Read the state content
	stateContent, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error reading state content: %s", err))
		return 1
	}

	// Validate it's valid JSON
	var jsonCheck interface{}
	if err := json.Unmarshal(stateContent, &jsonCheck); err != nil {
		c.Ui.Error(fmt.Sprintf("Error: downloaded content is not valid JSON: %s", err))
		return 1
	}

	// If output file specified, write to file
	if c.outputFile != "" {
		if err := os.WriteFile(c.outputFile, stateContent, 0644); err != nil {
			c.Ui.Error(fmt.Sprintf("Error writing state file: %s", err))
			return 1
		}

		c.Ui.Info(fmt.Sprintf("State file downloaded successfully to: %s", c.outputFile))
		c.Ui.Info(fmt.Sprintf("State version: %s", stateVersion.ID))
		c.Ui.Info(fmt.Sprintf("Serial: %d", stateVersion.Serial))
		c.Ui.Info(fmt.Sprintf("Resources: %d", len(stateVersion.Resources)))
	} else {
		// Otherwise, print to stdout
		fmt.Println(string(stateContent))
	}

	return 0
}

// Help returns help text for the state download command
func (c *StateDownloadCommand) Help() string {
	helpText := `
Usage: hcptf state download [options]

  Download state file contents as JSON.

Options:

  -org=<name>           Organization name (required with -workspace)
  -workspace=<name>     Workspace name (required with -org)
  -id=<state-version>   State version ID (optional, defaults to current)
  -output=<file>        Output file path (optional, prints to stdout if omitted)

Examples:

  # Print current state from workspace to stdout
  hcptf state download -org=my-org -workspace=my-workspace

  # Download current state to file
  hcptf state download -org=my-org -workspace=my-workspace -output=state.json

  # Pipe to jq for analysis
  hcptf state download -org=my-org -workspace=my-workspace | jq '.resources | length'
`
	return strings.TrimSpace(helpText)
}

func (c *StateDownloadCommand) stateService(client *client.Client) stateVersionReader {
	if c.stateDownloadSvc != nil {
		return c.stateDownloadSvc
	}
	return client.StateVersions
}

func (c *StateDownloadCommand) workspaceService(client *client.Client) workspaceReader {
	if c.workspaceDownloadSvc != nil {
		return c.workspaceDownloadSvc
	}
	return client.Workspaces
}

// Synopsis returns a short synopsis for the state download command
func (c *StateDownloadCommand) Synopsis() string {
	return "Download state file contents as JSON"
}
