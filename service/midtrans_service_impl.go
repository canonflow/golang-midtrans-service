package service

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"golang-midtrans-service/helper"
	"golang-midtrans-service/model"
	"net/http"
	"os"
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

func (service *MidtransServiceImpl) VerifySignatureKey(request model.MidtransNotification) bool {
	err := service.Validate.Struct(request)

	if err != nil {
		helper.PanicIfError(err)
	}

	// Get all attributes
	orderId := request.OrderId
	statusCode := request.StatusCode
	grossAmount := request.GrossAmount
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	combineSignature := orderId + statusCode + grossAmount + serverKey
	hash := sha512.New()
	hash.Write([]byte(combineSignature))
	hasString := hex.EncodeToString(hash.Sum(nil))

	return request.SignatureKey == hasString
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

func (service *MidtransServiceImpl) Notification(c *gin.Context, request model.MidtransNotification) model.WebResponse {
	err := service.Validate.Struct(request)

	if err != nil {
		helper.PanicIfError(err)
	}
	var message string
	paymentType := request.PaymentType
	orderId := request.OrderId

	// Check Status
	switch request.TransactionStatus {
	case "capture":
		if request.FraudStatus == "accept" {
			// Send Request to PHP (success)
			message = "success"
		}
		break
	case "settlement":
		// Send Request to PHP (success)
		message = "success"
		break
	case "cancel", "expire":
		// Send Request to PHP (failure)
		message = "expired"
		break
	case "pending":
		// Send Request to PHP (pending)
		message = "pending"
		break
	}

	// SEND TO PHP
	requestData := map[string]string{
		"order_id":     orderId,
		"status":       message,
		"payment_type": paymentType,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		helper.PanicIfError(err)
	}

	// Prepare POST request
	req, err := http.NewRequest(
		"POST",
		"https://ubaya.xyz/flutter/160422065/project/updatestatuspembelian.php",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		helper.PanicIfError(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		helper.PanicIfError(err)
	}
	defer resp.Body.Close()

	// Unmarshal into WebResponse
	var webResp model.WebResponse
	if err := json.NewDecoder(resp.Body).Decode(&webResp); err != nil {
		helper.PanicIfError(err)
	}

	return webResp

	//return model.WebResponse{
	//	Code:   http.StatusOK,
	//	Status: "OK",
	//	Data: map[string]string{
	//		"message": message,
	//	},
	//}
}
