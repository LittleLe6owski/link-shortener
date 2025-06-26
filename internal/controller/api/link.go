package api

import (
	"context"

	apiv1 "github.com/LittleLe6owski/link-shortener/api/v1"
)

func (s *LinkShortenerService) CreateLink(
	ctx context.Context, request *apiv1.CreateLinkRequest,
) (*apiv1.Link, error) {
	m, err := s.LinkShortenerService.CreateLink(ctx, fromProtoCreateLinkRequest(request))
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return toProtoCreateLinkResponse(m)
}

func (s *LinkShortenerService) GetLink(
	ctx context.Context, request *apiv1.GetLinkRequest,
) (*apiv1.Link, error) {
	entity, err := fromProtoGetLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	m, err := s.LinkShortenerService.GetLink(ctx, entity)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return toProtoCreateLinkResponse(m)
}

func (s *LinkShortenerService) PatchLink(
	ctx context.Context, request *apiv1.PatchLinkRequest,
) (*apiv1.PatchLinkResponse, error) {
	domain, err := fromProtoPatchLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	m, err := s.LinkShortenerService.UpdateExpiresAt(ctx, domain)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &apiv1.PatchLinkResponse{IsSuccessfully: m}, nil
}

func (s *LinkShortenerService) PutLink(
	ctx context.Context, request *apiv1.PutLinkRequest,
) (*apiv1.PutLinkResponse, error) {
	domain, err := fromProtoPutLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	m, err := s.LinkShortenerService.UpdateLink(ctx, domain)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &apiv1.PutLinkResponse{IsSuccessfully: m}, nil
}

func (s *LinkShortenerService) DeleteLink(
	ctx context.Context, request *apiv1.CreateLinkRequest,
) (*apiv1.Link, error) {
	m, err := s.LinkShortenerService.CreateLink(ctx, fromProtoCreateLinkRequest(request))
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return toProtoCreateLinkResponse(m)
}
