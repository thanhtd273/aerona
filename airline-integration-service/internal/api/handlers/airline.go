package handlers

import (
	"fmt"
	"net/http"
	"time"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"
	"aerona.thanhtd.com/airline-integration-service/internal/api/services"

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

func (h *AirlineHandler) CreateAirline(c *gin.Context) {
	start := time.Now()

	var rawData models.RawAirline
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", "Invalid request payload"))
		return
	}
	airline, err := h.service.CreateAirline(rawData)
	if err != nil {
		h.logger.Error("Failed to create airline", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "fail", fmt.Sprintf("Failed to create airline, error: %v", err.Error())))
		return
	}

	h.logger.Info("Airline created successfully", zap.String("airline_id", rawData.IATACode))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusCreated, models.NewApiResponse(http.StatusCreated, "success", "Airline created successfully", airline, took))
}

func (h *AirlineHandler) GetAllAirlines(c *gin.Context) {
	start := time.Now()
	airlines, err := h.service.GetAllAirlines()
	if err != nil {
		h.logger.Error("Failed to get airlines", zap.Error(err))

		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Debug("Retrieved airlines", zap.Int("count", len(airlines)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrived airlines, count: %v", len(airlines)), airlines, took))
}

func (h *AirlineHandler) FindByAirlineId(c *gin.Context) {
	start := time.Now()
	airlineId := c.Param("airlineId")
	airline, err := h.service.FindByAirlineId(airlineId)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to find airline with id=%s", airlineId, zap.Error(err))
		c.JSON(http.StatusNotFound,
			models.NewErrorHandler(http.StatusNotFound, "Not found", err.Error()))
		return
	}

	h.logger.Info("Find airline success", zap.String("airlineId", airlineId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Find by airlineId=%s success", airlineId), airline, took))
}

func (h *AirlineHandler) UpdateAirline(c *gin.Context) {
	start := time.Now()

	airlineId := c.Param("airline")
	var rawData models.RawAirline
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn("Invalid request payload", zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", "Invalid request payload"))
		return
	}
	airline, err := h.service.UpdateByAirlineId(airlineId, rawData)
	if err != nil {
		h.logger.Error("Failed to update airline", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorHandler(http.StatusInternalServerError, "fail", err.Error()))
		return
	}

	h.logger.Info("Airline updated successfully", zap.String("airline_id", airlineId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Airline updated successfully", airline, took))
}

func (h *AirlineHandler) DeleteByAirlineId(c *gin.Context) {
	start := time.Now()

	airlineId := c.Param("airlineId")
	err := h.service.DeleteByAirlineId(airlineId)
	if err != nil {
		h.logger.Error("Failed to delete airline: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Deleted airline successfully", zap.String("airlineId", airlineId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Deleted airline successfully", "Deleted", took))
}

func (h *AirlineHandler) ImportAirlines(c *gin.Context) {
	start := time.Now()

	airlineId := c.Param("airlineId")
	err := h.service.ImportAirlineData("")
	if err != nil {
		h.logger.Error("Failed to delete airline: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Deleted airline successfully", zap.String("airlineId", airlineId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Deleted airline successfully", "Deleted", took))
}
