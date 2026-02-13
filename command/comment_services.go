package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type commentLister interface {
	List(ctx context.Context, runID string) (*tfe.CommentList, error)
}

type commentCreator interface {
	Create(ctx context.Context, runID string, options tfe.CommentCreateOptions) (*tfe.Comment, error)
}

type commentReader interface {
	Read(ctx context.Context, commentID string) (*tfe.Comment, error)
}
