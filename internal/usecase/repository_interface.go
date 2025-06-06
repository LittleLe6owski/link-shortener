package usecase

import (
	"context"

	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/google/uuid"
)

type MyItemCacheManager interface {
	Create(ctx context.Context, myItem *domain.Link) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Link, error)
	Update(ctx context.Context, myItem *domain.Link) (*domain.Link, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Clear(ctx context.Context) error
}

type MyItemStorageManager interface {
	Create(ctx context.Context, myItem *domain.Link) error
	Get(ctx context.Context, id uuid.UUID) (*domain.Link, error)
	Update(ctx context.Context, myItem *domain.Link) (*domain.Link, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, pageSize int, pageNumber int) ([]*domain.Link, int, error)
}
