package dto

type CreateRewardPayload struct {
	RecipientAddress string `json:"recipient_address"`
	TokenAmount      uint64 `json:"token_amount"`
}
