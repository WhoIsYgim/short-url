package delivery

import (
	"short-link/internal/delivery/http/dto"
	"short-link/internal/entities"
)

type LinkUsecase interface {
	GetOriginalLink(token string) (string, error)
	CreateShortLink(request *dto.CreateLinkRequest) (*entities.Link, error)
}
