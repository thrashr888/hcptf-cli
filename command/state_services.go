package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type stateVersionReader interface {
	ReadCurrent(ctx context.Context, workspaceID string) (*tfe.StateVersion, error)
}
