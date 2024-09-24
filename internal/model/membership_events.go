package model

import (
	"time"

	"github.com/genefriendway/onchain-handler/internal/dto"
)

type MembershipEvents struct {
	ID              uint64    `json:"id" gorm:"primaryKey"`
	UserAddress     string    `json:"user_address" gorm:"column:user_address"`
	OrderID         uint64    `json:"order_id" gorm:"column:order_id"`
	TransactionHash string    `json:"transaction_hash" gorm:"column:transaction_hash"`
	Amount          string    `json:"amount" gorm:"column:amount"`
	Status          uint8     `json:"status" gorm:"column:status"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at"`
	EndDuration     time.Time `json:"end_duration" gorm:"column:end_duration"`
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
