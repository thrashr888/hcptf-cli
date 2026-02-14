package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type registryProviderVersionCreator interface {
	Create(ctx context.Context, providerID tfe.RegistryProviderID, options tfe.RegistryProviderVersionCreateOptions) (*tfe.RegistryProviderVersion, error)
}

type registryProviderVersionDeleter interface {
	Delete(ctx context.Context, versionID tfe.RegistryProviderVersionID) error
}
