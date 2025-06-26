package api

import (
	"time"

	apiv1 "github.com/LittleLe6owski/link-shortener/api/v1"
	"github.com/LittleLe6owski/link-shortener/internal/domain"
	"github.com/google/uuid"
	"github.com/samber/mo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func fromProtoCreateLinkRequest(req *apiv1.CreateLinkRequest) domain.CreateLinkRequest {
	expAt := mo.None[time.Time]()
	if expAtOrNil := req.GetExpiresAt(); expAtOrNil != nil {
		expAt = mo.Some(expAtOrNil.AsTime())
	}

	return domain.CreateLinkRequest{FullURI: req.FullUri, ExpiresAt: expAt}
}

func fromProtoGetLinkRequest(req *apiv1.GetLinkRequest) (domain.GetLinkRequest, error) {
	id, err := uuid.Parse(req.GetId())

	return domain.GetLinkRequest{ID: id}, err
}

func fromProtoDeleteLinkRequest(req *apiv1.DeleteLinkRequest) (domain.DeleteLinkRequest, error) {
	id, err := uuid.Parse(req.GetId())

	return domain.DeleteLinkRequest{ID: id}, err
}

func fromProtoPatchLinkRequest(req *apiv1.PatchLinkRequest) (domain.PatchLinkRequest, error) {
	id, err := uuid.Parse(req.GetId())

	return domain.PatchLinkRequest{ID: id, ExpiresAt: req.GetNewExpiresAt().AsTime()}, err
}

func fromProtoPutLinkRequest(req *apiv1.PutLinkRequest) (domain.PutLinkRequest, error) {
	id, err := uuid.Parse(req.GetId())

	return domain.PutLinkRequest{
		ID: id, FullURI: req.FullUri, ExpiresAt: mo.Some(req.GetNewExpiresAt().AsTime()),
	}, err
}

func toProtoCreateLinkResponse(link domain.Link) (*apiv1.Link, error) {
	return &apiv1.Link{
		Id:        link.ID.String(),
		FullUri:   link.FullURI,
		ShortUri:  link.ShortURI,
		CreatedAt: timestamppb.New(link.CreatedAt),
		UpdatedAt: timestamppb.New(link.UpdatedAt),
	}, nil
}
