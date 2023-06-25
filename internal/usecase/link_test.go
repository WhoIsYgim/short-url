package usecase

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"short-link/config"
	"short-link/internal/delivery/http/dto"
	"short-link/internal/entities"
	mock_usecase "short-link/internal/usecase/mocks"
	"short-link/pkg/errs"
	"testing"
)

var (
	cfg = &config.Config{
		ServiceConfig: config.ServiceConfig{
			Host:           "localhost",
			Port:           8080,
			RecalcInterval: 10,
		},
		LinkConfig: config.LinkConfig{
			RecreateRetries: 2,
		},
	}
	prefix = "http://localhost:8080/url/"
)

func TestLinkService_GetOriginalLink(t *testing.T) {

	tests := []struct {
		name          string
		expectedLink  string
		expectedError error
		token         string

		mockBehaviour func(repository *mock_usecase.MockLinkRepository,
			cache *mock_usecase.MockTokenCache, generator *mock_usecase.MockGenerator,
			token, link string)
	}{
		{
			name:         "Success",
			expectedLink: "http://wikipedia.org",
			token:        "qwerty123_",
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository,
				cache *mock_usecase.MockTokenCache, generator *mock_usecase.MockGenerator,
				token, link string) {
				linkReturned := &entities.Link{
					OriginalLink: link,
				}
				repository.EXPECT().GetLink(token).Return(linkReturned, nil)
			},
		},
		{
			name:          "Not found",
			expectedLink:  "",
			expectedError: errs.NotFoundError(),
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository,
				cache *mock_usecase.MockTokenCache, generator *mock_usecase.MockGenerator,
				token, link string) {
				repository.EXPECT().GetLink(token).Return(nil, errs.NotFoundError())
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_usecase.NewMockLinkRepository(ctrl)
			mockCache := mock_usecase.NewMockTokenCache(ctrl)
			mockGenerator := mock_usecase.NewMockGenerator(ctrl)

			test.mockBehaviour(mockRepo, mockCache, mockGenerator, test.token, test.expectedLink)

			usecase := LinkService{
				repo:            mockRepo,
				tokenCache:      mockCache,
				generator:       mockGenerator,
				cfg:             cfg,
				shortlinkPrefix: prefix,
			}

			link, err := usecase.GetOriginalLink(test.token)
			if test.expectedError != nil {
				require.ErrorAs(t, err, &test.expectedError)
				return
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, test.expectedLink, link)
		})
	}
}

func TestLinkService_CreateShortLink(t *testing.T) {
	tests := []struct {
		name          string
		expectedLink  *entities.Link
		expectedError error
		token         string
		dto           *dto.CreateLinkRequest

		mockBehaviour func(repository *mock_usecase.MockLinkRepository,
			cache *mock_usecase.MockTokenCache, generator *mock_usecase.MockGenerator,
			dto *dto.CreateLinkRequest, link *entities.Link)
	}{
		{
			name: "Success",
			dto: &dto.CreateLinkRequest{
				Link: "http://wikipedia.org",
			},
			expectedLink: &entities.Link{
				OriginalLink: "http://wikipedia.org",
				Token:        "qwerty123_",
				ExpiresAt:    gomock.Any().String(),
				ShortLink:    prefix + "qwerty123_",
			},
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository, cache *mock_usecase.MockTokenCache, generator *mock_usecase.MockGenerator, dto *dto.CreateLinkRequest, link *entities.Link) {
				repository.EXPECT().GetLinkByOriginal(dto.Link).Return(nil, nil)
				generator.EXPECT().GenString().Return(link.Token).AnyTimes()
				cache.EXPECT().Exists(link.Token).Return(false).AnyTimes()
				repository.EXPECT().StoreLink(gomock.Any()).Return(nil)

				cache.EXPECT().Store(link.Token)
			},
		}, {
			name: "unable to create token",
			dto: &dto.CreateLinkRequest{
				Link: "http://wikipedia.org",
			},
			expectedLink:  nil,
			expectedError: errs.NewAppError(errs.UnableToCreateLink, nil),
			mockBehaviour: func(repository *mock_usecase.MockLinkRepository, cache *mock_usecase.MockTokenCache, generator *mock_usecase.MockGenerator, dto *dto.CreateLinkRequest, link *entities.Link) {
				repository.EXPECT().GetLinkByOriginal(dto.Link).Return(nil, nil)
				generator.EXPECT().GenString().Return("qwerty123_").AnyTimes()
				cache.EXPECT().Exists("qwerty123_").Return(true).AnyTimes()
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_usecase.NewMockLinkRepository(ctrl)
			mockCache := mock_usecase.NewMockTokenCache(ctrl)
			mockGenerator := mock_usecase.NewMockGenerator(ctrl)

			test.mockBehaviour(mockRepo, mockCache, mockGenerator, test.dto, test.expectedLink)

			usecase := LinkService{
				repo:            mockRepo,
				tokenCache:      mockCache,
				generator:       mockGenerator,
				cfg:             cfg,
				shortlinkPrefix: prefix,
			}

			link, err := usecase.CreateShortLink(test.dto)
			if test.expectedError != nil {
				require.ErrorAs(t, err, &test.expectedError)
				return
			} else {
				require.NoError(t, err)
			}
			if link != nil {
				link.ExpiresAt = ""
			}
			require.Equal(t, test.expectedLink.ShortLink, link.ShortLink)
			require.Equal(t, test.expectedLink.ShortLink, link.ShortLink)
		})
	}
}
