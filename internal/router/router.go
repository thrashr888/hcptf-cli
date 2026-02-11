package router

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// Router handles URL-like argument routing
type Router struct {
	client *tfe.Client
}

// NewRouter creates a new router
func NewRouter(client *tfe.Client) *Router {
	return &Router{client: client}
}

// TranslateArgs converts URL-like args to command args
// Examples:
//   - "myorg" -> ["organization", "show", "-org=myorg"]
//   - "myorg workspaces" -> ["workspace", "list", "-org=myorg"]
//   - "myorg myworkspace" -> ["workspace", "read", "-org=myorg", "-workspace=myworkspace"]
//   - "myorg myworkspace runs list" -> ["run", "list", "-org=myorg", "-workspace=myworkspace"]
func (r *Router) TranslateArgs(args []string) ([]string, error) {
	// If no args or first arg starts with "-", use default routing
	if len(args) == 0 || strings.HasPrefix(args[0], "-") {
		return args, nil
	}

	// Check if first arg is a known command
	if r.isKnownCommand(args[0]) {
		return args, nil
	}

	// URL-like pattern detected
	org := args[0]

	// Just org: show org details
	if len(args) == 1 {
		return []string{"organization", "show", "-org=" + org}, nil
	}

	// Check for resource type as second arg
	second := args[1]

	// Handle known subcommands that list resources in an org
	switch second {
	case "workspaces":
		return []string{"workspace", "list", "-org=" + org}, nil
	case "projects":
		return []string{"project", "list", "-org=" + org}, nil
	case "teams":
		return []string{"team", "list", "-org=" + org}, nil
	case "policies":
		return []string{"policy", "list", "-org=" + org}, nil
	case "policysets":
		return []string{"policyset", "list", "-org=" + org}, nil
	case "variables":
		// This would need a workspace context, skip for now
		return args, nil
	case "runs":
		// This would need a workspace context, skip for now
		return args, nil
	}

	// If we have 2 args and second is not a known subcommand,
	// assume it's a workspace
	if len(args) == 2 {
		workspace := args[1]
		return []string{"workspace", "read", "-org=" + org, "-workspace=" + workspace}, nil
	}

	// 3+ args: org, workspace, resource
	if len(args) >= 3 {
		workspace := args[1]
		third := args[2]

		switch third {
		case "runs":
			if len(args) == 3 || (len(args) == 4 && args[3] == "list") {
				return []string{"run", "list", "-org=" + org, "-workspace=" + workspace}, nil
			}
			// run show/apply/etc would be: org workspace runs <run-id> <action>
			if len(args) >= 4 {
				runID := args[3]
				action := "show"
				if len(args) >= 5 {
					action = args[4]
				}
				return []string{"run", action, "-id=" + runID}, nil
			}
		case "variables":
			if len(args) == 3 || (len(args) == 4 && args[3] == "list") {
				return []string{"variable", "list", "-org=" + org, "-workspace=" + workspace}, nil
			}
		case "state":
			if len(args) == 3 || (len(args) == 4 && args[3] == "list") {
				return []string{"state", "list", "-org=" + org, "-workspace=" + workspace}, nil
			}
			if len(args) == 4 && args[3] == "outputs" {
				return []string{"state", "outputs", "-org=" + org, "-workspace=" + workspace}, nil
			}
		}
	}

	// If we couldn't translate, return as-is
	return args, nil
}

// isKnownCommand checks if the arg is a known command
func (r *Router) isKnownCommand(arg string) bool {
	knownCommands := []string{
		"account", "login", "logout", "version",
		"workspace", "run", "organization", "variable", "team", "project",
		"state", "policy", "policyset", "sshkey", "notification",
		"variableset", "agentpool", "runtask", "oauthclient", "oauthtoken",
		"runtrigger", "plan", "apply", "configversion", "teamaccess",
		"projectteamaccess", "registrymodule", "registryprovider",
		"registryproviderversion", "registryproviderplatform", "gpgkey",
		"stack", "stackconfiguration", "stackdeployment", "stackstate",
		"audittrail", "audittrailtoken", "organizationtoken", "usertoken",
		"teamtoken", "organizationmembership", "organizationmember",
		"organizationtag", "reservedtagkey", "comment", "policycheck",
		"policyevaluation", "policysetoutcome", "policysetparameter",
		"awsoidc", "azureoidc", "gcpoidc", "vaultoidc",
		"workspaceresource", "workspacetag", "queryrun", "queryworkspace",
		"changerequest", "assessmentresult", "hyok", "hyokkey",
		"vcsevent", "planexport", "agent",
	}

	for _, cmd := range knownCommands {
		if arg == cmd {
			return true
		}
	}
	return false
}

// ValidateOrg checks if the org exists (optional, for better UX)
func (r *Router) ValidateOrg(ctx context.Context, org string) error {
	if r.client == nil {
		return nil // Skip validation if no client
	}

	_, err := r.client.Organizations.Read(ctx, org)
	if err != nil {
		return fmt.Errorf("organization %q not found: %w", org, err)
	}
	return nil
}

// ValidateWorkspace checks if the workspace exists (optional, for better UX)
func (r *Router) ValidateWorkspace(ctx context.Context, org, workspace string) error {
	if r.client == nil {
		return nil // Skip validation if no client
	}

	_, err := r.client.Workspaces.Read(ctx, org, workspace)
	if err != nil {
		return fmt.Errorf("workspace %q not found in organization %q: %w", workspace, org, err)
	}
	return nil
}
