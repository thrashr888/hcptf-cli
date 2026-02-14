package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type policyLister interface {
	List(ctx context.Context, organization string, options *tfe.PolicyListOptions) (*tfe.PolicyList, error)
}

type policyReader interface {
	Read(ctx context.Context, policyID string) (*tfe.Policy, error)
}

type policyDownloader interface {
	Download(ctx context.Context, policyID string) ([]byte, error)
}
