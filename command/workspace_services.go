package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type workspaceLister interface {
	List(ctx context.Context, organization string, options *tfe.WorkspaceListOptions) (*tfe.WorkspaceList, error)
}

type workspaceReader interface {
	Read(ctx context.Context, organization, workspace string) (*tfe.Workspace, error)
}

type workspaceCreator interface {
	Create(ctx context.Context, organization string, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error)
}

type workspaceUpdater interface {
	Update(ctx context.Context, organization, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error)
}

type workspaceDeleter interface {
	Delete(ctx context.Context, organization, workspace string) error
}
