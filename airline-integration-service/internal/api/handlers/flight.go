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

type FlightHandler struct {
	service *services.FlightService
	logger  *zap.Logger
}

const INVALID_PAYLOAD = "Invalid payload"

func NewFlightHandler(service *services.FlightService, logger *zap.Logger) *FlightHandler {
	return &FlightHandler{service: service, logger: logger}
}

func (h *FlightHandler) CreateFlight(c *gin.Context) {
	start := time.Now()

	var rawData models.RawFlightData
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn(INVALID_PAYLOAD, zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", INVALID_PAYLOAD))
		return
	}
	flight, err := h.service.CreateFlight(rawData)
	if err != nil {
		h.logger.Error("Failed to create flight", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "fail", fmt.Sprintf("Failed to create flight, error: %v", err.Error())))
		return
	}

	h.logger.Info("Flight created successfully", zap.String("flight_id", rawData.Flight.IATA))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusCreated, models.NewApiResponse(http.StatusCreated, "success", "Flight created successfully", flight, took))
}

func (h *FlightHandler) GetAllFlights(c *gin.Context) {
	start := time.Now()
	flights, err := h.service.GetAllFlights()
	if err != nil {
		h.logger.Error("Failed to get flights", zap.Error(err))

		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Debug("Retrieved flights", zap.Int("count", len(flights)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved flights, count: %v", len(flights)), flights, took))
}

func (h *FlightHandler) FindByFlightId(c *gin.Context) {
	start := time.Now()
	flightId := c.Param("flightId")
	flight, err := h.service.FindByFlightId(flightId)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to find flight with id=%s", flightId, zap.Error(err))
		c.JSON(http.StatusNotFound,
			models.NewErrorHandler(http.StatusNotFound, "Not found", err.Error()))
		return
	}

	h.logger.Info("Find flight success", zap.String("flightId", flightId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Find by flightId=%s success", flightId), flight, took))
}

func (h *FlightHandler) UpdateFlight(c *gin.Context) {
	start := time.Now()

	flightId := c.Param("flight")
	var rawData models.RawFlightData
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn(INVALID_PAYLOAD, zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", INVALID_PAYLOAD))
		return
	}
	flight, err := h.service.UpdateByFlightId(flightId, rawData)
	if err != nil {
		h.logger.Error("Failed to update flight", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorHandler(http.StatusInternalServerError, "fail", err.Error()))
		return
	}

	h.logger.Info("Flight updated successfully", zap.String("flight_id", flightId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Flight updated successfully", flight, took))
}

func (h *FlightHandler) DeleteByFlightId(c *gin.Context) {
	start := time.Now()

	flightId := c.Param("flightId")
	err := h.service.DeleteByFlightId(flightId)
	if err != nil {
		h.logger.Error("Failed to delete flight: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Deleted flight successfully", zap.String("flightId", flightId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Deleted flight successfully", "Deleted", took))
}
