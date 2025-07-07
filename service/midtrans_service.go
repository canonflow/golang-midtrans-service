package service

import (
	"github.com/gin-gonic/gin"
	"golang-midtrans-service/model"
)

type MidtransService interface {
	Create(c *gin.Context, request model.MidtransRequest) model.MidtransResponse
	VerifySignatureKey(request model.MidtransNotification) bool
	Notification(c *gin.Context, request model.MidtransNotification) model.WebResponse
}
