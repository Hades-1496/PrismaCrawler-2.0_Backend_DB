package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"prismacrawler/internal/models"
	"prismacrawler/pkg/db"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

// Helper to compute stripe signature
func computeSignature(t time.Time, payload []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(fmt.Sprintf("%d.", t.Unix())))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestStripeWebhook(t *testing.T) {
	_ = godotenv.Load("../../.env.test")
	testDBUrl := os.Getenv("DIRECT_URL")
	if testDBUrl == "" {
		testDBUrl = os.Getenv("DATABASE_URL")
	}
	db.ConnectDB(testDBUrl)
	db.DB.Exec("DROP TABLE IF EXISTS processed_stripe_events, wallets CASCADE;")
	db.DB.AutoMigrate(&models.ProcessedStripeEvent{}, &models.Wallet{})

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/webhook", StripeWebhook)

	secret := "whsec_test_secret"
	os.Setenv("STRIPE_WEBHOOK_SECRET", secret)

	userID := uint(1)
	wallet := models.Wallet{UserID: userID, Gems: 0}
	db.DB.Create(&wallet)

	payload := []byte(`{"id":"evt_test123","type":"checkout.session.completed","api_version":"2024-04-10","data":{"object":{"client_reference_id":"1","metadata":{"amount":"500","package_id":"gems_500"},"payment_status":"paid"}}}`)

	now := time.Now()
	sig := computeSignature(now, payload, secret)
	header := fmt.Sprintf("t=%d,v1=%s", now.Unix(), sig)

	req, _ := http.NewRequest("POST", "/webhook", bytes.NewBuffer(payload))
	req.Header.Set("Stripe-Signature", header)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedWallet models.Wallet
	db.DB.First(&updatedWallet, "user_id = ?", userID)
	assert.Equal(t, 500, updatedWallet.Gems)

	var event models.ProcessedStripeEvent
	err := db.DB.First(&event, "event_id = ?", "evt_test123").Error
	assert.NoError(t, err)
}
