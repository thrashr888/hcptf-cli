package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type azureOIDCCreator interface {
	Create(ctx context.Context, organization string, options tfe.AzureOIDCConfigurationCreateOptions) (*tfe.AzureOIDCConfiguration, error)
}

type azureOIDCReader interface {
	Read(ctx context.Context, configurationID string) (*tfe.AzureOIDCConfiguration, error)
}

type azureOIDCUpdater interface {
	Update(ctx context.Context, configurationID string, options tfe.AzureOIDCConfigurationUpdateOptions) (*tfe.AzureOIDCConfiguration, error)
}

type azureOIDCDeleter interface {
	Read(ctx context.Context, configurationID string) (*tfe.AzureOIDCConfiguration, error)
	Delete(ctx context.Context, configurationID string) error
}
