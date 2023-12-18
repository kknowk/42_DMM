package repository

import (
	"context"

	"yatter-backend-go/app/domain/object"
)

type Status interface {

	FindByID(ctx context.Context, id object.StatusID) (*object.Status, error)

	AddStatus(ctx context.Context, status *object.Status) error

	DeleteStatus(ctx context.Context, status *object.Status) error

	FindAllPublicStatuses(ctx context.Context, onlyMedia bool, maxID int64, sinceID int64, limit int) ([]*object.Status, error)
}