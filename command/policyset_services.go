package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type policySetLister interface {
	List(ctx context.Context, organization string, options *tfe.PolicySetListOptions) (*tfe.PolicySetList, error)
}

type policySetReader interface {
	Read(ctx context.Context, policySetID string) (*tfe.PolicySet, error)
}

type policySetReaderWithOptions interface {
	ReadWithOptions(ctx context.Context, policySetID string, options *tfe.PolicySetReadOptions) (*tfe.PolicySet, error)
}

type policySetWorkspaceAdder interface {
	AddWorkspaces(ctx context.Context, policySetID string, options tfe.PolicySetAddWorkspacesOptions) error
}

type policySetWorkspaceRemover interface {
	RemoveWorkspaces(ctx context.Context, policySetID string, options tfe.PolicySetRemoveWorkspacesOptions) error
}

type policySetWorkspaceExclusionAdder interface {
	AddWorkspaceExclusions(ctx context.Context, policySetID string, options tfe.PolicySetAddWorkspaceExclusionsOptions) error
}

type policySetWorkspaceExclusionRemover interface {
	RemoveWorkspaceExclusions(ctx context.Context, policySetID string, options tfe.PolicySetRemoveWorkspaceExclusionsOptions) error
}

type policySetProjectAdder interface {
	AddProjects(ctx context.Context, policySetID string, options tfe.PolicySetAddProjectsOptions) error
}

type policySetProjectRemover interface {
	RemoveProjects(ctx context.Context, policySetID string, options tfe.PolicySetRemoveProjectsOptions) error
}
