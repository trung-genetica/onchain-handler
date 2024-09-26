package transfer

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

type TransferHandler struct {
	UCase interfaces.TransferUCase
}

// NewRewardHandler initializes the RewardHandler
func NewTransferHandler(ucase interfaces.TransferUCase) *TransferHandler {
	return &TransferHandler{
		UCase: ucase,
	}
}

// Transfer handles the distribution of tokens
// @Summary Transfer
// @Description Transfer
// @Tags 	transfer
// @Accept	json
// @Produce json
// @Param 	payload	body 			[]dto.TransferTokenPayloadDTO true "Request transfer tokens, required"
// @Success 200 		{object}	[]dto.TransferTokenPayloadDTO "When success, return {"success": true}"
// @Failure 400 		{object}	util.GeneralError "Invalid payload"
// @Failure 500 		{object}	util.GeneralError "Internal server error"
// @Router 	/api/v1/transfer [post]
func (h *TransferHandler) Transfer(ctx *gin.Context) {
	var req []dto.TransferTokenPayloadDTO

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.LG.Errorf("%s: %v", "Invalid payload", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid payload",
			"details": err.Error(),
		})
		return
	}

	for _, payload := range req {
		if payload.TxType != "PURCHASE" && payload.TxType != "COMMISSION" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid tx_type",
				"details": "TxType must be either PURCHASE or COMMISSION",
			})
			return
		}
	}

	if err := h.UCase.DistributeTokens(ctx, req); err != nil {
		log.LG.Errorf("Failed to distribute rewards: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to distribute rewards",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
