package reward

import (
	"context"
	"fmt"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/model"
)

type rewardUCase struct {
	RewardRepository interfaces.RewardRepository
}

func NewRewardUCase(rewardRepository interfaces.RewardRepository) interfaces.RewardUCase {
	return &rewardUCase{
		RewardRepository: rewardRepository,
	}
}

func (u *rewardUCase) CreateRewardsHistory(ctx context.Context, payloads []dto.Reward) error {
	// Convert DTO rewards to model rewards
	rewardModels := make([]model.Reward, len(payloads))
	for i, payload := range payloads {
		rewardModels[i] = model.Reward{
			RewardAddress:    payload.RewardAddress,
			RecipientAddress: payload.RecipientAddress,
			TransactionHash:  payload.TransactionHash,
			TokenAmount:      payload.TokenAmount,
			Status:           1,
		}
	}

	err := u.RewardRepository.CreateRewardsHistory(ctx, rewardModels)
	if err != nil {
		return fmt.Errorf("cannot create rewards history: %w", err)
	}

	return nil
}
