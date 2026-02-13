package command

import "context"

type notificationDeleter interface {
	Delete(ctx context.Context, notificationConfigurationID string) error
}
