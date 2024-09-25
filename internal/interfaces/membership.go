package interfaces

import (
	"context"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/model"
)

type MembershipRepository interface {
	CreateMembershipEventHistory(ctx context.Context, model model.MembershipEvents) error
	GetMembershipEventByOrderID(ctx context.Context, orderID uint64) (*model.MembershipEvents, error)
	GetMembershipEventsByOrderIDs(ctx context.Context, orderIDs []uint64) ([]model.MembershipEvents, error)
}

type MembershipUCase interface {
	GetMembershipEventByOrderID(ctx context.Context, orderID uint64) (*dto.MembershipEventsDTO, error)
	GetMembershipEventsByOrderIDs(ctx context.Context, orderIDs []uint64) ([]dto.MembershipEventsDTO, error)
}
