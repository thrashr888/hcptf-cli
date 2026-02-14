package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type agentPoolCreator interface {
	Create(ctx context.Context, organization string, options tfe.AgentPoolCreateOptions) (*tfe.AgentPool, error)
}

type agentPoolDeleter interface {
	Delete(ctx context.Context, agentPoolID string) error
}
