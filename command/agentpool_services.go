package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type agentPoolCreator interface {
	Create(ctx context.Context, organization string, options tfe.AgentPoolCreateOptions) (*tfe.AgentPool, error)
}
