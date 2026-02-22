package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type runLister interface {
	List(ctx context.Context, workspaceID string, options *tfe.RunListOptions) (*tfe.RunList, error)
}

type runCreator interface {
	Create(ctx context.Context, options tfe.RunCreateOptions) (*tfe.Run, error)
}

type runApplier interface {
	Apply(ctx context.Context, runID string, options tfe.RunApplyOptions) error
}

type runCanceler interface {
	Cancel(ctx context.Context, runID string, options tfe.RunCancelOptions) error
	ForceCancel(ctx context.Context, runID string, options tfe.RunForceCancelOptions) error
}

type runDiscarder interface {
	Discard(ctx context.Context, runID string, options tfe.RunDiscardOptions) error
}

type runReader interface {
	Read(ctx context.Context, runID string) (*tfe.Run, error)
}

type runForceExecutor interface {
	ForceExecute(ctx context.Context, runID string) error
}
