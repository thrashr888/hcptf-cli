package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type projectLister interface {
	List(ctx context.Context, organization string, options *tfe.ProjectListOptions) (*tfe.ProjectList, error)
}

type projectReader interface {
	Read(ctx context.Context, projectID string) (*tfe.Project, error)
}
