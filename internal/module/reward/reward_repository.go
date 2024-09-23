package reward

import (
	"context"

	"gorm.io/gorm"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/model"
)

type rewardRepository struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) interfaces.RewardRepository {
	return &rewardRepository{
		db: db,
	}
}

func (r rewardRepository) CreateRewardsHistory(ctx context.Context, models []model.Reward) error {
	// TODO: implement here
	return nil
}
