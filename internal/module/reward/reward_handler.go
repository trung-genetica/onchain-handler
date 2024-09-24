package reward

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

type RewardHandler struct {
	UCase interfaces.RewardUCase
}

// NewRewardHandler initializes the RewardHandler
func NewRewardHandler(ucase interfaces.RewardUCase) *RewardHandler {
	return &RewardHandler{
		UCase: ucase,
	}
}

// Reward handles the distribution of reward tokens
// @Summary Reward
// @Description Reward
// @Tags 	reward
// @Accept	json
// @Produce json
// @Param 	payload	body 			[]dto.CreateRewardPayloadDTO true "Request reward tokens, required"
// @Success 200 		{object}	[]dto.CreateRewardPayloadDTO "When success, return {"success": true}"
// @Failure 424 		{object}	util.GeneralError
// @Failure 417 		{object}	util.GeneralError
// @Router 	/api/v1/rewards [post]
func (h *RewardHandler) Reward(ctx *gin.Context) {
	var req []dto.CreateRewardPayloadDTO
	if err := ctx.BindJSON(&req); err != nil {
		log.LG.Errorf("Failed to parse reward payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Call the use case to distribute rewards
	err := h.UCase.DistributeRewards(ctx, req)
	if err != nil {
		log.LG.Errorf("Failed to distribute rewards: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
