package dto

import (
	"time"

	"github.com/google/uuid"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}
type AccessTokenResponse struct {
	TokenResponse
	IssuedAt time.Time
}

type AcceptPaymentRequest struct {
	MerchantRequestID string `json:"MerchantRequestID"`
	BusinessShortCode string `json:"BusinessShortCode"`
	Password          string `json:"Password"`
	Timestamp         string `json:"Timestamp"`
	TransactionType   string `json:"TransactionType"`
	Amount            int    `json:"Amount"`
	PartyA            string `json:"PartyA"`
	PartyB            string `json:"PartyB"`
	PhoneNumber       string `json:"PhoneNumber"`
	CallBackURL       string `json:"CallBackURL"`
	AccountReference  string `json:"AccountReference"`
	TransactionDesc   string `json:"TransactionDesc"`
}

type AcceptPaymentResponse struct {
	MerchantRequestID string `json:"MerchantRequestID"`
	CheckoutRequestID string `json:"CheckoutRequestID"`
	ResponseCode      string `json:"ResponseCode"`
	ResponseDesc      string `json:"ResponseDescription"`
	CustomerMessage   string `json:"CustomerMessage"`
}

type Payment struct {
	ID                uuid.UUID `db:"id"`
	CheckoutRequestID string    `db:"checkout_request_id"`
	MerchantRequestID string    `db:"merchant_request_id"`
	PhoneNumber       string    `db:"phone_number"`
	Amount            float64   `db:"amount"`
	MpesaReceipt      string    `db:"mpesa_receipt"`
	TransactionDate   time.Time `db:"transaction_date"`
	ResultCode        int       `db:"result_code"`
	ResultDesc        string    `db:"result_desc"`
	CreatedAt         time.Time `db:"created_at"`
}
