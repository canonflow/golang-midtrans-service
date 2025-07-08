package model

type MidtransRequest struct {
	OrderId string `json:"order_id" binding:"required"`
	Amount  int64  `json:"amount" binding:"required"`
	Email   string `json:"email" binding:"required"`
}
