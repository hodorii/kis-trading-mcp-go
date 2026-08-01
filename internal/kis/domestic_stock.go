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

// ETFConstituentResponse ETF 구성종목 응답
type ETFConstituentResponse struct {
	Output []struct {
		StckIscd           string `json:"stck_iscd"`            // 종목코드
		IscdName           string `json:"iscd_name"`            // 종목명
		EtfCnfgIssuAvls    string `json:"etf_cnfg_issu_avls"`   // 시가총액
		EtfCnfgIssuRlim    string `json:"etf_cnfg_issu_rlim"`   // 비중(%)
		EtfVltnAmt         string `json:"etf_vltn_amt"`         // 평가금액
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	Msg1  string `json:"msg1"`
}

// InquireETFConstituents ETF 구성종목 시세 조회
func (c *Client) InquireETFConstituents(iscd string) (*ETFConstituentResponse, error) {
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
			"tr_id":         "FHKST121600C0",
		}).
		SetQueryParams(map[string]string{
			"FID_COND_MRKT_DIV_CODE": "J",
			"FID_INPUT_ISCD":         iscd,
		}).
		SetResult(&ETFConstituentResponse{}).
		Get(c.BaseURL + "/uapi/etfetn/v1/quotations/inquire-component-stock-price")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*ETFConstituentResponse), nil
}

// TimeOHLCVResponse 분봉 응답
type TimeOHLCVResponse struct {
	Output1 struct {
		PrdyVrss       string `json:"prdy_vrss"`
		PrdyVrssSign   string `json:"prdy_vrss_sign"`
		PrdyCtrt       string `json:"prdy_ctrt"`
		StckPrdyClpr   string `json:"stck_prdy_clpr"`
		AcmlVol        string `json:"acml_vol"`
		AcmlTrPbmn     string `json:"acml_tr_pbmn"`
		HtsKorIsnm     string `json:"hts_kor_isnm"`
		StckPrpr       string `json:"stck_prpr"`
	} `json:"output1"`
	Output2 []struct {
		StckBsopDate string `json:"stck_bsop_date"`
		StckCntgHour string `json:"stck_cntg_hour"`
		StckPrpr     string `json:"stck_prpr"`
		StckOprc     string `json:"stck_oprc"`
		StckHgpr     string `json:"stck_hgpr"`
		StckLwpr     string `json:"stck_lwpr"`
		CntgVol      string `json:"cntg_vol"`
		AcmlTrPbmn   string `json:"acml_tr_pbmn"`
	} `json:"output2"`
	RtCd  string `json:"rt_cd"`
	Msg1  string `json:"msg1"`
}

// InquireTimeDailyChartPrice 주식 분봉 조회
func (c *Client) InquireTimeDailyChartPrice(iscd, date, hour string) (*TimeOHLCVResponse, error) {
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
			"tr_id":         "FHKST03010230",
		}).
		SetQueryParams(map[string]string{
			"FID_COND_MRKT_DIV_CODE": "J",
			"FID_INPUT_ISCD":         iscd,
			"FID_INPUT_DATE_1":       date,
			"FID_INPUT_HOUR_1":       hour,
			"FID_PW_DATA_INCU_YN":    "Y",
			"FID_FAKE_TICK_INCU_YN":  "",
		}).
		SetResult(&TimeOHLCVResponse{}).
		Get(c.BaseURL + "/uapi/domestic-stock/v1/quotations/inquire-time-dailychartprice")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*TimeOHLCVResponse), nil
}

// OHLCVResponse 일봉/주봉/월봉 응답
type OHLCVResponse struct {
	Output []struct {
		StckBsopDate string `json:"stck_bsop_date"` // 영업일자
		StckOprc     string `json:"stck_oprc"`      // 시가
		StckHgpr     string `json:"stck_hgpr"`      // 고가
		StckLwpr     string `json:"stck_lwpr"`      // 저가
		StckClpr     string `json:"stck_clpr"`      // 종가
		AcmlVol      string `json:"acml_vol"`       // 누적거래량
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	Msg1  string `json:"msg1"`
}

// InquireDailyPrice 주식 기간별 시세 조회 (OHLCV)
func (c *Client) InquireDailyPrice(iscd, periodDiv string) (*OHLCVResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	// periodDiv: D(일봉), W(주봉), M(월봉)
	trID := "FHKST03010100" // 기본 일봉

	resp, err := c.httpClient.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
			"appkey":        c.AppKey,
			"appsecret":     c.AppSecret,
			"tr_id":         trID,
		}).
		SetQueryParams(map[string]string{
			"FID_COND_MRKT_DIV_CODE": "J",
			"FID_INPUT_ISCD":         iscd,
			"FID_PERIOD_DIV_CODE":    periodDiv,
		}).
		SetResult(&OHLCVResponse{}).
		Get(c.BaseURL + "/uapi/domestic-stock/v1/quotations/inquire-daily-itemchartprice")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*OHLCVResponse), nil
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
