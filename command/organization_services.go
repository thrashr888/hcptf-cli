package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type organizationLister interface {
	List(ctx context.Context, options *tfe.OrganizationListOptions) (*tfe.OrganizationList, error)
}

type organizationReader interface {
	Read(ctx context.Context, organization string) (*tfe.Organization, error)
}
