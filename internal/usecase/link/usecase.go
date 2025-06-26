package link

import (
	"context"
	"fmt"
	"time"

	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/LittleLe6owski/link-shortener/internal/usecase"
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type LinkManager struct {
	linkRepoPG       usecase.LinkStorageManager
	linkCache        usecase.LinkCacheManager
	snowflakeManager usecase.SnowflakeManager
	logger           zerolog.Logger
}

func NewLinkManger(
	linkRepoPG usecase.LinkStorageManager,
	linkCache usecase.LinkCacheManager,
	snowflakeManager usecase.SnowflakeManager,
	logger zerolog.Logger,
) LinkManager {
	return LinkManager{
		linkRepoPG: linkRepoPG, linkCache: linkCache, logger: logger, snowflakeManager: snowflakeManager,
	}
}

func (s *LinkManager) CreateLink(
	ctx context.Context, request domain.CreateLinkRequest,
) (domain.Link, error) {
	tn := time.Now()

	expAt := tn.Add(domain.DefaultOffsetTTL)
	if request.ExpiresAt.IsPresent() {
		expAt = request.ExpiresAt.MustGet()
	}

	nodeID, err := s.snowflakeManager.CaptureNodeID(ctx)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to capture node: %w", err)
	}

	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get new node by id: %d, %w", nodeID, err)
	}

	createLink := domain.Link{
		ID:        uuid.New(),
		FullURI:   request.FullURI,
		ShortURI:  node.Generate().String(),
		CreatedAt: tn,
		UpdatedAt: tn,
		ExpiresAt: expAt,
	} //exhaustruct:enforce

	createdLink, err := s.linkRepoPG.Create(ctx, createLink)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to create link, %w", err)
	}

	return createdLink, nil
}

func (s *LinkManager) GetLink(ctx context.Context, request domain.GetLinkRequest) (domain.Link, error) {
	link, err := s.linkCache.GetByID(ctx, request.ID)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get link in cache: %w", err)
	}

	if !link.IsEmpty() {
		return link, nil
	}

	link, err = s.linkRepoPG.GetByID(ctx, request.ID)
	if err != nil {
		return domain.Link{}, fmt.Errorf("failed to get link in storage: %w", err)
	}

	if err = s.linkCache.Set(ctx, link); err != nil {
		return domain.Link{}, fmt.Errorf("failed to set link in cache: %w", err)
	}

	return link, nil
}

func (s *LinkManager) UpdateExpiresAt(ctx context.Context, request domain.PatchLinkRequest) (bool, error) {
	if err := s.linkRepoPG.UpdateExpiresAt(ctx, request.ID, request.ExpiresAt); err != nil {
		return false, fmt.Errorf("failed to update exp at from link, %w", err)
	}

	return true, nil
}

func (s *LinkManager) UpdateLink(ctx context.Context, request domain.PutLinkRequest) (bool, error) {
	tn := time.Now()

	expAt := tn.Add(domain.DefaultOffsetTTL)
	if request.ExpiresAt.IsPresent() {
		expAt = request.ExpiresAt.MustGet()
	}

	nodeID, err := s.snowflakeManager.CaptureNodeID(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to capture node: %w", err)
	}

	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return false, fmt.Errorf("failed to get new node by id: %d, %w", nodeID, err)
	}

	newLink := domain.Link{
		ID:        request.ID,
		FullURI:   request.FullURI,
		ShortURI:  node.Generate().String(),
		UpdatedAt: time.Now().UTC(),
		ExpiresAt: expAt,
	} //exhaustruct:enforce

	if err := s.linkCache.Set(ctx, newLink); err != nil {
		return false, fmt.Errorf("failed to update link in cache, %w", err)
	}

	if _, err = s.linkRepoPG.UpdateByID(ctx, newLink); err != nil {
		if deleteErr := s.linkCache.DeleteByID(ctx, request.ID); deleteErr != nil {
			return false, fmt.Errorf(
				"%w failed to delete link in cache, date is inconsistent: %w", domain.ErrUnexpected, deleteErr,
			)
		}

		return false, fmt.Errorf("failed to update link in storage, %w", err)
	}

	return true, nil
}

func (s *LinkManager) DeleteLink(ctx context.Context, request domain.DeleteLinkRequest) (bool, error) {
	if err := s.linkRepoPG.MarkIsDeletedByID(ctx, request.ID); err != nil {
		return false, fmt.Errorf("failed to delete link in storage: %w", err)
	}

	if err := s.linkCache.DeleteByID(ctx, request.ID); err != nil {
		return false, fmt.Errorf("failed to delete link in cache: %w", err)
	}

	return true, nil
}
