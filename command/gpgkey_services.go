package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type gpgKeyCreator interface {
	Create(ctx context.Context, registryName tfe.RegistryName, options tfe.GPGKeyCreateOptions) (*tfe.GPGKey, error)
}

type gpgKeyLister interface {
	ListPrivate(ctx context.Context, options tfe.GPGKeyListOptions) (*tfe.GPGKeyList, error)
}

type gpgKeyReader interface {
	Read(ctx context.Context, keyID tfe.GPGKeyID) (*tfe.GPGKey, error)
}

type gpgKeyUpdater interface {
	Update(ctx context.Context, keyID tfe.GPGKeyID, options tfe.GPGKeyUpdateOptions) (*tfe.GPGKey, error)
}

type gpgKeyDeleter interface {
	Delete(ctx context.Context, keyID tfe.GPGKeyID) error
}
