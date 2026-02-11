package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type registryProviderLister interface {
	List(ctx context.Context, organization string, options *tfe.RegistryProviderListOptions) (*tfe.RegistryProviderList, error)
}

type registryProviderCreator interface {
	Create(ctx context.Context, organization string, options tfe.RegistryProviderCreateOptions) (*tfe.RegistryProvider, error)
}

type registryProviderReader interface {
	Read(ctx context.Context, providerID tfe.RegistryProviderID, options *tfe.RegistryProviderReadOptions) (*tfe.RegistryProvider, error)
}

type registryProviderDeleter interface {
	Delete(ctx context.Context, providerID tfe.RegistryProviderID) error
}
