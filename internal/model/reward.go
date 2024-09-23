package model

import (
	"math/big"
	"time"

	"github.com/genefriendway/onchain-handler/internal/dto"
)

const tableName = "reward"

type Reward struct {
	ID               uint64    `json:"id" gorm:"primaryKey"`
	RewardAddress    string    `json:"reward_address"`
	RecipientAddress string    `json:"recipient_address"`
	TransactionHash  string    `json:"transaction_hash"`
	TokenAmount      *big.Int  `json:"token_amount"`
	Status           int16     `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (m *Reward) TableName() string {
	return tableName
}

func (m *Reward) ToDto() dto.Reward {
	return dto.Reward{
		ID:               m.ID,
		RewardAddress:    m.RewardAddress,
		RecipientAddress: m.RecipientAddress,
		TransactionHash:  m.TransactionHash,
		TokenAmount:      m.TokenAmount,
		Status:           m.Status,
	}
}
