package dto

import "math/big"

type Reward struct {
	ID               uint64   `json:"id"`
	RewardAddress    string   `json:"reward_address"`
	RecipientAddress string   `json:"recipient_address"`
	TransactionHash  string   `json:"transaction_hash"`
	TokenAmount      *big.Int `json:"token_amount"`
	Status           int16    `json:"status"`
}
