package command

import (
	"context"
	"io"

	tfe "github.com/hashicorp/go-tfe"
)

type planReader interface {
	Read(ctx context.Context, planID string) (*tfe.Plan, error)
}

type planLogReader interface {
	Logs(ctx context.Context, planID string) (io.Reader, error)
}
