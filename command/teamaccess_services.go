package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type teamAccessLister interface {
	List(ctx context.Context, options *tfe.TeamAccessListOptions) (*tfe.TeamAccessList, error)
}

type teamAccessReader interface {
	Read(ctx context.Context, teamAccessID string) (*tfe.TeamAccess, error)
}

type teamAccessCreator interface {
	Add(ctx context.Context, options tfe.TeamAccessAddOptions) (*tfe.TeamAccess, error)
}

type teamAccessUpdater interface {
	Update(ctx context.Context, teamAccessID string, options tfe.TeamAccessUpdateOptions) (*tfe.TeamAccess, error)
}

type teamAccessDeleter interface {
	Remove(ctx context.Context, teamAccessID string) error
}
