package model

import (
	"time"

	"github.com/genefriendway/onchain-handler/internal/dto"
)

type MembershipEvents struct {
	ID              uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserAddress     string    `json:"user_address"`
	OrderID         uint64    `json:"order_id"`
	TransactionHash string    `json:"transaction_hash"`
	Amount          string    `json:"amount"`
	Status          uint8     `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	EndDuration     time.Time `json:"end_duration"`
}

func (m *MembershipEvents) TableName() string {
	return "membership_events"
}

func (m *MembershipEvents) ToDto() dto.MembershipEventsDTO {
	return dto.MembershipEventsDTO{
		ID:              m.ID,
		UserAddress:     m.UserAddress,
		OrderID:         m.OrderID,
		TransactionHash: m.TransactionHash,
		Amount:          m.Amount,
		Status:          m.Status,
		EndDuration:     m.EndDuration,
	}
}
