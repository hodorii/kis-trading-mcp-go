package kis

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// TokenResponse KIS API 토큰 응답 구조체
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// Client KIS API 클라이언트
type Client struct {
	httpClient *resty.Client
	AppKey     string
	AppSecret  string
	BaseURL    string
	Token      string
}

// NewClient KIS 클라이언트 생성
func NewClient(appKey, appSecret, baseURL string) *Client {
	return &Client{
		httpClient: resty.New(),
		AppKey:     appKey,
		AppSecret:  appSecret,
		BaseURL:    baseURL,
	}
}

// FetchToken 토큰 발급 (POST /oauth2/tokenP)
func (c *Client) FetchToken() error {
	payload := map[string]string{
		"grant_type": "client_credentials",
		"appkey":     c.AppKey,
		"appsecret":  c.AppSecret,
	}

	resp, err := c.httpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post(c.BaseURL + "/oauth2/tokenP")

	if err != nil {
		return fmt.Errorf("HTTP 요청 실패: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("API 에러: %s (상태코드: %d)", resp.String(), resp.StatusCode())
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil {
		return fmt.Errorf("JSON 파싱 실패: %w", err)
	}

	c.Token = tokenResp.AccessToken
	return nil
}

// GetToken 토큰 반환 (없으면 발급)
func (c *Client) GetToken() (string, error) {
	if c.Token == "" {
		if err := c.FetchToken(); err != nil {
			return "", err
		}
	}
	return c.Token, nil
}
