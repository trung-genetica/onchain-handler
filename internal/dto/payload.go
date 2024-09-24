package dto

type CreateRewardPayloadDTO struct {
	RecipientAddress string `json:"recipient_address"`
	TokenAmount      string `json:"token_amount"`
}
