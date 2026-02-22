package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/jsonapi"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// WorkspaceUpdateCommand is a command to update a workspace
type WorkspaceUpdateCommand struct {
	Meta
	organization string
	name         string
	newName      string
	description  string
	projectID    string
	format       string
	workspaceSvc workspaceUpdater

	// Basic settings
	terraformVersion string
	executionMode    string
	workingDirectory string
	agentPoolID      string

	// Boolean toggle flags (string "true"/"false", empty = unset)
	autoApply                  string
	allowDestroyPlan           string
	assessmentsEnabled         string
	autoApplyRunTrigger        string
	fileTriggersEnabled        string
	globalRemoteState          string
	projectRemoteState         string
	queueAllRuns               string
	speculativeEnabled         string
	structuredRunOutputEnabled string
	inheritsProjectAutoDestroy string
	hyokEnabled                string

	// Trigger lists (comma-separated)
	triggerPrefixes string
	triggerPatterns string

	// VCS flags
	vcsIdentifier        string
	vcsBranch            string
	vcsOAuthTokenID      string
	vcsIngressSubmodules string
	vcsTagsRegex         string
	vcsGHAInstallationID string
	removeVCS            bool

	// Auto-destroy
	autoDestroyAt               string
	autoDestroyActivityDuration string
}

