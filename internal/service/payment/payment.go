package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/errors"
	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
	"github.com/Kirubel-Enyew27/safari-payment/internal/service"
	"github.com/Kirubel-Enyew27/safari-payment/internal/storage"
	"github.com/Kirubel-Enyew27/safari-payment/utils"
	"go.uber.org/zap"
)

type payment struct {
	storage storage.Payment
	token   dto.AccessTokenResponse
	logger  *zap.Logger
}

func InitService(paymentStorage storage.Payment, token dto.AccessTokenResponse, log *zap.Logger) service.Payment {
	return &payment{
		storage: paymentStorage,
		token:   token,
		logger:  log,
	}

}

func (p *payment) AcceptPayment(ctx context.Context, payload dto.AcceptPaymentRequest) (dto.AcceptPaymentResponse, error) {
	shortcode := os.Getenv("SAFARI_BUSINESS_SHORT_CODE")
	password := os.Getenv("SAFARI_PASSWORD")
	if shortcode == "" || password == "" {
		err := errors.ErrInternalServerError.New("fialed to read business short code or password")
		p.logger.Error("failed to read business short code or password from config", zap.String("short code", shortcode), zap.String("password", password))
		return dto.AcceptPaymentResponse{}, err
	}
	timestamp := time.Now().Format("20060102150405")

	payload.BusinessShortCode = shortcode
	payload.Password = password
	payload.Timestamp = timestamp
	payload.PartyB = shortcode

	tokenExpirationDuration, err := time.ParseDuration(p.token.ExpiresIn + "s")
	if err != nil {
		err := errors.ErrUnExpectedError.Wrap(err, "failed to parse token expiration duration")
		p.logger.Error("Failed to parse token expiration",
			zap.String("expires_in", p.token.ExpiresIn),
			zap.Error(err),
		)
		return dto.AcceptPaymentResponse{}, err
	}

	// Add (5 seconds) to prevent using expired token
	expired := time.Since(p.token.IssuedAt.Add(5*time.Second)) > tokenExpirationDuration
	if expired || p.token.AccessToken == "" {
		newToken, err := utils.GetSafariAccessToken()
		if err != nil {
			wrappedErr := errors.ErrAccessToken.Wrap(err, "failed to refresh access token")
			p.logger.Error("Failed to refresh access token", zap.Error(wrappedErr))
			return dto.AcceptPaymentResponse{}, wrappedErr
		}
		p.token = newToken
	}

	resp, err := p.processRequest(ctx, payload, p.token.AccessToken)
	if err != nil {
		return dto.AcceptPaymentResponse{}, err
	}

	return resp, nil

}

func (p *payment) processRequest(ctx context.Context, payload dto.AcceptPaymentRequest, token string) (dto.AcceptPaymentResponse, error) {
	err := payload.Validate()
	if err != nil {
		err := errors.ErrCreateRequest.Wrap(err, "validation failed")
		p.logger.Error("validation failed", zap.Error(err))
		return dto.AcceptPaymentResponse{}, err
	}

	jsonPayload, _ := json.Marshal(payload)
	safari_base_url := os.Getenv("SAFARI_BASE_URL")
	if safari_base_url == "" {
		err := errors.ErrInternalServerError.New("fialed to read base url")
		p.logger.Error("failed to read safari base url from config", zap.String("safari_base_url", safari_base_url))
		return dto.AcceptPaymentResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", safari_base_url+"/mpesa/stkpush/v3/processrequest", bytes.NewBuffer(jsonPayload))
	if err != nil {
		err := errors.ErrCreateRequest.Wrap(err, "fialed to create request")
		p.logger.Error("fialed to create request to safari", zap.Error(err))
		return dto.AcceptPaymentResponse{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err := errors.ErrCreateRequest.Wrap(err, "fialed to send request")
		p.logger.Error("fialed to send request to safari", zap.Error(err))
		return dto.AcceptPaymentResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var stkResp dto.AcceptPaymentResponse
	err = json.Unmarshal(body, &stkResp)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "failed to unmarshal response")
		p.logger.Error("failed to unmarshal response", zap.Error(err))
		return dto.AcceptPaymentResponse{}, err
	}

	return stkResp, nil
}

func (p *payment) StorePayment(ctx context.Context, payload map[string]any) (dto.Payment, error) {
	body, ok := payload["Body"].(map[string]any)
	if !ok {
		err := errors.ErrFailedPayment.New("invalid structure")
		p.logger.Error("invalid structure", zap.Any("body", body))
		return dto.Payment{}, err
	}

	stkCallback := body["stkCallback"].(map[string]any)
	resultCode := int(stkCallback["ResultCode"].(float64))
	resultDesc := stkCallback["ResultDesc"].(string)

	if resultCode != 0 {
		err := errors.ErrFailedPayment.New("payment failed")
		p.logger.Error("payment failed", zap.Int("resultCode", resultCode), zap.String("resultDesc", resultDesc))
		return dto.Payment{}, err
	}

	callbackMetadata := stkCallback["CallbackMetadata"].(map[string]any)["Item"].([]any)
	data := map[string]any{}
	for _, item := range callbackMetadata {
		entry := item.(map[string]any)
		name := entry["Name"].(string)
		value := entry["Value"]
		data[name] = value
	}

	payment := dto.Payment{
		CheckoutRequestID: stkCallback["CheckoutRequestID"].(string),
		MerchantRequestID: stkCallback["MerchantRequestID"].(string),
		PhoneNumber:       data["PhoneNumber"].(string),
		MpesaReceipt:      data["MpesaReceiptNumber"].(string),
		Amount:            data["Amount"].(float64),
		TransactionDate:   utils.ParseMpesaDate(data["TransactionDate"]),
		ResultCode:        resultCode,
		ResultDesc:        resultDesc,
	}

	return p.storage.SavePayment(ctx, payment)

}
