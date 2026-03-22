package handler

import (
	"net/http"

	reportService "idc-finance/internal/modules/report/service"
	appErrors "idc-finance/internal/platform/errors"
	"idc-finance/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *reportService.Service
}

func New(service *reportService.Service) *Handler {
	return &Handler{service: service}
}

func (handler *Handler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/workbench", handler.getWorkbench)
	router.GET("/reports/overview", handler.getReportOverview)
}

func (handler *Handler) getWorkbench(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.GetWorkbench(), middleware.GetRequestID(c)))
}

func (handler *Handler) getReportOverview(c *gin.Context) {
	c.JSON(http.StatusOK, appErrors.Ok(handler.service.GetReportOverview(), middleware.GetRequestID(c)))
}
