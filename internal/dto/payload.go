package dto

type CreateRewardPayload struct {
	RecipientAddress string `json:"recipient_address"`
	TokenAmount      string `json:"token_amount"`
}
