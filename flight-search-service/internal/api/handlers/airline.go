package handlers

import (
	"fmt"
	"net/http"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/services"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AirlineHandler struct {
	service *services.AirlineService
	logger  *zap.Logger
}

func NewAirlineHandler(service *services.AirlineService, logger *zap.Logger) *AirlineHandler {
	return &AirlineHandler{service: service, logger: logger}
}

func (h *AirlineHandler) GetAllAirlines(c *gin.Context) {
	start := time.Now()
	airlines, err := h.service.GetAllAirlines()
	if err != nil {
		h.logger.Error("Failed to retrieve all airlines", zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Successfully retrieved all airlines", zap.Int("count", len(airlines)), zap.Duration("duration", time.Since(start)))
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved airlines, count: %v", len(airlines)), airlines, time.Since(start).Milliseconds()))
}

func (h *AirlineHandler) FindByAirlineId(c *gin.Context) {
	start := time.Now()
	airlineId := c.Param("airlineId")
	airline, err := h.service.FindByAirlineId(airlineId)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to find airline with id=%s", airlineId, zap.Error(err))
		c.JSON(http.StatusNotFound, utils.NewErrorHandler(http.StatusNotFound, "Not found", err.Error()))
		return
	}

	h.logger.Info("Successfully found airline", zap.String("airlineId", airlineId), zap.Duration("duration", time.Since(start)))
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Found airline with id=%s", airlineId), airline, time.Since(start).Milliseconds()))
}

func (h *AirlineHandler) DeleteByAirlineId(c *gin.Context) {
	start := time.Now()
	airlineId := c.Param("airlineId")
	err := h.service.DeleteByAirlineId(airlineId)
	if err != nil {
		h.logger.Error("Failed to delete airline", zap.String("airlineId", airlineId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Successfully deleted airline", zap.String("airlineId", airlineId), zap.Duration("duration", time.Since(start)))
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", "Deleted airline successfully", "Deleted", time.Since(start).Milliseconds()))
}
