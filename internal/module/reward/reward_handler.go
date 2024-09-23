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

// NewRewardHandler initializes the RewardHandler
func NewRewardHandler(ucase interfaces.RewardUCase, config *conf.Configuration) *RewardHandler {
	// Connect to eth client
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

// Reward handles the distribution of reward tokens
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
	// Parse the request body
	var req []dto.CreateRewardPayload
	if err := ctx.BindJSON(&req); err != nil {
		log.LG.Errorf("Failed to parse reward payload: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Convert the payload into the recipients map (address -> token amount)
	recipients, err := convertToRecipients(req)
	if err != nil {
		log.LG.Errorf("Failed to convert recipients: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Perform the reward distribution
	txHash, err := DistributeReward(h.ETHClient, h.Config, recipients)
	if err != nil {
		log.LG.Errorf("Failed to distribute rewards: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Reward distribution failed"})
		return
	}

	// Prepare the rewards history
	rewards, err := h.prepareRewardHistory(req, *txHash)
	if err != nil {
		log.LG.Errorf("Error preparing reward history: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Save the rewards history to the database
	if err := h.UCase.CreateRewardsHistory(ctx, rewards); err != nil {
		log.LG.Errorf("Failed to save rewards history: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save rewards history"})
		return
	}

	// Return a success response
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

// Helper to convert CreateRewardPayload to recipients map
func convertToRecipients(req []dto.CreateRewardPayload) (map[string]*big.Int, error) {
	recipients := make(map[string]*big.Int)

	for _, payload := range req {
		// Check for duplicate recipient addresses
		if _, exists := recipients[payload.RecipientAddress]; exists {
			return nil, fmt.Errorf("duplicate recipient address: %s", payload.RecipientAddress)
		}

		// Convert token amount to big.Int
		tokenAmount := new(big.Int)
		if _, success := tokenAmount.SetString(payload.TokenAmount, 10); !success {
			return nil, fmt.Errorf("invalid token amount: %s", payload.TokenAmount)
		}

		recipients[payload.RecipientAddress] = tokenAmount
	}

	return recipients, nil
}

// prepareRewardHistory creates a slice of Reward structs from the request payload
func (h *RewardHandler) prepareRewardHistory(req []dto.CreateRewardPayload, txHash string) ([]dto.Reward, error) {
	rewardAddress := h.Config.Blockchain.RewardAddress
	var rewards []dto.Reward

	for _, payload := range req {
		// Convert the token amount from string to *big.Int
		tokenAmount := new(big.Int)
		if _, success := tokenAmount.SetString(payload.TokenAmount, 10); !success {
			return nil, fmt.Errorf("invalid token amount: %s", payload.TokenAmount)
		}

		// Create the Reward object
		reward := dto.Reward{
			RecipientAddress: payload.RecipientAddress,
			RewardAddress:    rewardAddress,
			TokenAmount:      tokenAmount,
			TransactionHash:  txHash,
		}
		rewards = append(rewards, reward)
	}

	return rewards, nil
}
