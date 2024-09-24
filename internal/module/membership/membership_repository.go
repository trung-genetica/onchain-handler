package membership

import (
	"context"

	"gorm.io/gorm"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/model"
)

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) interfaces.MembershipRepository {
	return &membershipRepository{
		db: db,
	}
}

func (r *membershipRepository) CreateMembershipEventHistory(ctx context.Context, membershipEvent model.MembershipEvents) error {
	if err := r.db.WithContext(ctx).Create(&membershipEvent).Error; err != nil {
		return err
	}
	return nil
}

func (r *membershipRepository) GetMembershipEventByOrderID(ctx context.Context, orderID uint64) (*model.MembershipEvents, error) {
	var membershipEvent model.MembershipEvents
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&membershipEvent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &membershipEvent, nil
}
