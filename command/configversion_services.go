package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type configVersionLister interface {
	List(ctx context.Context, workspaceID string, options *tfe.ConfigurationVersionListOptions) (*tfe.ConfigurationVersionList, error)
}

type configVersionReader interface {
	Read(ctx context.Context, configurationID string) (*tfe.ConfigurationVersion, error)
	ReadWithOptions(ctx context.Context, configurationID string, options *tfe.ConfigurationVersionReadOptions) (*tfe.ConfigurationVersion, error)
}

type configVersionCreator interface {
	Create(ctx context.Context, workspaceID string, options tfe.ConfigurationVersionCreateOptions) (*tfe.ConfigurationVersion, error)
}

type configVersionUploader interface {
	Upload(ctx context.Context, url, path string) error
}
