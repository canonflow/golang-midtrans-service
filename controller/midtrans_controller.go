package controller

import "github.com/gin-gonic/gin"

type MidtransController interface {
	CreateSnapToken(c *gin.Context)
	ListenNotification(c *gin.Context)
}
