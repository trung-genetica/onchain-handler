package dto

type RewardDTO struct {
	ID               uint64 `json:"id"`
	RewardAddress    string `json:"reward_address"`
	RecipientAddress string `json:"recipient_address"`
	TransactionHash  string `json:"transaction_hash"`
	TokenAmount      string `json:"token_amount"`
	Status           int16  `json:"status"`
	ErrorMessage     string `json:"error_message"`
}
