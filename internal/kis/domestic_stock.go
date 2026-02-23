package kis

import (
	"fmt"
)

// PriceResponse 주식 현재가 응답 구조체 (필요한 필드 위주)
type PriceResponse struct {
	Output struct {
		IscdStatClrCode string `json:"iscd_stat_clr_code"`
		MargRate        string `json:"marg_rate"`
		StckPrpr        string `json:"stck_prpr"` // 현재가
		PrdyVrss        string `json:"prdy_vrss"` // 전일대비
		PrdyVrssSign    string `json:"prdy_vrss_sign"`
		PrdyCtrt        string `json:"prdy_ctrt"` // 전일대비율
		StckRefp        string `json:"stck_refp"` // 기준가
		StckOprc        string `json:"stck_oprc"` // 시가
		StckHgpr        string `json:"stck_hgpr"` // 고가
		StckLwpr        string `json:"stck_lwpr"` // 저가
		StckSdpr        string `json:"stck_sdpr"` // 상한가
		StckLdpr        string `json:"stck_ldpr"` // 하한가
		AcmlVol         string `json:"acml_vol"`  // 누적 거래량
		AcmlTrPbmn      string `json:"acml_tr_pbmn"` // 누적 거래대금
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	MsgCd string `json:"msg_cd"`
	Msg1  string `json:"msg1"`
}

// InquirePrice 주식 현재가 조회
func (c *Client) InquirePrice(marketDiv, iscd string) (*PriceResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
			"appkey":        c.AppKey,
			"appsecret":     c.AppSecret,
			"tr_id":         "FHKST01010100",
		}).
		SetQueryParams(map[string]string{
			"FID_COND_MRKT_DIV_CODE": marketDiv,
			"FID_INPUT_ISCD":         iscd,
		}).
		SetResult(&PriceResponse{}).
		Get(c.BaseURL + "/uapi/domestic-stock/v1/quotations/inquire-price")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API 에러: %s", resp.String())
	}

	return resp.Result().(*PriceResponse), nil
}
