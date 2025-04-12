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

type CityHandler struct {
	service *services.CityService
	logger  *zap.Logger
}

func NewCityHandler(service *services.CityService, logger *zap.Logger) *CityHandler {
	return &CityHandler{service: service, logger: logger}
}

func (h *CityHandler) CreateCity(c *gin.Context) {
	start := time.Now()

	var rawData models.RawCity
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn(INVALID_PAYLOAD, zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", INVALID_PAYLOAD))
		return
	}
	city, err := h.service.CreateCity(rawData)
	if err != nil {
		h.logger.Error("Failed to create city", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "fail", fmt.Sprintf("Failed to create city, error: %v", err.Error())))
		return
	}

	h.logger.Info("City created successfully", zap.String("city_id", rawData.IATACode))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusCreated, models.NewApiResponse(http.StatusCreated, "success", "City created successfully", city, took))
}

func (h *CityHandler) GetAllCities(c *gin.Context) {
	start := time.Now()
	cities, err := h.service.GetAllCities()
	if err != nil {
		h.logger.Error("Failed to get cities", zap.Error(err))

		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Debug("Retrieved cities", zap.Int("count", len(cities)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Retrieved cities, count: %v", len(cities)), cities, took))
}

func (h *CityHandler) FindByCityId(c *gin.Context) {
	start := time.Now()
	cityId := c.Param("cityId")
	city, err := h.service.FindByCityId(cityId)
	if err != nil {
		h.logger.Sugar().Errorf("Failed to find city with id=%s", cityId, zap.Error(err))
		c.JSON(http.StatusNotFound,
			models.NewErrorHandler(http.StatusNotFound, "Not found", err.Error()))
		return
	}

	h.logger.Info("Find city success", zap.String("cityId", cityId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", fmt.Sprintf("Find by cityId=%s success", cityId), city, took))
}

func (h *CityHandler) UpdateCity(c *gin.Context) {
	start := time.Now()

	cityId := c.Param("cityId")
	var rawData models.RawCity
	if err := c.ShouldBindJSON(&rawData); err != nil {
		h.logger.Warn(INVALID_PAYLOAD, zap.Error(err))
		c.JSON(http.StatusBadRequest,
			models.NewErrorHandler(http.StatusBadRequest, "Bad request", INVALID_PAYLOAD))
		return
	}
	city, err := h.service.UpdateCity(cityId, rawData)
	if err != nil {
		h.logger.Error("Failed to update city", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.NewErrorHandler(http.StatusInternalServerError, "fail", err.Error()))
		return
	}

	h.logger.Info("City updated successfully", zap.String("city_id", cityId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "City updated successfully", city, took))
}

func (h *CityHandler) DeleteByCityId(c *gin.Context) {
	start := time.Now()

	cityId := c.Param("cityId")
	err := h.service.DeleteByCityId(cityId)
	if err != nil {
		h.logger.Error("Failed to delete city: ", zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			models.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Deleted city successfully", zap.String("cityId", cityId))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, models.NewApiResponse(http.StatusOK, "success", "Deleted city successfully", "Deleted", took))
}
