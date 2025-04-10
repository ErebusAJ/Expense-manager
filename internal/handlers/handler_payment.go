package handlers

import (
	"os"

	"github.com/ErebusAJ/expense-manager/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	razorpay "github.com/razorpay/razorpay-go"
	razorpay_utils "github.com/razorpay/razorpay-go/utils"
)

// createRazorPayOrder
// creates a payment orderID for razorpay
func (cfg *apiConfig) createRazorPayOrder(c *gin.Context) {
	var reqDetails struct {
		Amount     float64 `json:"amount" binding:"required"`
		Currency   string  `json:"currency" binding:"required"`
		PartialPay bool    `json:"partial_pay"`
	}

	err := c.BindJSON(&reqDetails)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	godotenv.Load()
	key := os.Getenv("RAZOR_KEY")
	secret := os.Getenv("RAZOR_SECRET")
	if key == "" || secret == "" {
		utils.ErrorJSON(c, 500, utils.InternalError, "error retrieving env vars", nil)
		return
	}

	client := razorpay.NewClient(key, secret)

	razorPayOptions := map[string]interface{}{
		"amount":          reqDetails.Amount,
		"currency":        reqDetails.Currency,
		"partial_payment": reqDetails.PartialPay,
	}

	body, err := client.Order.Create(razorPayOptions, nil)
	if err != nil {
		utils.ErrorJSON(c, 500, utils.InternalError, "error creating razorpay orderID", err)
		return
	}

	c.IndentedJSON(200, body)
}

// verifyRazorPayPayment
// verifies payment made by user
func (cfg *apiConfig) verifyRazorPayPayment(c *gin.Context) {
	var reqDetails struct {
		Signature string `json:"razorpay_signature" binding:"required"`
		OrderID   string `json:"razorpay_order_id" binding:"required"`
		PaymentID string `json:"razorpay_payment_id" binding:"required"`
	}

	err := c.BindJSON(&reqDetails)
	if err != nil {
		utils.ErrorJSON(c, 400, utils.InvalidError, utils.RequestBodyError, err)
		return
	}

	godotenv.Load()
	secret := os.Getenv("RAZOR_SECRET")
	if secret == "" {
		utils.ErrorJSON(c, 500, utils.InternalError, "error getting env vars", nil)
		return
	}

	verifyParams := map[string]interface{}{
		"razorpay_order_id":   reqDetails.OrderID,
		"razorpay_payment_id": reqDetails.PaymentID,
	}

	if !razorpay_utils.VerifyPaymentSignature(verifyParams, reqDetails.Signature, secret) {
		utils.ErrorJSON(c, 400, "Payment Failed", "error verifying payment signature", nil)
		return
	}

	c.IndentedJSON(200, utils.MessageObj("Payment Success"))
}
