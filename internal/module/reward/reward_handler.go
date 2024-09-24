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
// @Failure 400 		{object}	util.GeneralError "Invalid payload"
// @Failure 500 		{object}	util.GeneralError "Internal server error"
// @Router 	/api/v1/rewards [post]
func (h *RewardHandler) Reward(ctx *gin.Context) {
	var req []dto.CreateRewardPayloadDTO

	// Step 1: Parse and validate the request payload
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.LG.Errorf("%s: %v", "Invalid payload", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payload",
			"details": err.Error(),
		})
		return
	}

	if err := h.UCase.DistributeRewards(ctx, req); err != nil {
		log.LG.Errorf("Failed to distribute rewards: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to distribute rewards",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
