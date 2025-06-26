package domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

const (
	// defaultOffsetTTL - кол-во времени, сколько должна жить ссылка,
	// если явно не указывали при создании
	DefaultOffsetTTL = time.Duration((time.Hour * 24) * 365)
)

type Link struct {
	ID          uuid.UUID
	FullURI     string
	ShortURI    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ExpiresAt   time.Time
}

func (l Link) IsEmpty() bool {
	return l.ID == uuid.Nil
}

type CreateLinkRequest struct {
	FullURI   string
	ExpiresAt mo.Option[time.Time]
}

type GetLinkRequest struct {
	ID uuid.UUID
}

type PatchLinkRequest struct {
	ID        uuid.UUID
	ExpiresAt time.Time
}

type PutLinkRequest struct {
	ID        uuid.UUID
	FullURI   string
	ExpiresAt mo.Option[time.Time]
}

type DeleteLinkRequest struct {
	ID uuid.UUID
}

type DeleteLinkResponse struct {
	IsSuccessfully bool
}
