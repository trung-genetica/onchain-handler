package reward

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"

	"github.com/genefriendway/onchain-handler/conf"
	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/utils/log"
)

type RewardHandler struct {
	UCase     interfaces.RewardUCase
	ETHClient *ethclient.Client
	Config    *conf.Configuration
}

func NewRewardHandler(ucase interfaces.RewardUCase, config *conf.Configuration) *RewardHandler {
	// Initialize the eth client
	client, err := ethclient.Dial(config.Blockchain.RpcUrl)
	if err != nil {
		log.LG.Fatalf("failed to connect to eth client: %v", err)
		return nil
	}

	return &RewardHandler{
		UCase:     ucase,
		ETHClient: client,
		Config:    config,
	}
}

// Reward Distribute reward tokens
// @Summary Reward
// @Description Reward
// @Tags 	reward
// @Accept	json
// @Produce json
// @Param 	payload	body 			[]dto.CreateRewardPayload true "Request reward tokens, required"
// @Success 200 		{object}	[]dto.CreateRewardPayload "When success, return {"success": true}"
// @Failure 424 		{object}	util.GeneralError
// @Failure 417 		{object}	util.GeneralError
// @Router 	/api/v1/rewards [post]
func (h *RewardHandler) Reward(ctx *gin.Context) {
	// Parse the incoming JSON request
	var req []dto.CreateRewardPayload
	if err := ctx.BindJSON(&req); err != nil {
		log.LG.Errorf("Failed to parse reward payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payload",
		})
		return
	}

	// Convert the payload into the recipients map (address -> token amount)
	recipients, err := convertToRecipients(req)
	if err != nil {
		log.LG.Errorf("Failed to convert recipients: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process reward recipients",
		})
		return
	}

	// Distribute the rewards using the bulk transfer function
	txHash, err := DistributeReward(h.ETHClient, h.Config, recipients)
	if err != nil {
		log.LG.Errorf("Failed to distribute rewards: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Reward distribution failed",
		})
		return
	}

	// Prepare reward history data
	var rewards []dto.Reward
	rewardAddress := h.Config.Blockchain.RewardAddress
	for _, payload := range req {
		reward := dto.Reward{
			RecipientAddress: payload.RecipientAddress,
			RewardAddress:    rewardAddress,
			TokenAmount:      payload.TokenAmount,
			TransactionHash:  *txHash,
		}
		rewards = append(rewards, reward)
	}

	// Save the rewards history in the database
	if err := h.UCase.CreateRewardsHistory(ctx, rewards); err != nil {
		log.LG.Errorf("Failed to save rewards history: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save rewards history",
		})
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// Helper to convert CreateRewardPayload to recipients map
func convertToRecipients(req []dto.CreateRewardPayload) (map[string]*big.Int, error) {
	recipients := make(map[string]*big.Int)
	for _, payload := range req {
		amount := new(big.Int).SetUint64(payload.TokenAmount)
		if _, ok := recipients[payload.RecipientAddress]; ok {
			return nil, fmt.Errorf("duplicate recipient address: %s", payload.RecipientAddress)
		}
		recipients[payload.RecipientAddress] = amount
	}
	return recipients, nil
}
