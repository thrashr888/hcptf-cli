package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type accountReader interface {
	ReadCurrent(ctx context.Context) (*tfe.User, error)
}
