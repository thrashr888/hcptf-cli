package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type policySetLister interface {
	List(ctx context.Context, organization string, options *tfe.PolicySetListOptions) (*tfe.PolicySetList, error)
}

type policySetReader interface {
	Read(ctx context.Context, policySetID string) (*tfe.PolicySet, error)
}
