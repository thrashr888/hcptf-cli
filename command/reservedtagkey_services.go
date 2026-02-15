package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type reservedTagKeyCreator interface {
	Create(ctx context.Context, organization string, options tfe.ReservedTagKeyCreateOptions) (*tfe.ReservedTagKey, error)
}

type reservedTagKeyUpdater interface {
	Update(ctx context.Context, reservedTagKeyID string, options tfe.ReservedTagKeyUpdateOptions) (*tfe.ReservedTagKey, error)
}
