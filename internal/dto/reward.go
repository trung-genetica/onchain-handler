package dto

type Reward struct {
	ID               uint64 `json:"id"`
	RewardAddress    string `json:"reward_address"`
	RecipientAddress string `json:"recipient_address"`
	TransactionHash  string `json:"transaction_hash"`
	TokenAmount      uint64 `json:"token_amount"`
	Status           uint64 `json:"status"`
}
