package usecase

import (
	"fmt"
	"net/url"
	"short-link/config"
	"short-link/internal/delivery"
	"short-link/internal/delivery/http/dto"
	"short-link/internal/entities"
	"short-link/internal/utils"
	"short-link/pkg/errs"
	"time"
)

//go:generate mockgen --destination=mocks/link.go short-link/internal/usecase LinkRepository,TokenCache,Generator

type LinkRepository interface {
	GetLink(token string) (*entities.Link, error)
	GetLinkByOriginal(origLink string) (*entities.Link, error)
	StoreLink(link *entities.Link) error
	StartRecalculation(interval time.Duration, deleted chan []string)
	ShutDown() error
}

type TokenCache interface {
	Exists(token string) bool
	Store(token string)
	SetRecalculationChan(chan []string)
}

type Generator interface {
	GenString() string
}

const DeletedChanBufferSize = 10

type LinkService struct {
	repo            LinkRepository
	tokenCache      TokenCache
	cfg             *config.Config
	generator       Generator
	shortlinkPrefix string
	deleteChan      chan []string
}

func (l *LinkService) GetOriginalLink(token string) (string, error) {
	link, err := l.repo.GetLink(token)
	if err != nil {
		return "", err
	}
	return link.OriginalLink, nil
}

func (l *LinkService) CreateShortLink(request *dto.CreateLinkRequest) (*entities.Link, error) {
	_, err := url.ParseRequestURI(request.Link)
	if err != nil {
		return nil, errs.NewAppError(errs.UrlNotValid, err)
	}
	link, _ := l.repo.GetLinkByOriginal(request.Link)
	if link != nil {
		link.ShortLink = l.shortlinkPrefix + link.Token
		return link, nil
	}

	retries := 0
	exists := true
	token := ""

	for exists && retries < l.cfg.LinkConfig.RecreateRetries {
		token = l.generator.GenString()
		exists = l.tokenCache.Exists(token)
		retries++
	}

	if retries == l.cfg.LinkConfig.RecreateRetries {
		return nil, errs.NewAppError(errs.UnableToCreateLink, nil)
	}

	newLink := &entities.Link{
		OriginalLink: request.Link,
		Token:        token,
		ExpiresAt:    utils.ExpireTimeString(l.cfg.LinkConfig.Expiration),
		ShortLink:    fmt.Sprintf(l.shortlinkPrefix + token),
	}

	err = l.repo.StoreLink(newLink)
	if err != nil {
		return nil, err
	}
	l.tokenCache.Store(newLink.Token)
	return newLink, err
}

func NewLinkService(repo LinkRepository, cfg *config.Config, cache TokenCache, strGenerator Generator) delivery.LinkUsecase {
	deleteChan := make(chan []string)
	cache.SetRecalculationChan(deleteChan)
	repo.StartRecalculation(time.Duration(cfg.ServiceConfig.RecalcInterval)*time.Hour, deleteChan)
	prefix := fmt.Sprintf("http://%s:%d/url/", cfg.ServiceConfig.Host, cfg.ServiceConfig.Port)
	return &LinkService{
		repo:            repo,
		cfg:             cfg,
		generator:       strGenerator,
		tokenCache:      cache,
		shortlinkPrefix: prefix,
		deleteChan:      deleteChan,
	}
}
