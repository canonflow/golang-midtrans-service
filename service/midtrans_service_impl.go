package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"golang-midtrans-service/helper"
	"golang-midtrans-service/model"
	"strconv"
)

type MidtransServiceImpl struct {
	Validate   *validator.Validate
	SnapClient *snap.Client
}

func NewMidtransServiceImpl(validate *validator.Validate, client *snap.Client) *MidtransServiceImpl {
	return &MidtransServiceImpl{
		Validate:   validate,
		SnapClient: client,
	}
}

func (service *MidtransServiceImpl) Create(c *gin.Context, request model.MidtransRequest) model.MidtransResponse {
	err := service.Validate.Struct(request)

	if err != nil {
		helper.PanicIfError(err)
	}

	// Get OrderId
	orderId := strconv.Itoa(request.OrderId)

	// Create Snap Request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderId,
			GrossAmt: request.Amount,
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.Email,
			Email: request.Email,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Expiry: &snap.ExpiryDetails{
			Unit:     "day",
			Duration: 1,
		},
	}

	snapToken, errSnap := service.SnapClient.CreateTransactionToken(req)
	if errSnap != nil {
		helper.PanicIfError(errSnap)
	}

	// Generate Resposne
	midtransResponse := model.MidtransResponse{
		SnapToken: snapToken,
	}

	return midtransResponse
}
