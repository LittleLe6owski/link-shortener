package pglink

import (
	"errors"
	"fmt"

	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx"
)

type Link struct {
	ID        pgtype.UUID        `db:"id"`
	CreatedAt pgtype.Timestamptz `db:"created_at"`
	UpdatedAt pgtype.Timestamptz `db:"updated_at"`
	ExpiresAt pgtype.Timestamptz `db:"expires_at"`
	FullURI   string             `db:"full_uri"`
	ShortURI  string             `db:"short_uri"`
}

func fromDomainMyItem(in domain.Link) (Link, error) {
	link := Link{
		ID:       pgtype.UUID{Bytes: in.ID, Status: pgtype.Present},
		FullURI:  in.FullURI,
		ShortURI: in.ShortURI,
	}

	if err := link.CreatedAt.Scan(in.CreatedAt); err != nil {
		return Link{}, err
	}

	if err := link.UpdatedAt.Scan(in.UpdatedAt); err != nil {
		return Link{}, err
	}

	if err := link.ExpiresAt.Scan(in.ExpiresAt); err != nil {
		return Link{}, err
	}

	return link, nil
}

func toDomainMyItem(in Link) domain.Link {
	id, err := uuid.FromBytes(in.ID.Bytes[:])
	if err != nil {
		id = uuid.Nil
	}

	return domain.Link{
		ID:        id,
		ShortURI:  in.ShortURI,
		FullURI:   in.FullURI,
		CreatedAt: in.CreatedAt.Time,
		UpdatedAt: in.UpdatedAt.Time,
		ExpiresAt: in.ExpiresAt.Time,
	}
}

func ToEntityError(repoErr error) error {
	if repoErr == nil {
		return nil
	}

	if errors.Is(repoErr, pgx.ErrNoRows) {
		return fmt.Errorf("%w:%s", domain.ErrNotFound, repoErr.Error())
	}

	return fmt.Errorf("%w:%s", domain.ErrUnexpected, repoErr.Error())
}
