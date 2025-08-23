package dto

import "github.com/LucasLCabral/moranGo-pay/internal/domain"

type TransactionRequest struct {
	UserID      string                 `json:"user_id"`
	Amount      float64                `json:"amount"`
	Type        domain.TransactionType `json:"type"`
	Description string                 `json:"description"`

	RecipientID  *string `json:"recipient_id,omitempty"`
	MerchantID   *string `json:"merchant_id,omitempty"`
	OriginalTxID *string `json:"original_tx_id,omitempty"`
}

type TransferRequest struct {
	FromUserID  string  `json:"from_user_id"`
	RecipientID string  `json:"recipient_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}
