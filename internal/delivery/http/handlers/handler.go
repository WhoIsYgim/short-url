package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
	"net/http"
	"short-link/internal/delivery"
	"short-link/internal/delivery/http/dto"
	"short-link/pkg/errs"
)

type ShortLinkHandler struct {
	usecase delivery.LinkUsecase
}

func NewShortLinkHandler(lu delivery.LinkUsecase) *ShortLinkHandler {
	return &ShortLinkHandler{
		usecase: lu,
	}
}

func (sh *ShortLinkHandler) GetLink(ctx *gin.Context) {
	token := ctx.Param("key")

	if token == "" {
		_ = ctx.Error(errs.BadRequestError())
		return
	}

	origLink, err := sh.usecase.GetOriginalLink(token)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Redirect(http.StatusFound, origLink)
}

func (sh *ShortLinkHandler) CreateLink(ctx *gin.Context) {
	input := new(dto.CreateLinkRequest)
	if err := easyjson.UnmarshalFromReader(ctx.Request.Body, input); err != nil {
		_ = ctx.Error(errs.BadRequestError())
		return
	}

	if input.Link == "" {
		_ = ctx.Error(errs.BadRequestError())
		return
	}
	link, err := sh.usecase.CreateShortLink(input)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	response := &dto.CreateLinkResponse{
		ShortLink: link.ShortLink,
		ExpiresAt: link.ExpiresAt,
	}

	responseJSON, err := response.MarshalJSON()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Data(http.StatusOK, "application/json; charset=utf-8", responseJSON)
}
