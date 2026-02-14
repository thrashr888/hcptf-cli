package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type auditTrailTokenDeleter interface {
	DeleteWithOptions(ctx context.Context, organization string, options tfe.OrganizationTokenDeleteOptions) error
}
