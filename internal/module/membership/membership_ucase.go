package membership

import (
	"github.com/genefriendway/onchain-handler/internal/interfaces"
)

type membershipPurchaseUCase struct {
	MembershipPurchaseRepository interfaces.MembershipPurchaseRepository
}

func NewMembershipPurchaseUCase(membershipPurchaseRepository interfaces.MembershipPurchaseRepository) interfaces.MembershipPurchaseUCase {
	return &membershipPurchaseUCase{
		MembershipPurchaseRepository: membershipPurchaseRepository,
	}
}
