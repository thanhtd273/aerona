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

type AirportHandler struct {
	service *services.AirportService
	logger  *zap.Logger
}

func NewAirportHandler(service *services.AirportService, logger *zap.Logger) *AirportHandler {
	return &AirportHandler{service: service, logger: logger}
}

func (h *AirportHandler) CreateAirport(c *gin.Context) {
	start := time.Now()

	var rawData models.RawAirport
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn(INVALID_PAYLOAD, zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", INVALID_PAYLOAD))
		return
	}
	airport, err := h.service.CreateAirport(rawData)
	if err != nil {
		h.logger.Error("Failed to create airport", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "fail", fmt.Sprintf("Failed to create airport, error: %v", err.Error())))
		return
	}

	h.logger.Info("Airport created successfully", zap.String("airport_id", rawData.IATACode))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusCreated, models.NewApiResponse(http.StatusCreated, "success", "Airport created successfully", airport, took))
}

func (h *AirportHandler) GetAllAirports(c *gin.Context) {
	start := time.Now()
	airports, err := h.service.GetAllAirports()
	if err != nil {
		h.logger.Error("Failed to get airports", zap.Error(err))

		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Debug("Retrieved airports", zap.Int("count", len(airports)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved airports, count: %v", len(airports)), airports, took))
}

func (h *AirportHandler) FindByAirportId(c *gin.Context) {
	start := time.Now()
	airportId := c.Param("airportId")
	airport, err := h.service.FindByAirportId(airportId)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to find airport with id=%s", airportId, zap.Error(err))
		c.JSON(http.StatusNotFound,
			models.NewErrorHandler(http.StatusNotFound, "Not found", err.Error()))
		return
	}

	h.logger.Info("Find airport success", zap.String("airportId", airportId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Find by airportId=%s success", airportId), airport, took))
}

func (h *AirportHandler) UpdateAirport(c *gin.Context) {
	start := time.Now()

	airportId := c.Param("airportId")
	var rawData models.RawAirport
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn(INVALID_PAYLOAD, zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", INVALID_PAYLOAD))
		return
	}
	airport, err := h.service.UpdateAirport(airportId, rawData)
	if err != nil {
		h.logger.Error("Failed to update airport", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorHandler(http.StatusInternalServerError, "fail", err.Error()))
		return
	}

	h.logger.Info("Airport updated successfully", zap.String("airport_id", airportId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Airport updated successfully", airport, took))
}

func (h *AirportHandler) DeleteByAirportId(c *gin.Context) {
	start := time.Now()

	airportId := c.Param("airportId")
	err := h.service.DeleteByAirportId(airportId)
	if err != nil {
		h.logger.Error("Failed to delete airport: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Deleted airport successfully", zap.String("airportId", airportId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Deleted airport successfully", "Deleted", took))
}
