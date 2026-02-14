package command

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

type notificationDeleter interface {
	Delete(ctx context.Context, notificationConfigurationID string) error
}

type notificationReadDeleter interface {
	notificationDeleter
	Read(ctx context.Context, notificationConfigurationID string) (*tfe.NotificationConfiguration, error)
}
