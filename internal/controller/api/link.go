package api

import (
	"context"

	apiv1 "github.com/LittleLe6owski/link-shortener/api/v1"
)

func (s *LinkShortenerService) CreateLink(
	ctx context.Context, request *apiv1.CreateLinkRequest,
) (*apiv1.Link, error) {
	entity, err := fromProtoCreateLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	m, err := s.LinkShortenerService.CreateLink(ctx, entity)
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

	link, err := s.LinkShortenerService.GetLink(ctx, entity)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return toProtoCreateLinkResponse(link)
}

func (s *LinkShortenerService) PatchLink(
	ctx context.Context, request *apiv1.PatchLinkRequest,
) (*apiv1.PatchLinkResponse, error) {
	domain, err := fromProtoPatchLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	isUpdated, err := s.LinkShortenerService.UpdateExpiresAt(ctx, domain)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &apiv1.PatchLinkResponse{IsSuccessfully: isUpdated}, nil
}

func (s *LinkShortenerService) PutLink(
	ctx context.Context, request *apiv1.PutLinkRequest,
) (*apiv1.PutLinkResponse, error) {
	domain, err := fromProtoPutLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	isUpdated, err := s.LinkShortenerService.UpdateLink(ctx, domain)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &apiv1.PutLinkResponse{IsSuccessfully: isUpdated}, nil
}

func (s *LinkShortenerService) DeleteLink(
	ctx context.Context, request *apiv1.DeleteLinkRequest,
) (*apiv1.DeleteLinkResponse, error) {
	deleteReq, err := fromProtoDeleteLinkRequest(request)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	isDeleted, err := s.LinkShortenerService.DeleteLink(ctx, deleteReq)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &apiv1.DeleteLinkResponse{IsSuccessfully: isDeleted}, nil
}

func (s *LinkShortenerService) Redirect(
	ctx context.Context, request *apiv1.RedirectRequest,
) (*apiv1.RedirectResponse, error) {
	fullURI, err := s.LinkShortenerService.GetByRedirect(ctx, request.ShortLink)
	if err != nil {
		return nil, toGRPCErr(err)
	}

	return &apiv1.RedirectResponse{FullUri: fullURI}, nil
}
