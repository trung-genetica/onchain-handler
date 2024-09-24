package membership

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

type MembershipHandler struct {
	UCase interfaces.MembershipUCase
}

// NewMembershipHandler initializes the MembershipHandler
func NewMembershipHandler(ucase interfaces.MembershipUCase) *MembershipHandler {
	return &MembershipHandler{
		UCase: ucase,
	}
}

// GetMembershipEventByOrderID retrieves a membership event by order ID
// @Summary Get membership event by order ID
// @Description Get membership event by order ID
// @Tags 	membership
// @Accept	json
// @Produce json
// @Param orderId query string true "Order ID to query"
// @Success 200 		{object}	dto.MembershipEventsDTO
// @Failure 400 		{object}	util.GeneralError "Invalid Order ID"
// @Failure 404 		{object}	util.GeneralError "Membership event not found"
// @Failure 500 		{object}	util.GeneralError "Internal server error"
// @Router 	/api/v1/membership [get]
func (h *MembershipHandler) GetMembershipEventByOrderID(ctx *gin.Context) {
	// Extract orderId from query params and convert to uint64
	orderIdStr := ctx.Query("orderId")
	if orderIdStr == "" {
		log.LG.Errorf("Order ID is required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
		return
	}

	// Convert orderId to uint64
	orderId, err := strconv.ParseUint(orderIdStr, 10, 64)
	if err != nil {
		log.LG.Errorf("Invalid Order ID: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order ID"})
		return
	}

	// Fetch the membership event using the use case
	event, err := h.UCase.GetMembershipEventByOrderID(ctx, orderId)
	if err != nil {
		log.LG.Errorf("Failed to retrieve membership event: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// If no event is found, return a 404 response
	if event == nil {
		log.LG.Errorf("Membership event not found for Order ID: %d", orderId)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Membership event not found"})
		return
	}

	// Return the event data as a JSON response
	ctx.JSON(http.StatusOK, event)
}
