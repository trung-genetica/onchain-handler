package reward

import (
	"context"
	"fmt"

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
	// Use GORM to perform a bulk insert
	err := r.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return fmt.Errorf("failed to create rewards history: %w", err)
	}
	return nil
}
