package membership

import (
	"github.com/genefriendway/onchain-handler/internal/interfaces"
)

type membershipUCase struct {
	MembershipPurchaseRepository interfaces.MembershipRepository
}

func NewMembershipPurchaseUCase(membershipPurchaseRepository interfaces.MembershipRepository) interfaces.MembershipUCase {
	return &membershipUCase{
		MembershipPurchaseRepository: membershipPurchaseRepository,
	}
}
