package model

import (
	"time"

	"github.com/genefriendway/onchain-handler/internal/dto"
)

type Reward struct {
	ID               uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	RewardAddress    string    `json:"reward_address"`
	RecipientAddress string    `json:"recipient_address"`
	TransactionHash  string    `json:"transaction_hash"`
	TokenAmount      string    `json:"token_amount"`
	Status           int16     `json:"status"`
	ErrorMessage     string    `json:"error_message"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (m *Reward) TableName() string {
	return "reward"
}

func (m *Reward) ToDto() dto.RewardDTO {
	return dto.RewardDTO{
		ID:               m.ID,
		RewardAddress:    m.RewardAddress,
		RecipientAddress: m.RecipientAddress,
		TransactionHash:  m.TransactionHash,
		TokenAmount:      m.TokenAmount,
		Status:           m.Status,
		ErrorMessage:     m.ErrorMessage,
	}
}
