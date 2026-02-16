package router

import (
	"context"
	"fmt"
	"strings"

	tfe "github.com/hashicorp/go-tfe"
)

// Router handles URL-like argument routing
type Router struct {
	client      *tfe.Client
	commandTree *CommandTree
}

// NewRouter creates a new router
func NewRouter(client *tfe.Client, commandPaths []string) *Router {
	return &Router{
		client:      client,
		commandTree: NewCommandTree(commandPaths),
	}
}

// TranslateArgs converts URL-like args to command args
// Examples:
//   - "myorg" -> ["organization", "show", "-org=myorg"]
//   - "myorg workspaces" -> ["workspace", "list", "-org=myorg"]
//   - "myorg myworkspace" -> ["workspace", "read", "-org=myorg", "-workspace=myworkspace"]
//   - "myorg myworkspace runs list" -> ["run", "list", "-org=myorg", "-workspace=myworkspace"]
//   - "myorg -h" -> ["organization:context", "-org=myorg"]
//   - "myorg myworkspace -h" -> ["workspace:context", "-org=myorg", "-workspace=myworkspace"]
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

	// Check if help is requested at any position
	hasHelp := r.hasHelpFlag(args)

	// Just org: show org details or org context help
	if len(args) == 1 {
		return []string{"organization", "show", "-name=" + org}, nil
	}

	// org -h: show org context
	if len(args) == 2 && hasHelp {
		return []string{"organization:context", "-org=" + org}, nil
	}

	// Check for resource type as second arg
	second := args[1]

	// Handle known subcommands that list resources in an org
	if namespace, ok := r.commandTree.OrgCollectionNamespace(second); ok {
		if len(args) == 2 {
			return []string{namespace, "list", "-org=" + org}, nil
		}

		third := args[2]
		if hasHelp && (third == "-h" || third == "--help" || third == "-help") {
			return []string{namespace, "-h"}, nil
		}
		if third == "list" {
			return appendRemaining([]string{namespace, "list", "-org=" + org}, args, 3), nil
		}

		// Forward additional namespace actions while preserving URL-style org context.
		return appendRemaining([]string{namespace, third, "-org=" + org}, args, 3), nil
	}

	switch second {
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
		return []string{"workspace", "read", "-org=" + org, "-name=" + workspace}, nil
	}

	// org workspace -h: show workspace context
	if len(args) == 3 && hasHelp && !r.isResourceKeyword(args[1]) {
		workspace := args[1]
		return []string{"workspace:context", "-org=" + org, "-workspace=" + workspace}, nil
	}

	// 3+ args: org, workspace, resource/run-id
	if len(args) >= 3 {
		workspace := args[1]
		third := args[2]

		// Check if third arg is a run ID (format: run-xxx)
		if strings.HasPrefix(third, "run-") {
			runID := third
			if len(args) == 3 {
				return appendRemaining([]string{"run", "show", "-id=" + runID}, args, 3), nil
			}

			if len(args) >= 4 {
				return r.translateRunAction(org, workspace, runID, args[3], args, 4), nil
			}
		}

		switch third {
		case "runs":
			if len(args) == 3 {
				return appendRemaining([]string{"run", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) == 4 && args[3] == "list" {
				return appendRemaining([]string{"run", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
			// run show/apply/etc would be: org workspace runs <run-id> <action>
			if len(args) >= 4 {
				runID := args[3]
				if len(args) == 4 {
					return appendRemaining([]string{"run", "show", "-id=" + runID}, args, 4), nil
				}

				if len(args) >= 5 {
					return r.translateRunAction(org, workspace, runID, args[4], args, 5), nil
				}
			}
		case "variables":
			if len(args) == 3 {
				return appendRemaining([]string{"variable", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"variable", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		case "state":
			if len(args) == 3 {
				return appendRemaining([]string{"state", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"state", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
			if len(args) == 4 && args[3] == "outputs" {
				return appendRemaining([]string{"state", "outputs", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		case "resources":
			if len(args) == 3 {
				return appendRemaining([]string{"workspaceresource", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"workspaceresource", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		case "assessments":
			if len(args) == 3 {
				return appendRemaining([]string{"assessmentresult", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"assessmentresult", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		case "changerequests":
			if len(args) == 3 {
				return appendRemaining([]string{"changerequest", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"changerequest", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		case "configversions":
			if len(args) == 3 {
				return appendRemaining([]string{"configversion", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"configversion", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		case "tags":
			if len(args) == 3 {
				return appendRemaining([]string{"workspacetag", "list", "-org=" + org, "-workspace=" + workspace}, args, 3), nil
			}
			if len(args) >= 4 && args[3] == "list" {
				return appendRemaining([]string{"workspacetag", "list", "-org=" + org, "-workspace=" + workspace}, args, 4), nil
			}
		}
	}

	// If we couldn't translate, return as-is
	return args, nil
}

// isKnownCommand checks if the arg is a known command
func (r *Router) isKnownCommand(arg string) bool {
	if r.commandTree == nil {
		return false
	}
	return r.commandTree.HasRoot(arg)
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

// translateRunAction maps a run sub-resource action to the appropriate command args.
// This is shared between the short-form (org workspace run-xxx action) and
// long-form (org workspace runs run-xxx action) dispatch paths.
func (r *Router) translateRunAction(org, workspace, runID, action string, args []string, consumed int) []string {
	switch action {
	case "plan":
		return appendRemaining([]string{"plan", "read", "-id=" + runID}, args, consumed)
	case "logs", "planlogs":
		return appendRemaining([]string{"plan", "logs", "-id=" + runID}, args, consumed)
	case "applylogs":
		return appendRemaining([]string{"apply", "logs", "-id=" + runID}, args, consumed)
	case "applyread", "applydetails":
		return appendRemaining([]string{"apply", "read", "-id=" + runID}, args, consumed)
	case "comments":
		return appendRemaining([]string{"comment", "list", "-run-id=" + runID}, args, consumed)
	case "policychecks":
		return appendRemaining([]string{"policycheck", "list", "-run-id=" + runID}, args, consumed)
	case "state", "stateversions":
		return appendRemaining([]string{"state", "list", "-org=" + org, "-workspace=" + workspace}, args, consumed)
	case "outputs":
		return appendRemaining([]string{"state", "outputs", "-org=" + org, "-workspace=" + workspace}, args, consumed)
	case "configversion":
		return appendRemaining([]string{"configversion", "read", "-run-id=" + runID}, args, consumed)
	case "assessment":
		return appendRemaining([]string{"assessmentresult", "list", "-org=" + org, "-workspace=" + workspace}, args, consumed)
	default:
		return appendRemaining([]string{"run", action, "-id=" + runID}, args, consumed)
	}
}

// hasHelpFlag checks if help flag is present in args
func (r *Router) hasHelpFlag(args []string) bool {
	for _, arg := range args {
		if arg == "-h" || arg == "--help" || arg == "-help" {
			return true
		}
	}
	return false
}

func (r *Router) isResourceKeyword(arg string) bool {
	if r.commandTree == nil {
		return false
	}
	return r.commandTree.IsResourceKeyword(arg)
}

func appendRemaining(base []string, args []string, consumed int) []string {
	if consumed >= len(args) {
		return base
	}
	return append(base, args[consumed:]...)
}
