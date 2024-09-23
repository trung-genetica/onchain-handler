package interfaces

import (
	"context"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/model"
)

type RewardRepository interface {
	CreateRewardsHistory(ctx context.Context, models []model.Reward) error
}

type RewardUCase interface {
	CreateRewardsHistory(ctx context.Context, payloads []dto.Reward) error
}
