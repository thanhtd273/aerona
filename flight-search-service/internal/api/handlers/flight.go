package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/dto"
	"aerona.thanhtd.com/flight-search-service/internal/api/services"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"aerona.thanhtd.com/flight-search-service/internal/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FlightHandler struct {
	service *services.FlightService
	logger  *zap.Logger
}

func NewFlightHandler(service *services.FlightService, logger *zap.Logger) *FlightHandler {
	return &FlightHandler{service: service, logger: logger}
}

func (h *FlightHandler) SearchFlights(c *gin.Context) {
	start := time.Now()
	offsetStr := c.DefaultQuery("offset", "1")
	if !validator.IsNumeric(offsetStr) {
		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "fail", "Invalid offset"))
		return
	}
	offset, _ := strconv.Atoi(offsetStr)
	limitStr := c.DefaultQuery("limit", "10")
	if !validator.IsNumeric(limitStr) {
		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "fail", "Invalid limit"))
		return
	}
	limit, _ := strconv.Atoi(limitStr)

	amountStr := c.DefaultQuery("amount", "1")
	if !validator.IsNumeric(amountStr) {
		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "fail", "Invalid amount"))
		return
	}
	amount, _ := strconv.Atoi(amountStr)

	searchInfo := dto.SearchInfo{
		From:            c.Query("from"),
		To:              c.Query("to"),
		DepartureDate:   c.Query("departure_date"),
		NumOfPassengers: amount,
	}
	h.logger.Debug("flight search query", zap.Any("query", searchInfo))
	browserId := utils.GenerateBrowserId(c.Request)
	flights, err := h.service.SearchFlights(browserId, searchInfo, offset, limit)
	if err != nil {
		h.logger.Error("Failed to search flights", zap.Error(err))

		c.JSON(http.StatusInternalServerError,
			utils.NewErrorHandler(http.StatusInternalServerError, "Failed to search flights", err.Error()))
		return
	}

	h.logger.Debug("Retrieved flights", zap.Int("count", len(flights)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved flights, count: %v", len(flights)), flights, took))
}

func (h *FlightHandler) FilterFlights(c *gin.Context) {
	start := time.Now()

	var searchInfo dto.SearchInfo
	if err := c.ShouldBindQuery(&searchInfo); err != nil {
		h.logger.Warn("Invalid query parameters for filtering flights", zap.Error(err))
		c.JSON(http.StatusBadRequest, utils.NewErrorHandler(http.StatusBadRequest, "fail", "Invalid query parameters"))
		return
	}

	if searchInfo.From == "" || searchInfo.To == "" || searchInfo.DepartureDate == "" {
		h.logger.Warn("Missing required fields in filter flights request", zap.Any("searchInfo", searchInfo))
		c.JSON(http.StatusBadRequest, utils.NewErrorHandler(400, "Bad Request", "Missing required fields"))
		return
	}

	flights, err := h.service.FilterFlights(searchInfo)

	if err != nil {
		h.logger.Error("Failed to search flights", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}
	h.logger.Info("Successfully retrieved filtered flights", zap.Int("count", len(flights)), zap.Duration("duration", time.Since(start)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved flights, count: %v", len(flights)), flights, took))
}

func (h *FlightHandler) GetAllFlights(c *gin.Context) {
	start := time.Now()
	flights, err := h.service.GetAllFlights()
	if err != nil {
		h.logger.Error("Failed to get flights", zap.Error(err))

		c.JSON(http.StatusInternalServerError,
			utils.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Debug("Retrieved flights", zap.Int("count", len(flights)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved flights, count: %v", len(flights)), flights, took))
}

func (h *FlightHandler) FindByFlightId(c *gin.Context) {
	start := time.Now()
	flightId := c.Param("flightId")
	flight, err := h.service.FindByFlightId(flightId)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to find flight with id=%s", flightId, zap.Error(err))
		c.JSON(http.StatusNotFound,
			utils.NewErrorHandler(http.StatusNotFound, "Not found", err.Error()))
		return
	}

	h.logger.Info("Find flight success", zap.String("flightId", flightId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Find by flightId=%s success", flightId), flight, took))
}
