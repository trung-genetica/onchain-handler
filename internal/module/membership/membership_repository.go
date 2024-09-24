package membership

import (
	"context"

	"gorm.io/gorm"

	"github.com/genefriendway/onchain-handler/internal/interfaces"
	"github.com/genefriendway/onchain-handler/internal/model"
)

type membershipPurchaseRepository struct {
	db *gorm.DB
}

func NewMembershipPurchaseRepository(db *gorm.DB) interfaces.MembershipPurchaseRepository {
	return &membershipPurchaseRepository{
		db: db,
	}
}

func (r membershipPurchaseRepository) CreateMembershipEventHistory(ctx context.Context, model model.MembershipEvents) error {
	// TODO: implement here
	return nil
}
