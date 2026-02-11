package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type planExportCreator interface {
	Create(ctx context.Context, options tfe.PlanExportCreateOptions) (*tfe.PlanExport, error)
}

type planExportReader interface {
	Read(ctx context.Context, planExportID string) (*tfe.PlanExport, error)
}

type planExportDownloader interface {
	Download(ctx context.Context, planExportID string) ([]byte, error)
}
