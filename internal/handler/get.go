package handler

import (
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) ShortLinkRedirect(ctx *ginext.Context) {
	urlID, fullURL, err := h.srvc.GetFullURL(ctx.Param("short_code"))
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка получения полного url: %v", err)
		statusCode := getErrorCode(err)
		ctx.JSON(statusCode, ginext.H{"err": err.Error()})
		return
	}
	userAgent := ctx.Request.UserAgent()
	nowUTC := time.Now().UTC()
	if err := h.srvc.WriteAnalytics(urlID, userAgent, nowUTC); err != nil {
		zlog.Logger.Error().Msgf("ошибка записи аналитики: %v", err)
	}
	ctx.Redirect(http.StatusFound, fullURL)
}

func (h *Handler) GetAnalytics(ctx *ginext.Context) {
	aggregateBy := ctx.Query("aggregate_by")

	analytics, err := h.srvc.GetAnalytics(ctx.Param("short_code"), aggregateBy)
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка получения аналитики: %v", err)
		ctx.JSON(getErrorCode(err), ginext.H{"err": err})
		return
	}
	ctx.JSON(http.StatusOK, ginext.H{"res": analytics})
}

func (h *Handler) GetAllLinks(ctx *ginext.Context) {
	links, err := h.srvc.GetAllLinks()
	if err != nil {
		zlog.Logger.Error().Msgf("ошибка запроса всех ссылок: %v", err)
		statusCode := getErrorCode(err)
		ctx.JSON(statusCode, ginext.H{"err": err})
		return
	}
	ctx.JSON(http.StatusOK, ginext.H{"res": links})
}
