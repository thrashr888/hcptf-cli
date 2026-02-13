package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type teamLister interface {
	List(ctx context.Context, organization string, options *tfe.TeamListOptions) (*tfe.TeamList, error)
}

type teamReader interface {
	Read(ctx context.Context, teamName string) (*tfe.Team, error)
}
