package usecase

import (
	"context"
	"time"

	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/google/uuid"
)

type LinkCacheManager interface {
	Set(context.Context, domain.Link) error
	GetByID(context.Context, uuid.UUID) (domain.Link, error)
	DeleteByID(context.Context, uuid.UUID) error
	Clear(context.Context) error
}

type LinkStorageManager interface {
	Create(context.Context, domain.Link) (domain.Link, error)
	GetByID(context.Context, uuid.UUID) (domain.Link, error)
	UpdateByID(context.Context, domain.Link) (domain.Link, error)
	UpdateExpiresAt(context.Context, uuid.UUID, time.Time) error
	MarkIsDeletedByID(context.Context, uuid.UUID) error
	DeleteManyByIDs(context.Context, []uuid.UUID) error
}

type SnowflakeManager interface {
	CaptureNodeID(context.Context) (int64, error)
}
