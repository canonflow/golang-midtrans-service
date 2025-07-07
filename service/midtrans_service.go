package service

import (
	"github.com/gin-gonic/gin"
	"golang-midtrans-service/model"
)

type MidtransService interface {
	Create(c *gin.Context, request model.MidtransRequest) model.MidtransResponse
}
