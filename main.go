package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"golang-midtrans-service/controller"
	"golang-midtrans-service/initializer"
	"golang-midtrans-service/middleware"
	"golang-midtrans-service/model"
	"golang-midtrans-service/service"
	"net/http"
	"os"
)

func init() {
	initializer.LoadEnv()
}

func main() {
	// Snap Client
	var snapClient = snap.Client{}
	snapClient.New(os.Getenv("MIDTRANS_SERVER_KEY"), midtrans.Sandbox)
	snapClient.Options.SetPaymentOverrideNotification(os.Getenv("MIDTRANS_CALLBACK_URL"))

	// Service
	validate := validator.New()
	midtransService := service.NewMidtransServiceImpl(validate, &snapClient)
	midtransController := controller.NewMidtransControllerImpl(midtransService)

	// Router
	router := gin.Default()
	router.Use(middleware.ErrorHandle())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":   http.StatusNotFound,
			"status": "Not Found",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.WebResponse{
			Code:   http.StatusOK,
			Status: "OK",
			Data: map[string]any{
				"Message": "Welcome to Midtrans Service by canonflow",
			},
		})
	})

	midtransRouting := router.Group("/midtrans")
	{
		midtransRouting.POST("/create-snap-token", midtransController.CreateSnapToken)
		midtransRouting.POST("/listen-notification", midtransController.ListenNotification)
	}

	err := router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		panic(err)
	}
}
