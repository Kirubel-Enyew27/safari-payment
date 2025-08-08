package dto

import (
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
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

func (acp *AcceptPaymentRequest) Validate() error {
	err := validation.ValidateStruct(acp,
		validation.Field(&acp.MerchantRequestID, validation.Required),
		validation.Field(&acp.BusinessShortCode,
			validation.Required,
			validation.Length(5, 6),
			is.Digit),
		validation.Field(&acp.Password, validation.Required),
		validation.Field(&acp.Timestamp,
			validation.Required,
			is.Digit),
		validation.Field(&acp.TransactionType,
			validation.Required,
			validation.In("CustomerPayBillOnline", "CustomerBuyGoodsOnline")),
		validation.Field(&acp.Amount,
			validation.Required,
			validation.Min(1)),
		validation.Field(&acp.PartyA,
			validation.Required,
			validation.Length(12, 12),
			validation.Match(regexp.MustCompile(`^2517\d{8}$`)).Error("must be a valid 12-digit Safaricom phone number starting with 2517")),
		validation.Field(&acp.PartyB,
			validation.Required,
			validation.Length(5, 6),
			is.Digit),
		validation.Field(&acp.PhoneNumber,
			validation.Required,
			validation.Length(12, 12),
			validation.Match(regexp.MustCompile(`^2517\d{8}$`)).Error("must be a valid 12-digit Safaricom phone number starting with 2517")),
		validation.Field(&acp.CallBackURL,
			validation.Required,
			is.URL),
		validation.Field(&acp.AccountReference,
			validation.Required,
			validation.Length(1, 12)),
		validation.Field(&acp.TransactionDesc,
			validation.Required,
			validation.Length(1, 13)),
	)

	return err
}
