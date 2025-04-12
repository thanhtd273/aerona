package handlers

import (
	"net/http"
	"time"

	"aerona.thanhtd.com/airline-integration-service/internal/api/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MonitorHandler struct {
	logger *zap.Logger
}

func NewMonitorHandler(logger *zap.Logger) *MonitorHandler {
	return &MonitorHandler{
		logger: logger,
	}
}

func (h *MonitorHandler) HealthCheck(c *gin.Context) {
	start := time.Now()
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusCreated, models.NewApiResponse(http.StatusOK, "success", "Service is healthy", "Health", took))
}
