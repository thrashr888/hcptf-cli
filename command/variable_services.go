package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type variableCreator interface {
	Create(ctx context.Context, workspaceID string, options tfe.VariableCreateOptions) (*tfe.Variable, error)
}

type variableUpdater interface {
	Update(ctx context.Context, workspaceID string, variableID string, options tfe.VariableUpdateOptions) (*tfe.Variable, error)
}

type variableDeleter interface {
	Delete(ctx context.Context, workspaceID string, variableID string) error
}
