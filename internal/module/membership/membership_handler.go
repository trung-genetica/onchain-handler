package membership

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

// GetMembershipEventsByOrderIDs retrieves membership events by a list of order IDs.
// @Summary Get membership events by order IDs
// @Description Get a list of membership events based on provided order IDs
// @Tags membership
// @Accept json
// @Produce json
// @Param orderIds query string true "Comma-separated list of Order IDs"
// @Success 200 {array} dto.MembershipEventsDTO
// @Failure 400 {object} util.GeneralError "Invalid Order IDs"
// @Failure 404 {object} util.GeneralError "Membership events not found"
// @Failure 500 {object} util.GeneralError "Internal server error"
// @Router /api/v1/membership/events [get]
func (h *MembershipHandler) GetMembershipEventsByOrderIDs(ctx *gin.Context) {
	// Extract order IDs from query params and split by comma
	orderIDsStr := ctx.Query("orderIds")
	if orderIDsStr == "" {
		log.LG.Errorf("Order IDs are required")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order IDs are required"})
		return
	}

	// Split the comma-separated IDs and parse them into uint64
	orderIDs, err := parseOrderIDs(orderIDsStr)
	if err != nil {
		log.LG.Errorf("Invalid Order IDs: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Order IDs"})
		return
	}

	// Fetch the membership events using the use case
	events, err := h.UCase.GetMembershipEventsByOrderIDs(ctx, orderIDs)
	if err != nil {
		log.LG.Errorf("Failed to retrieve membership events: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// If no events are found, return a 404 response
	if events == nil {
		log.LG.Error("No membership events found for provided Order IDs")
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Membership events not found"})
		return
	}

	// Return the event data as a JSON response
	ctx.JSON(http.StatusOK, events)
}

// parseOrderIDs parses a comma-separated string of order IDs into a slice of uint64.
func parseOrderIDs(orderIDsStr string) ([]uint64, error) {
	var orderIDs []uint64
	idStrs := strings.Split(orderIDsStr, ",")
	for _, idStr := range idStrs {
		orderID, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid Order ID: %v", err)
		}
		orderIDs = append(orderIDs, orderID)
	}
	return orderIDs, nil
}
