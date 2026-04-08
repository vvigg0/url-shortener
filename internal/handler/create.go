package handler

import (
	"net/http"
	"net/url"

	"github.com/vvigg0/l3/url-shortener/internal/dto"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) CreateShortLink(ctx *ginext.Context) {
	var req dto.PostRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		zlog.Logger.Error().Msgf("ошибка парсинга json: %v", err)
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": "invalid JSON"})
		return
	}
	if req.URL == "" {
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": "invalid URL"})
		return
	}
	_, err := url.Parse(req.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ginext.H{"err": "invalid URL"})
		return
	}
	shortLink, err := h.srvc.CreateShortLink(req.URL, req.CustomCode)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка создания сокращённой ссылки: %v", err)
		statusCode := getErrorCode(err)
		ctx.JSON(statusCode, ginext.H{"err": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, ginext.H{"res": shortLink})
}
