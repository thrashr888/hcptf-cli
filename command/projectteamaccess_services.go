package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type projectTeamAccessLister interface {
	List(ctx context.Context, options tfe.TeamProjectAccessListOptions) (*tfe.TeamProjectAccessList, error)
}

type projectTeamAccessReader interface {
	Read(ctx context.Context, teamProjectAccessID string) (*tfe.TeamProjectAccess, error)
}

type projectTeamAccessCreator interface {
	Add(ctx context.Context, options tfe.TeamProjectAccessAddOptions) (*tfe.TeamProjectAccess, error)
}

type projectTeamAccessUpdater interface {
	Update(ctx context.Context, teamProjectAccessID string, options tfe.TeamProjectAccessUpdateOptions) (*tfe.TeamProjectAccess, error)
}

type projectTeamAccessDeleter interface {
	Remove(ctx context.Context, teamProjectAccessID string) error
}
