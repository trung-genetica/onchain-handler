package reward

import (
	"context"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/interfaces"
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
	// TODO: implement here
	return nil
}
