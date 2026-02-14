package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type userReader interface {
	ReadCurrent(ctx context.Context) (*tfe.User, error)
}

type userTokenLister interface {
	List(ctx context.Context, userID string) (*tfe.UserTokenList, error)
}
