package kis

import "fmt"

// FutureOptionPriceResponse 국내 선물옵션 현재가 응답
type FutureOptionPriceResponse struct {
	Output struct {
		StckPrpr string `json:"stck_prpr"` // 현재가
		PrdyVrss string `json:"prdy_vrss"` // 전일대비
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	Msg1  string `json:"msg1"`
}

// InquireFutureOptionPrice 국내 선물옵션 현재가 조회
func (c *Client) InquireFutureOptionPrice(iscd string) (*FutureOptionPriceResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	trID := "FHKIF01010100"
	if c.BaseURL == "https://openapivts.koreainvestment.com:29443" {
		trID = "VHKIF01010100"
	}

	resp, err := c.httpClient.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
			"appkey":        c.AppKey,
			"appsecret":     c.AppSecret,
			"tr_id":         trID,
		}).
		SetQueryParams(map[string]string{
			"FID_COND_MRKT_DIV_CODE": "U",
			"FID_INPUT_ISCD":         iscd,
		}).
		SetResult(&FutureOptionPriceResponse{}).
		Get(c.BaseURL + "/uapi/domestic-futureoption/v1/quotations/inquire-price")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*FutureOptionPriceResponse), nil
}
