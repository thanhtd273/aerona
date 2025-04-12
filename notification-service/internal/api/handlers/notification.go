package handlers

import (
	"fmt"
	"net/http"
	"time"

	"aerona.thanhtd.com/notification-service/internal/api/services"
	"aerona.thanhtd.com/notification-service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	logger  *zap.Logger
	service *services.NotificationService
}

func NewNotificationHandler(logger *zap.Logger, service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		logger:  logger,
		service: service,
	}
}

func (h *NotificationHandler) FindById(c *gin.Context) {
	start := time.Now()
	notificationId := c.Param("notificationId")
	notification, err := h.service.FindById(c.Request.Context(), notificationId)
	if err != nil {
		h.logger.Error("Not found notification", zap.String("notificationId", notificationId), zap.Error(err))
		c.JSON(http.StatusNotFound, utils.NewErrorHandler(http.StatusNotFound, "fail", fmt.Sprintf("Failed to find notification with id=%s, error: %v", notificationId, err)))
	}

	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", "", notification, took))
}
