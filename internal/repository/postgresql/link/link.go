package pglink

import (
	"context"
	"time"

	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/LittleLe6owski/link-shortener/internal/repository/postgresql"
	"github.com/LittleLe6owski/link-shortener/internal/usecase"
	"github.com/google/uuid"
)

var (
	_ usecase.LinkStorageManager = (*LinkRepository)(nil)
)

type LinkRepository struct {
	connector postgresql.Connect
}

func NewRepository(connector postgresql.Connect) LinkRepository {
	return LinkRepository{connector: connector}
}

func (l LinkRepository) Create(ctx context.Context, link domain.Link) (domain.Link, error) {
	var result Link

	insertQu := `INSERT INTO link_shortener.links(
		id, full_uri, short_uri, created_at, updated_at, expires_at) VALUES  ($1, $2, $3, $4, $5, $6)
			RETURNING id, full_uri, short_uri, created_at, updated_at, expires_at`

	err := l.connector.Pgx.QueryRow(
		ctx, insertQu, link.ID, link.FullURI, link.ShortURI, link.CreatedAt, link.UpdatedAt, link.ExpiresAt,
	).Scan(&result.ID, &result.FullURI, &result.ShortURI, &result.CreatedAt, &result.UpdatedAt, &result.ExpiresAt)
	if err != nil {
		return domain.Link{}, ToEntityError(err)
	}

	return toDomainMyItem(result), nil
}

func (l LinkRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Link, error) {
	var result Link

	selectQuery := `SELECT id, full_uri, short_uri, expires_at, created_at, updated_at
						FROM
						    link_shortener.links
						WHERE
							id = $1 AND is_deleted = false`

	if err := l.connector.ScanAPI.Get(ctx, l.connector.Pgx, &result, selectQuery, id); err != nil {
		return domain.Link{}, ToEntityError(err)
	}

	return toDomainMyItem(result), nil
}

func (l LinkRepository) GetByShortURI(ctx context.Context, shortURI int64) (domain.Link, error) {
	var result Link

	selectQuery := `SELECT id, full_uri, short_uri, expires_at, created_at, updated_at
						FROM
						    link_shortener.links
						WHERE
							short_uri = $1 AND is_deleted = false`

	if err := l.connector.ScanAPI.Get(ctx, l.connector.Pgx, &result, selectQuery, shortURI); err != nil {
		return domain.Link{}, ToEntityError(err)
	}

	return toDomainMyItem(result), nil
}

func (l LinkRepository) UpdateByID(ctx context.Context, link domain.Link) (domain.Link, error) {
	var result Link

	updateQuery := `UPDATE link_shortener.links SET 
		full_uri = $1, short_uri = $2, expires_at = $3, updated_at = $4, WHERE id = $3
			RETURNING id, full_uri, short_uri, expires_at, created_at, updated_at`

	err := l.connector.ScanAPI.Select(ctx, l.connector.Pgx, result, updateQuery,
		link.FullURI, link.ShortURI, link.ExpiresAt, link.UpdatedAt, link.ID,
	)
	if err != nil {
		return domain.Link{}, ToEntityError(err)
	}

	return toDomainMyItem(result), nil
}

func (l LinkRepository) UpdateExpiresAt(ctx context.Context, id uuid.UUID, newExp time.Time) error {
	updateQuery := `UPDATE link_shortener.links SET expires_at = $1, updated_at = $2, WHERE id = $3`

	_, err := l.connector.Pgx.Exec(ctx, updateQuery, newExp, time.Now().UTC(), id)
	if err != nil {
		return ToEntityError(err)
	}

	return nil
}

func (l LinkRepository) MarkIsDeletedByID(ctx context.Context, id uuid.UUID) error {
	updateQuery := `UPDATE link_shortener.links SET is_deleted = true, updated_at = $1 WHERE id = $2`

	if _, err := l.connector.Pgx.Exec(ctx, updateQuery, time.Now().UTC(), id.String()); err != nil {
		return ToEntityError(err)
	}

	return nil
}

func (l LinkRepository) DeleteManyByIDs(ctx context.Context, id []uuid.UUID) error {
	return nil
}
