package handlers

import (
	"net/http"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/services"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AirportHandler struct {
	service *services.AirportService
	logger  *zap.Logger
}

func NewAirportHandler(service *services.AirportService, logger *zap.Logger) *AirportHandler {
	return &AirportHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AirportHandler) GetPopularAirports(c *gin.Context) {
	start := time.Now()
	airports, err := h.service.GetPopularAirports()
	if err != nil {
		h.logger.Error("Failed to get popular airports", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "fail", err.Error()))
		return
	}

	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", "Get popular airports successfully", airports, took))
}
