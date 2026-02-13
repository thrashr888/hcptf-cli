package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type runTaskCreator interface {
	Create(ctx context.Context, organization string, options tfe.RunTaskCreateOptions) (*tfe.RunTask, error)
}

type runTaskUpdater interface {
	Update(ctx context.Context, runTaskID string, options tfe.RunTaskUpdateOptions) (*tfe.RunTask, error)
}
