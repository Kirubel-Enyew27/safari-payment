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
	"github.com/avast/retry-go"
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

	safariBaseURL := os.Getenv("SAFARI_BASE_URL")
	if safariBaseURL == "" {
		err := errors.ErrInternalServerError.New("failed to read base url")
		p.logger.Error("missing SAFARI_BASE_URL in environment", zap.String("safari_base_url", safariBaseURL))
		return dto.AcceptPaymentResponse{}, err
	}

	url := safariBaseURL + "/mpesa/stkpush/v3/processrequest"

	var (
		resp     *http.Response
		stkResp  dto.AcceptPaymentResponse
		respBody []byte
	)

	err = retry.Do(
		func() error {
			req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
			if err != nil {
				p.logger.Error("failed to create request", zap.Error(err))
				return retry.Unrecoverable(errors.ErrCreateRequest.Wrap(err, "failed to create request"))
			}

			req.Header.Set("Authorization", "Bearer "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				p.logger.Warn("temporary failure sending request to SafariCom", zap.Error(err))
				return err
			}

			// Retry on 5xx errors
			if resp.StatusCode >= 500 {
				p.logger.Warn("received 5xx from SafariCom", zap.Int("status_code", resp.StatusCode))
				return errors.ErrExternalService.New("server error: retrying")
			}

			// Don't retry on client errors
			if resp.StatusCode >= 400 {
				p.logger.Error("received 4xx from SafariCom", zap.Int("status_code", resp.StatusCode))
				return retry.Unrecoverable(errors.ErrCreateRequest.New("client error from SafariCom"))
			}

			respBody, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				p.logger.Error("failed to read response body", zap.Error(err))
				return errors.ErrInternalServerError.Wrap(err, "failed to read body")
			}

			err = json.Unmarshal(respBody, &stkResp)
			if err != nil {
				p.logger.Error("failed to unmarshal SafariCom response", zap.Error(err), zap.ByteString("body", respBody))
				return errors.ErrInternalServerError.Wrap(err, "unmarshal failed")
			}

			return nil
		},
		retry.Attempts(3),
		retry.DelayType(retry.BackOffDelay),
		retry.Delay(1*time.Second),
		retry.OnRetry(func(n uint, err error) {
			p.logger.Warn("retrying SafariCom request", zap.Uint("attempt", n+1), zap.Error(err))
		}),
	)

	if err != nil {
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

func (p *payment) GetPayments(ctx context.Context, limit int32, offset int32) ([]dto.Payment, error) {
	return p.storage.GetPayments(ctx, limit, offset)
}
