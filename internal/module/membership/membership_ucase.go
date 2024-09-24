package membership

import (
	"context"

	"github.com/genefriendway/onchain-handler/internal/dto"
	"github.com/genefriendway/onchain-handler/internal/interfaces"
)

type membershipUCase struct {
	MembershipRepository interfaces.MembershipRepository
}

func NewMembershipUCase(membershipRepository interfaces.MembershipRepository) interfaces.MembershipUCase {
	return &membershipUCase{
		MembershipRepository: membershipRepository,
	}
}

func (u *membershipUCase) GetMembershipEventByOrderID(ctx context.Context, orderID uint64) (*dto.MembershipEventsDTO, error) {
	membershipEvent, err := u.MembershipRepository.GetMembershipEventByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if membershipEvent == nil {
		return nil, nil
	}
	membershipEventDTO := membershipEvent.ToDto()
	return &membershipEventDTO, nil
}
