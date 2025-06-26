package api

import (
	"errors"

	apiv1 "github.com/LittleLe6owski/link-shortener/api/v1"
	"github.com/LittleLe6owski/link-shortener/internal/usecase/link"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LinkShortenerService struct {
	apiv1.UnimplementedLinkShortenerServiceServer
	LinkShortenerService link.LinkManager
}

func NewLinkShortenerServiceServer(service link.LinkManager) *LinkShortenerService {
	return &LinkShortenerService{LinkShortenerService: service}
}

func toGRPCErr(err error) error {
	if errors.Is(err, link.ErrInvalidRequest) {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if errors.Is(err, link.ErrMyItemNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}

	if errors.Is(err, link.ErrMyItemAlreadyExists) {
		return status.Error(codes.AlreadyExists, err.Error())
	}

	return status.Error(codes.Internal, err.Error())
}
