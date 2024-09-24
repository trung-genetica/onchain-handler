package interfaces

import (
	"context"

	"github.com/genefriendway/onchain-handler/internal/model"
)

type MembershipPurchaseRepository interface {
	CreateMembershipEventHistory(ctx context.Context, model model.MembershipEvents) error
}

type MembershipPurchaseUCase interface{}
