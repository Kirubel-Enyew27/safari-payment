package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Kirubel-Enyew27/safari-payment/internal/model/dto"
)

func GetSafariAccessToken() (dto.AccessTokenResponse, error) {
	consumerKey := os.Getenv("SAFARI_CONSUMER_KEY")
	consumerSecret := os.Getenv("SAFARI_CONSUMER_SECRET")
	safari_base_url := os.Getenv("SAFARI_BASE_URL")
	endpoint := safari_base_url + "/v1/token/generate?grant_type=client_credentials"

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return dto.AccessTokenResponse{}, err
	}
	req.SetBasicAuth(consumerKey, consumerSecret)

	resp, err := client.Do(req)
	if err != nil {
		return dto.AccessTokenResponse{}, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return dto.AccessTokenResponse{}, fmt.Errorf("failed to get token: %s", body)
	}

	var tokenResp dto.TokenResponse
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return dto.AccessTokenResponse{}, err
	}

	token := dto.AccessTokenResponse{
		TokenResponse: tokenResp,
		IssuedAt:      time.Now(),
	}

	return token, nil
}
