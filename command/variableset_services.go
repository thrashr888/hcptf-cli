package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type variableSetReader interface {
	Read(ctx context.Context, variableSetID string, options *tfe.VariableSetReadOptions) (*tfe.VariableSet, error)
}

type variableSetDeleter interface {
	Delete(ctx context.Context, variableSetID string) error
}

type variableSetLister interface {
	List(ctx context.Context, organization string, options *tfe.VariableSetListOptions) (*tfe.VariableSetList, error)
}

type variableSetCreator interface {
	Create(ctx context.Context, organization string, options *tfe.VariableSetCreateOptions) (*tfe.VariableSet, error)
}

type variableSetWorkspaceLister interface {
	ListForWorkspace(ctx context.Context, workspaceID string, options *tfe.VariableSetListOptions) (*tfe.VariableSetList, error)
}

type variableSetProjectLister interface {
	ListForProject(ctx context.Context, projectID string, options *tfe.VariableSetListOptions) (*tfe.VariableSetList, error)
}

type variableSetRemover interface {
	RemoveFromWorkspaces(ctx context.Context, variableSetID string, options *tfe.VariableSetRemoveFromWorkspacesOptions) error
	RemoveFromProjects(ctx context.Context, variableSetID string, options tfe.VariableSetRemoveFromProjectsOptions) error
	RemoveFromStacks(ctx context.Context, variableSetID string, options *tfe.VariableSetRemoveFromStacksOptions) error
}

type variableSetWorkspaceUpdater interface {
	UpdateWorkspaces(ctx context.Context, variableSetID string, options *tfe.VariableSetUpdateWorkspacesOptions) (*tfe.VariableSet, error)
}

type variableSetStackUpdater interface {
	UpdateStacks(ctx context.Context, variableSetID string, options *tfe.VariableSetUpdateStacksOptions) (*tfe.VariableSet, error)
}
