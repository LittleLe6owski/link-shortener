package domain

import (
	"time"

	"github.com/google/uuid"
)

type Link struct {
	ID          uuid.UUID
	URI         string
	ShortURI    string
	CreateAt    time.Time
	ExpiresAt   int
	CallCounter int
}