// Run executes the workspace update command
func (c *WorkspaceUpdateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace update")

	// Required
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")

	// Basic settings
	flags.StringVar(&c.newName, "new-name", "", "New workspace name")
	flags.StringVar(&c.description, "description", "", "Workspace description")
	flags.StringVar(&c.terraformVersion, "terraform-version", "", "Terraform version")
	flags.StringVar(&c.executionMode, "execution-mode", "", "Execution mode: remote, local, or agent")
	flags.StringVar(&c.workingDirectory, "working-directory", "", "Working directory")
	flags.StringVar(&c.projectID, "project-id", "", "Project ID to move the workspace into")
	flags.StringVar(&c.agentPoolID, "agent-pool-id", "", "Agent pool ID (required when execution-mode is agent)")

	// Boolean toggles
	flags.StringVar(&c.autoApply, "auto-apply", "", "Enable auto-apply (true/false)")
	flags.StringVar(&c.allowDestroyPlan, "allow-destroy-plan", "", "Allow destroy plans (true/false)")
	flags.StringVar(&c.assessmentsEnabled, "assessments-enabled", "", "Enable health assessments (true/false)")
	flags.StringVar(&c.autoApplyRunTrigger, "auto-apply-run-trigger", "", "Auto-apply run triggers (true/false)")
	flags.StringVar(&c.fileTriggersEnabled, "file-triggers-enabled", "", "Enable file triggers (true/false)")
	flags.StringVar(&c.globalRemoteState, "global-remote-state", "", "Enable global remote state (true/false)")
	flags.StringVar(&c.projectRemoteState, "project-remote-state", "", "Enable project remote state (true/false)")
	flags.StringVar(&c.queueAllRuns, "queue-all-runs", "", "Queue all runs (true/false)")
	flags.StringVar(&c.speculativeEnabled, "speculative-enabled", "", "Enable speculative plans (true/false)")
	flags.StringVar(&c.structuredRunOutputEnabled, "structured-run-output-enabled", "", "Enable structured run output (true/false)")
	flags.StringVar(&c.inheritsProjectAutoDestroy, "inherits-project-auto-destroy", "", "Inherit project auto-destroy settings (true/false)")
	flags.StringVar(&c.hyokEnabled, "hyok-enabled", "", "Enable HYOK (true only; cannot be disabled)")

	// VCS options
	flags.StringVar(&c.vcsIdentifier, "vcs-identifier", "", "VCS repository identifier (e.g. org/repo)")
	flags.StringVar(&c.vcsBranch, "vcs-branch", "", "VCS repository branch")
	flags.StringVar(&c.vcsOAuthTokenID, "vcs-oauth-token-id", "", "VCS OAuth token ID")
	flags.StringVar(&c.vcsIngressSubmodules, "vcs-ingress-submodules", "", "Enable VCS ingress submodules (true/false)")
	flags.StringVar(&c.vcsTagsRegex, "vcs-tags-regex", "", "VCS tags regex")
	flags.StringVar(&c.vcsGHAInstallationID, "vcs-gha-installation-id", "", "GitHub App installation ID")
	flags.BoolVar(&c.removeVCS, "remove-vcs", false, "Remove VCS connection from workspace")

	// Advanced
	flags.StringVar(&c.triggerPrefixes, "trigger-prefixes", "", "Comma-separated list of trigger prefixes")
	flags.StringVar(&c.triggerPatterns, "trigger-patterns", "", "Comma-separated list of trigger patterns")
	flags.StringVar(&c.autoDestroyAt, "auto-destroy-at", "", "Auto-destroy time (RFC3339) or 'none' to clear")
	flags.StringVar(&c.autoDestroyActivityDuration, "auto-destroy-activity-duration", "", "Auto-destroy activity duration (e.g. '24h') or 'none' to clear")

	// Output
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

	if c.name == "" {
		c.Ui.Error("Error: -name flag is required")
		c.Ui.Error(c.Help())
		return 1
	}

	// Check for conflict: -remove-vcs cannot be used with any -vcs-* flag
	if c.removeVCS {
		if c.vcsIdentifier != "" || c.vcsBranch != "" || c.vcsOAuthTokenID != "" ||
			c.vcsIngressSubmodules != "" || c.vcsTagsRegex != "" || c.vcsGHAInstallationID != "" {
			c.Ui.Error("Error: -remove-vcs cannot be used together with -vcs-* flags")
			return 1
		}
	}

	// Validate execution mode
	if c.executionMode != "" {
		switch c.executionMode {
		case "remote", "local", "agent":
		default:
			c.Ui.Error("Error: -execution-mode must be 'remote', 'local', or 'agent'")
			return 1
		}
	}

	// Validate HYOK: only "true" is allowed
	if c.hyokEnabled != "" {
		if c.hyokEnabled == "false" {
			c.Ui.Error("Error: HYOK cannot be disabled once enabled")
			return 1
		}
		if c.hyokEnabled != "true" {
			c.Ui.Error("Error: -hyok-enabled must be 'true'")
			return 1
		}
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Build update options
	options := tfe.WorkspaceUpdateOptions{}

	if c.newName != "" {
		options.Name = tfe.String(c.newName)
	}

	if c.description != "" {
		options.Description = tfe.String(c.description)
	}

	if c.terraformVersion != "" {
		options.TerraformVersion = tfe.String(c.terraformVersion)
	}

	if c.executionMode != "" {
		options.ExecutionMode = tfe.String(c.executionMode)
	}

	if c.workingDirectory != "" {
		options.WorkingDirectory = tfe.String(c.workingDirectory)
	}

	if c.agentPoolID != "" {
		options.AgentPoolID = tfe.String(c.agentPoolID)
	}

	if c.projectID != "" {
		options.Project = &tfe.Project{ID: c.projectID}
	}

	// Parse boolean toggle flags
	boolFlags := []struct {
		value    string
		flagName string
		target   **bool
	}{
		{c.autoApply, "auto-apply", &options.AutoApply},
		{c.allowDestroyPlan, "allow-destroy-plan", &options.AllowDestroyPlan},
		{c.assessmentsEnabled, "assessments-enabled", &options.AssessmentsEnabled},
		{c.autoApplyRunTrigger, "auto-apply-run-trigger", &options.AutoApplyRunTrigger},
		{c.fileTriggersEnabled, "file-triggers-enabled", &options.FileTriggersEnabled},
		{c.globalRemoteState, "global-remote-state", &options.GlobalRemoteState},
		{c.projectRemoteState, "project-remote-state", &options.ProjectRemoteState},
		{c.queueAllRuns, "queue-all-runs", &options.QueueAllRuns},
		{c.speculativeEnabled, "speculative-enabled", &options.SpeculativeEnabled},
		{c.structuredRunOutputEnabled, "structured-run-output-enabled", &options.StructuredRunOutputEnabled},
		{c.inheritsProjectAutoDestroy, "inherits-project-auto-destroy", &options.InheritsProjectAutoDestroy},
	}

	for _, bf := range boolFlags {
		parsed, parseErr := parseBoolFlag(bf.value, bf.flagName)
		if parseErr != nil {
			c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
			return 1
		}
		if parsed != nil {
			*bf.target = parsed
		}
	}

	// HYOK (already validated above)
	if c.hyokEnabled == "true" {
		options.HYOKEnabled = tfe.Bool(true)
	}

	// Trigger prefixes/patterns
	if c.triggerPrefixes != "" {
		options.TriggerPrefixes = splitCommaList(c.triggerPrefixes)
	}
	if c.triggerPatterns != "" {
		options.TriggerPatterns = splitCommaList(c.triggerPatterns)
	}

	// VCS repo options
	if c.removeVCS {
		options.VCSRepo = &tfe.VCSRepoOptions{}
	} else if c.vcsIdentifier != "" || c.vcsBranch != "" || c.vcsOAuthTokenID != "" ||
		c.vcsIngressSubmodules != "" || c.vcsTagsRegex != "" || c.vcsGHAInstallationID != "" {
		vcsRepo := &tfe.VCSRepoOptions{}
		if c.vcsIdentifier != "" {
			vcsRepo.Identifier = tfe.String(c.vcsIdentifier)
		}
		if c.vcsBranch != "" {
			vcsRepo.Branch = tfe.String(c.vcsBranch)
		}
		if c.vcsOAuthTokenID != "" {
			vcsRepo.OAuthTokenID = tfe.String(c.vcsOAuthTokenID)
		}
		if c.vcsIngressSubmodules != "" {
			parsed, parseErr := parseBoolFlag(c.vcsIngressSubmodules, "vcs-ingress-submodules")
			if parseErr != nil {
				c.Ui.Error(fmt.Sprintf("Error: %s", parseErr))
				return 1
			}
			vcsRepo.IngressSubmodules = parsed
		}
		if c.vcsTagsRegex != "" {
			vcsRepo.TagsRegex = tfe.String(c.vcsTagsRegex)
		}
		if c.vcsGHAInstallationID != "" {
			vcsRepo.GHAInstallationID = tfe.String(c.vcsGHAInstallationID)
		}
		options.VCSRepo = vcsRepo
	}

	// Auto-destroy-at
	if c.autoDestroyAt != "" {
		if c.autoDestroyAt == "none" {
			options.AutoDestroyAt = tfe.NullTime()
		} else {
			t, parseErr := time.Parse(time.RFC3339, c.autoDestroyAt)
			if parseErr != nil {
				c.Ui.Error(fmt.Sprintf("Error: -auto-destroy-at must be a valid RFC3339 time or 'none': %s", parseErr))
				return 1
			}
			options.AutoDestroyAt = tfe.NullableTime(t)
		}
	}

	// Auto-destroy-activity-duration
	if c.autoDestroyActivityDuration != "" {
		if c.autoDestroyActivityDuration == "none" {
			options.AutoDestroyActivityDuration = jsonapi.NewNullNullableAttr[string]()
		} else {
			options.AutoDestroyActivityDuration = jsonapi.NewNullableAttrWithValue(c.autoDestroyActivityDuration)
		}
	}

	// Update workspace
	workspace, err := c.workspaceService(client).Update(client.Context(), c.organization, c.name, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error updating workspace: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Workspace '%s' updated successfully", workspace.Name))
	}

	// Show workspace details
	projectID := ""
	if workspace.Project != nil {
		projectID = workspace.Project.ID
	}

	agentPoolID := ""
	if workspace.AgentPool != nil {
		agentPoolID = workspace.AgentPool.ID
	}

	vcsRepoID := ""
	if workspace.VCSRepo != nil {
		vcsRepoID = workspace.VCSRepo.Identifier
	}

	data := map[string]interface{}{
		"ID":                  workspace.ID,
		"Name":                workspace.Name,
		"Organization":        c.organization,
		"ProjectID":           projectID,
		"TerraformVersion":    workspace.TerraformVersion,
		"AutoApply":           workspace.AutoApply,
		"Description":         workspace.Description,
		"ExecutionMode":       workspace.ExecutionMode,
		"WorkingDirectory":    workspace.WorkingDirectory,
		"AgentPoolID":         agentPoolID,
		"AllowDestroyPlan":    workspace.AllowDestroyPlan,
		"FileTriggersEnabled": workspace.FileTriggersEnabled,
		"QueueAllRuns":        workspace.QueueAllRuns,
		"SpeculativeEnabled":  workspace.SpeculativeEnabled,
		"TriggerPrefixes":     strings.Join(workspace.TriggerPrefixes, ","),
		"TriggerPatterns":     strings.Join(workspace.TriggerPatterns, ","),
		"VCSRepo":             vcsRepoID,
		"UpdatedAt":           workspace.UpdatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *WorkspaceUpdateCommand) workspaceService(client *client.Client) workspaceUpdater {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace update command
func (c *WorkspaceUpdateCommand) Help() string {
	helpText := `
Usage: hcptf workspace update [options]

  Update workspace settings.

Required:

  -organization=<name>              Organization name (required)
  -org=<name>                       Alias for -organization
  -name=<name>                      Workspace name (required)

Basic Settings:

  -new-name=<name>                  New workspace name
  -description=<text>               Workspace description
  -terraform-version=<ver>          Terraform version to use
  -execution-mode=<mode>            Execution mode: remote, local, or agent
  -working-directory=<dir>          Working directory for Terraform operations
  -project-id=<id>                  Project ID to move the workspace into
  -agent-pool-id=<id>               Agent pool ID (required when execution-mode is agent)

Boolean Toggles (true/false):

  -auto-apply=<bool>                Enable auto-apply
  -allow-destroy-plan=<bool>        Allow destroy plans
  -assessments-enabled=<bool>       Enable health assessments (drift detection)
  -auto-apply-run-trigger=<bool>    Auto-apply run triggers
  -file-triggers-enabled=<bool>     Enable file triggers
  -global-remote-state=<bool>       Enable global remote state
  -project-remote-state=<bool>      Enable project remote state
  -queue-all-runs=<bool>            Queue all runs
  -speculative-enabled=<bool>       Enable speculative plans
  -structured-run-output-enabled=<bool>  Enable structured run output
  -inherits-project-auto-destroy=<bool>  Inherit project auto-destroy settings
  -hyok-enabled=<bool>              Enable HYOK (true only; cannot be disabled)

VCS Options:

  -vcs-identifier=<id>              VCS repository identifier (e.g. org/repo)
  -vcs-branch=<branch>              VCS repository branch
  -vcs-oauth-token-id=<id>          VCS OAuth token ID
  -vcs-ingress-submodules=<bool>    Enable VCS ingress submodules (true/false)
  -vcs-tags-regex=<regex>           VCS tags regex
  -vcs-gha-installation-id=<id>     GitHub App installation ID
  -remove-vcs                       Remove VCS connection (cannot be used with -vcs-* flags)

Advanced:

  -trigger-prefixes=<list>          Comma-separated list of trigger prefixes
  -trigger-patterns=<list>          Comma-separated list of trigger patterns
  -auto-destroy-at=<time>           Auto-destroy time (RFC3339) or 'none' to clear
  -auto-destroy-activity-duration=<dur>  Auto-destroy activity duration (e.g. '24h') or 'none' to clear

Output:

  -output=<format>                  Output format: table (default) or json

Example:

  hcptf workspace update -org=my-org -name=my-workspace -auto-apply=true
  hcptf workspace update -org=my-org -name=old-name -new-name=new-name
  hcptf workspace update -org=my-org -name=my-workspace -project-id=prj-abc123
  hcptf workspace update -org=my-org -name=my-workspace -execution-mode=agent -agent-pool-id=apool-123
  hcptf workspace update -org=my-org -name=my-workspace -vcs-identifier=org/repo -vcs-oauth-token-id=ot-123
  hcptf workspace update -org=my-org -name=my-workspace -remove-vcs
  hcptf workspace update -org=my-org -name=my-workspace -auto-destroy-at=none
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace update command
func (c *WorkspaceUpdateCommand) Synopsis() string {
	return "Update workspace settings"
}

// parseBoolFlag parses a string flag that accepts "true" or "false".
// Returns (value, error). If the input is empty, returns (nil, nil).
func parseBoolFlag(value, flagName string) (*bool, error) {
	if value == "" {
		return nil, nil
	}
	switch value {
	case "true":
		return tfe.Bool(true), nil
	case "false":
		return tfe.Bool(false), nil
	default:
		return nil, fmt.Errorf("-%s must be 'true' or 'false'", flagName)
	}
}

// splitCommaList splits a comma-separated string into a slice,
// trimming whitespace from each element.
func splitCommaList(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
