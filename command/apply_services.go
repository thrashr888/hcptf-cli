package command

import (
	"context"
	"io"

	tfe "github.com/hashicorp/go-tfe"
)

type applyReader interface {
	Read(ctx context.Context, applyID string) (*tfe.Apply, error)
}

type applyLogReader interface {
	Logs(ctx context.Context, applyID string) (io.Reader, error)
}
