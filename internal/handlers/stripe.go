package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"prismacrawler/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/webhook"
)

// Map of package IDs to Gem amounts and prices (in cents)
var gemPackages = map[string]struct {
	Amount int
	Price  int64
}{
	"gems_500":  {Amount: 500, Price: 499},  // $4.99
	"gems_1200": {Amount: 1200, Price: 999}, // $9.99
}

// CreateStripeCheckout creates a Stripe Checkout Session for buying Gems
func CreateStripeCheckout(c *gin.Context) {
	userID := utils.GetUserID(c)
	var req CreateCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid package")
		return
	}

	pkg, exists := gemPackages[req.PackageID]
	if !exists {
		utils.SendError(c, http.StatusBadRequest, "Unknown package")
		return
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		utils.SendError(c, http.StatusInternalServerError, "Stripe not configured")
		return
	}

	// Assuming the frontend will be running on a certain origin, but fallback to localhost
	origin := c.Request.Header.Get("Origin")
	if origin == "" {
		origin = "http://localhost:3000"
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(fmt.Sprintf("%d Prisma Gems", pkg.Amount)),
					},
					UnitAmount: stripe.Int64(pkg.Price),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:        stripe.String(origin + "/shop?success=true"),
		CancelURL:         stripe.String(origin + "/shop?canceled=true"),
		ClientReferenceID: stripe.String(strconv.Itoa(int(userID))),
		Metadata: map[string]string{
			"package_id": req.PackageID,
			"amount":     strconv.Itoa(pkg.Amount),
		},
	}

	s, err := session.New(params)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to create checkout session")
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": s.URL})
}

// StripeWebhook handles asynchronous payment confirmations from Stripe
func StripeWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signatureHeader := c.GetHeader("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, signatureHeader, webhookSecret)
	if err != nil {
		fmt.Printf("⚠️  Webhook signature verification failed. %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error parsing session JSON"})
			return
		}

		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			userIDStr := session.ClientReferenceID
			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			amountStr := session.Metadata["amount"]
			amount, err := strconv.Atoi(amountStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid amount in metadata"})
				return
			}

			// Add gems to user's wallet
			var wallet models.Wallet
			db.DB.Where("user_id = ?", uint(userID)).FirstOrCreate(&wallet, models.Wallet{UserID: uint(userID)})
			wallet.Gems += amount
			db.DB.Save(&wallet)
			fmt.Printf("✅ Added %d Gems to user %d via Stripe\n", amount, userID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}
