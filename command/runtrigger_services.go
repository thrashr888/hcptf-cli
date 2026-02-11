package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type runTriggerLister interface {
	List(ctx context.Context, workspaceID string, options *tfe.RunTriggerListOptions) (*tfe.RunTriggerList, error)
}

type runTriggerCreator interface {
	Create(ctx context.Context, workspaceID string, options tfe.RunTriggerCreateOptions) (*tfe.RunTrigger, error)
}

type runTriggerReader interface {
	Read(ctx context.Context, runTriggerID string) (*tfe.RunTrigger, error)
}

type runTriggerDeleter interface {
	Delete(ctx context.Context, runTriggerID string) error
}
