package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type sshKeyCreator interface {
	Create(ctx context.Context, organization string, options tfe.SSHKeyCreateOptions) (*tfe.SSHKey, error)
}

type sshKeyDeleter interface {
	Delete(ctx context.Context, sshKeyID string) error
}
