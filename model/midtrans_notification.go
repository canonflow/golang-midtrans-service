package model

type MidtransNotification struct {
	OrderId           string `json:"order_id" binding:"required"`
	GrossAmount       string `json:"gross_amount" binding:"required"`
	TransactionStatus string `json:"transaction_status"  binding:"required"`
	TransactionId     string `json:"transaction_id"  binding:"required"`
	StatusCode        string `json:"status_code"  binding:"required"`
	PaymentType       string `json:"payment_type"  binding:"required"`
	FraudStatus       string `json:"fraud_status"  binding:"required"`
	SignatureKey      string `json:"signature_key"  binding:"required"`
}
