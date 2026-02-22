package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/jsonapi"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/hcptf-cli/internal/client"
)

// WorkspaceCreateCommand is a command to create a workspace
type WorkspaceCreateCommand struct {
	Meta
	organization string
	name         string
	description  string
	projectID    string
	format       string
	workspaceSvc workspaceCreator

	// Basic settings
	terraformVersion     string
	executionMode        string
	workingDirectory     string
	agentPoolID          string
	autoApply            bool
	sourceName           string
	sourceURL            string
	migrationEnvironment string

	// Boolean toggle flags (string "true"/"false", empty = unset)
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

	// Tags (comma-separated)
	tags string

	// VCS flags
	vcsIdentifier        string
	vcsBranch            string
	vcsOAuthTokenID      string
	vcsIngressSubmodules string
	vcsTagsRegex         string
	vcsGHAInstallationID string

	// Auto-destroy
	autoDestroyAt               string
	autoDestroyActivityDuration string
}

// Run executes the workspace create command
func (c *WorkspaceCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet("workspace create")

	// Required
	flags.StringVar(&c.organization, "organization", "", "Organization name (required)")
	flags.StringVar(&c.organization, "org", "", "Organization name (alias)")
	flags.StringVar(&c.name, "name", "", "Workspace name (required)")

	// Basic settings
	flags.StringVar(&c.description, "description", "", "Workspace description")
	flags.StringVar(&c.terraformVersion, "terraform-version", "", "Terraform version")
	flags.StringVar(&c.executionMode, "execution-mode", "", "Execution mode: remote, local, or agent")
	flags.StringVar(&c.workingDirectory, "working-directory", "", "Working directory")
	flags.StringVar(&c.projectID, "project-id", "", "Project ID to assign the workspace to")
	flags.StringVar(&c.agentPoolID, "agent-pool-id", "", "Agent pool ID (required when execution-mode is agent)")
	flags.BoolVar(&c.autoApply, "auto-apply", false, "Enable auto-apply")
	flags.StringVar(&c.sourceName, "source-name", "", "Source name for workspace creation tracking")
	flags.StringVar(&c.sourceURL, "source-url", "", "Source URL for workspace creation tracking")
	flags.StringVar(&c.migrationEnvironment, "migration-environment", "", "Legacy TFE environment for migration")

	// Boolean toggles
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
	flags.StringVar(&c.hyokEnabled, "hyok-enabled", "", "Enable HYOK (true/false)")

	// VCS options
	flags.StringVar(&c.vcsIdentifier, "vcs-identifier", "", "VCS repository identifier (e.g. org/repo)")
	flags.StringVar(&c.vcsBranch, "vcs-branch", "", "VCS repository branch")
	flags.StringVar(&c.vcsOAuthTokenID, "vcs-oauth-token-id", "", "VCS OAuth token ID")
	flags.StringVar(&c.vcsIngressSubmodules, "vcs-ingress-submodules", "", "Enable VCS ingress submodules (true/false)")
	flags.StringVar(&c.vcsTagsRegex, "vcs-tags-regex", "", "VCS tags regex")
	flags.StringVar(&c.vcsGHAInstallationID, "vcs-gha-installation-id", "", "GitHub App installation ID")

	// Advanced
	flags.StringVar(&c.triggerPrefixes, "trigger-prefixes", "", "Comma-separated list of trigger prefixes")
	flags.StringVar(&c.triggerPatterns, "trigger-patterns", "", "Comma-separated list of trigger patterns")
	flags.StringVar(&c.tags, "tags", "", "Comma-separated list of tags")
	flags.StringVar(&c.autoDestroyAt, "auto-destroy-at", "", "Auto-destroy time (RFC3339)")
	flags.StringVar(&c.autoDestroyActivityDuration, "auto-destroy-activity-duration", "", "Auto-destroy activity duration (e.g. '24h')")

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

	// Validate execution mode
	if c.executionMode != "" {
		switch c.executionMode {
		case "remote", "local", "agent":
		default:
			c.Ui.Error("Error: -execution-mode must be 'remote', 'local', or 'agent'")
			return 1
		}
	}

	// Get API client
	client, err := c.Meta.Client()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error initializing client: %s", err))
		return 1
	}

	// Create workspace
	options := tfe.WorkspaceCreateOptions{
		Name:      tfe.String(c.name),
		AutoApply: tfe.Bool(c.autoApply),
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

	if c.sourceName != "" {
		options.SourceName = tfe.String(c.sourceName)
	}

	if c.sourceURL != "" {
		options.SourceURL = tfe.String(c.sourceURL)
	}

	if c.migrationEnvironment != "" {
		options.MigrationEnvironment = tfe.String(c.migrationEnvironment)
	}

	// Parse boolean toggle flags
	boolFlags := []struct {
		value    string
		flagName string
		target   **bool
	}{
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
		{c.hyokEnabled, "hyok-enabled", &options.HYOKEnabled},
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

	// Trigger prefixes/patterns
	if c.triggerPrefixes != "" {
		options.TriggerPrefixes = splitCommaList(c.triggerPrefixes)
	}
	if c.triggerPatterns != "" {
		options.TriggerPatterns = splitCommaList(c.triggerPatterns)
	}

	// Tags
	if c.tags != "" {
		tagNames := splitCommaList(c.tags)
		tags := make([]*tfe.Tag, len(tagNames))
		for i, name := range tagNames {
			tags[i] = &tfe.Tag{Name: name}
		}
		options.Tags = tags
	}

	// VCS repo options
	if c.vcsIdentifier != "" || c.vcsBranch != "" || c.vcsOAuthTokenID != "" ||
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
		t, parseErr := time.Parse(time.RFC3339, c.autoDestroyAt)
		if parseErr != nil {
			c.Ui.Error(fmt.Sprintf("Error: -auto-destroy-at must be a valid RFC3339 time: %s", parseErr))
			return 1
		}
		options.AutoDestroyAt = tfe.NullableTime(t)
	}

	// Auto-destroy-activity-duration
	if c.autoDestroyActivityDuration != "" {
		options.AutoDestroyActivityDuration = jsonapi.NewNullableAttrWithValue(c.autoDestroyActivityDuration)
	}

	workspace, err := c.workspaceService(client).Create(client.Context(), c.organization, options)
	if err != nil {
		c.Ui.Error(fmt.Sprintf("Error creating workspace: %s", err))
		return 1
	}

	// Format output
	formatter := c.Meta.NewFormatter(c.format)

	if c.format != "json" {
		c.Ui.Output(fmt.Sprintf("Workspace '%s' created successfully", workspace.Name))
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
		"TagNames":            strings.Join(workspace.TagNames, ","),
		"CreatedAt":           workspace.CreatedAt,
	}

	formatter.KeyValue(data)
	return 0
}

