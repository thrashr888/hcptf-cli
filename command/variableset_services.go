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
