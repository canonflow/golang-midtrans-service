package controller

import (
	"github.com/gin-gonic/gin"
	"golang-midtrans-service/helper"
	"golang-midtrans-service/model"
	"golang-midtrans-service/service"
	"net/http"
)

type MidtransControllerImpl struct {
	MidtransService service.MidtransService
}

func NewMidtransControllerImpl(midtransService service.MidtransService) *MidtransControllerImpl {
	return &MidtransControllerImpl{MidtransService: midtransService}
}

func (controller *MidtransControllerImpl) CreateSnapToken(c *gin.Context) {
	var request model.MidtransRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		helper.PanicIfError(err)
	}

	midtransResponse := controller.MidtransService.Create(c, request)
	webResponse := model.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   midtransResponse,
	}

	c.JSON(http.StatusOK, webResponse)
}
