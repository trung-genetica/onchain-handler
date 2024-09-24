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

func (r membershipRepository) CreateMembershipEventHistory(ctx context.Context, model model.MembershipEvents) error {
	// TODO: implement here
	return nil
}
