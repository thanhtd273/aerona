package handlers

import (
	"net/http"
	"time"

	"aerona.thanhtd.com/ticket-service/internal/api/services"
	"aerona.thanhtd.com/ticket-service/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TicketHandler struct {
	service *services.TicketService
	logger  *zap.Logger
}

func NewTicketHandler(service *services.TicketService, logger *zap.Logger) *TicketHandler {
	return &TicketHandler{
		service: service,
		logger:  logger,
	}
}

// func (h *TicketHandler) CreateTicket(c *gin.Context) {
// 	start := time.Now()
// 	var booking dto.Booking
// 	if err := c.ShouldBindJSON(&booking); err != nil {
// 		h.logger.Warn("Invalid request body for creating ticket", zap.Error(err))
// 		c.JSON(http.StatusBadRequest, utils.NewErrorHandler(http.StatusBadRequest, "fail", "Invalid request body for creating ticket"))
// 		return
// 	}
// 	ticket, err := h.service.CreateTicket(c.Request.Context(), booking)
// 	if err != nil {
// 		h.logger.Error("Failed to create ticket", zap.Error(err))
// 		c.JSON(http.StatusInternalServerError, utils.NewErrorHandler(http.StatusInternalServerError, "fail", "Internal server error"))
// 		return
// 	}
// 	h.logger.Debug("Successfully created ticket")
// 	took := time.Since(start).Milliseconds()
// 	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", "Successfully created ticket", ticket, took))
// }

func (h *TicketHandler) GetAllTickets(c *gin.Context) {
	start := time.Now()
	tickets, err := h.service.GetAllTickets(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get all tickets", zap.Error(err))
		c.JSON(http.StatusOK, utils.NewErrorHandler(http.StatusBadRequest, "fail", "Failed to get all tickets"))
	}

	h.logger.Debug("Successfully get all tickets", zap.Int("count", len(tickets)))
	took := time.Since(start).Milliseconds()
	c.JSON(http.StatusOK, utils.NewApiResponse(http.StatusOK, "success", "Get all tickets successfully", tickets, took))
}
