package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type registryModuleLister interface {
	List(ctx context.Context, organization string, options *tfe.RegistryModuleListOptions) (*tfe.RegistryModuleList, error)
}

type registryModuleCreator interface {
	Create(ctx context.Context, organization string, options tfe.RegistryModuleCreateOptions) (*tfe.RegistryModule, error)
}

type registryModuleReader interface {
	Read(ctx context.Context, moduleID tfe.RegistryModuleID) (*tfe.RegistryModule, error)
}

type registryModuleDeleter interface {
	Delete(ctx context.Context, organization, name string) error
}

type registryModuleVersionDeleter interface {
	DeleteVersion(ctx context.Context, moduleID tfe.RegistryModuleID, version string) error
}

type registryModuleVersionCreator interface {
	CreateVersion(ctx context.Context, moduleID tfe.RegistryModuleID, options tfe.RegistryModuleCreateVersionOptions) (*tfe.RegistryModuleVersion, error)
}
