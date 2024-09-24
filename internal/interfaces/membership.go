package interfaces

import (
	"context"

	"github.com/genefriendway/onchain-handler/internal/model"
)

type MembershipRepository interface {
	CreateMembershipEventHistory(ctx context.Context, model model.MembershipEvents) error
}

type MembershipUCase interface{}
