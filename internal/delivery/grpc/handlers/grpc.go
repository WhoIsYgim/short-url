package handlers

import (
	"context"
	"short-link/internal/delivery"
	"short-link/internal/delivery/http/dto"
	"short-link/pkg/errs"
	"short-link/pkg/grpc/api"
)

type LinkHandlerGrpc struct {
	usecase delivery.LinkUsecase
	api.UnimplementedShortLinkServiceServer
}

func NewLinkHandler(usecase delivery.LinkUsecase) *LinkHandlerGrpc {
	return &LinkHandlerGrpc{
		usecase: usecase,
	}
}

func (lh *LinkHandlerGrpc) GetOriginalLink(ctx context.Context, req *api.ShortLinkRequest) (*api.ShortLinkResponse, error) {
	if req.ShortLink == "" {
		return nil, errs.BadRequestError()
	}
	link, err := lh.usecase.GetOriginalLink(req.ShortLink)
	if err != nil {
		return nil, err
	}
	return &api.ShortLinkResponse{
		OriginalLink: link,
	}, nil
}
func (lh *LinkHandlerGrpc) CreateShortLink(ctx context.Context, req *api.CreateShortLinkRequest) (*api.CreateShortLinkResponse, error) {
	input := &dto.CreateLinkRequest{
		Link: req.OriginalLink,
	}
	link, err := lh.usecase.CreateShortLink(input)
	if err != nil {
		return nil, err
	}

	return &api.CreateShortLinkResponse{
		ShortLink: link.ShortLink,
		ExpiresAt: link.ExpiresAt,
	}, nil
}