func (c *WorkspaceCreateCommand) workspaceService(client *client.Client) workspaceCreator {
	if c.workspaceSvc != nil {
		return c.workspaceSvc
	}
	return client.Workspaces
}

// Help returns help text for the workspace create command
func (c *WorkspaceCreateCommand) Help() string {
	helpText := `
Usage: hcptf workspace create [options]

  Create a new workspace.

Required:

  -organization=<name>              Organization name (required)
  -org=<name>                       Alias for -organization
  -name=<name>                      Workspace name (required)

Basic Settings:

  -description=<text>               Workspace description
  -terraform-version=<ver>          Terraform version to use
  -execution-mode=<mode>            Execution mode: remote, local, or agent
  -working-directory=<dir>          Working directory for Terraform operations
  -project-id=<id>                  Project ID to assign the workspace to
  -agent-pool-id=<id>               Agent pool ID (required when execution-mode is agent)
  -auto-apply                       Enable auto-apply (default: false)
  -source-name=<name>               Source name for workspace creation tracking
  -source-url=<url>                 Source URL for workspace creation tracking
  -migration-environment=<env>      Legacy TFE environment for migration

Boolean Toggles (true/false):

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
  -hyok-enabled=<bool>              Enable HYOK

VCS Options:

  -vcs-identifier=<id>              VCS repository identifier (e.g. org/repo)
  -vcs-branch=<branch>              VCS repository branch
  -vcs-oauth-token-id=<id>          VCS OAuth token ID
  -vcs-ingress-submodules=<bool>    Enable VCS ingress submodules (true/false)
  -vcs-tags-regex=<regex>           VCS tags regex
  -vcs-gha-installation-id=<id>     GitHub App installation ID

Advanced:

  -trigger-prefixes=<list>          Comma-separated list of trigger prefixes
  -trigger-patterns=<list>          Comma-separated list of trigger patterns
  -tags=<list>                      Comma-separated list of tags
  -auto-destroy-at=<time>           Auto-destroy time (RFC3339)
  -auto-destroy-activity-duration=<dur>  Auto-destroy activity duration (e.g. '24h')

Output:

  -output=<format>                  Output format: table (default) or json

Example:

  hcptf workspace create -org=my-org -name=my-workspace
  hcptf workspace create -org=my-org -name=prod -auto-apply -terraform-version=1.5.0
  hcptf workspace create -org=my-org -name=prod -project-id=prj-abc123
  hcptf workspace create -org=my-org -name=prod -execution-mode=agent -agent-pool-id=apool-123
  hcptf workspace create -org=my-org -name=prod -vcs-identifier=org/repo -vcs-oauth-token-id=ot-123
  hcptf workspace create -org=my-org -name=prod -tags=env:prod,team:infra
`
	return strings.TrimSpace(helpText)
}

// Synopsis returns a short synopsis for the workspace create command
func (c *WorkspaceCreateCommand) Synopsis() string {
	return "Create a new workspace"
}
