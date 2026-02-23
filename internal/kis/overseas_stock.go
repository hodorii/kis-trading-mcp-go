package kis

import "fmt"

// OverseasPriceResponse 해외 주식 현재가 응답
type OverseasPriceResponse struct {
	Output struct {
		Rsym string `json:"rsym"` // 종목코드
		Last string `json:"last"` // 현재가
		Diff string `json:"diff"` // 대비
		Rate string `json:"rate"` // 등락률
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	Msg1  string `json:"msg1"`
}

// InquireOverseasPrice 해외 주식 현재가 조회
func (c *Client) InquireOverseasPrice(exch, iscd string) (*OverseasPriceResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	trID := "HHDFS00000300"
	if c.BaseURL == "https://openapivts.koreainvestment.com:29443" {
		trID = "VHTS3001"
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
			"AUTH": "",
			"EXCD": exch,
			"SYMB": iscd,
		}).
		SetResult(&OverseasPriceResponse{}).
		Get(c.BaseURL + "/uapi/overseas-stock/v1/quotations/price")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*OverseasPriceResponse), nil
}
