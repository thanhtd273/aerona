package handlers

import (
	"net/http"
	"time"

	"aerona.thanhtd.com/flight-search-service/internal/api/services"
	"aerona.thanhtd.com/flight-search-service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SearchHistoryHandler struct {
	service *services.SearchHistoryService
	logger  *zap.Logger
}

func NewSearchHistoryHandler(service *services.SearchHistoryService, logger *zap.Logger) *SearchHistoryHandler {
	return &SearchHistoryHandler{service: service, logger: logger}
}

func (h *SearchHistoryHandler) DeleteBySearchId(c *gin.Context) {
	start := time.Now()
	searchId := c.Param("searchId")

	err := h.service.DeleteBySearchId(searchId)
	if err != nil {
		h.logger.Error("Failed to delete search history record",
			zap.String("searchId", searchId),
			zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Successfully deleted search history record",
		zap.String("searchId", searchId),
		zap.Duration("duration", time.Since(start)))
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success",
		"Search history record deleted successfully", nil, time.Since(start).Milliseconds()))
}

func (h *SearchHistoryHandler) GetRecentSearches(c *gin.Context) {
	start := time.Now()
	recentSearches, err := h.service.GetRecentSearches()

	if err != nil {
		h.logger.Error("Failed to retrieve recent searches",
			zap.Error(err))
		c.JSON(http.StatusInternalServerError,
			utils.NewErrorHandler(http.StatusInternalServerError, "Internal server error", err.Error()))
		return
	}

	h.logger.Info("Successfully retrieved recent searches",
		zap.Int("count", len(recentSearches)),
		zap.Duration("duration", time.Since(start)))
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success",
		"Recent searches retrieved successfully", recentSearches, time.Since(start).Milliseconds()))
}
