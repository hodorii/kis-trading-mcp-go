package kis

import "fmt"

// BalanceResponse 국내 주식 잔고 응답
type BalanceResponse struct {
	Output1 []struct {
		Pdno         string `json:"pdno"`           // 종목번호
		PrdtName     string `json:"prdt_name"`      // 종목명
		HldgQty      string `json:"hldg_qty"`       // 보유수량
		PchsAvgPric  string `json:"pchs_avg_pric"`  // 매입평균가
		EvluPrl1Amt  string `json:"evlu_prl1_amt"`  // 평가손익
	} `json:"output1"`
	Output2 []struct {
		DncaTotAmt string `json:"dnca_tot_amt"` // 예수금총액
	} `json:"output2"`
	RtCd string `json:"rt_cd"`
	Msg1 string `json:"msg1"`
}

// InquireBalance 국내 주식 잔고 조회
func (c *Client) InquireBalance(cano string) (*BalanceResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	trID := "TTTC8434R"
	if c.BaseURL == "https://openapivts.koreainvestment.com:29443" {
		trID = "VTTC8434R"
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
			"CANO":         cano,
			"ACNT_PRDT_CD": "01",
			"AFHR_FLPR_YN": "N",
			"OFL_YN":       "",
			"INQR_DVSN":    "01", // 01: 단가, 02: 종목별
			"UNPR_DVSN":    "01",
			"FUND_STTL_ICLD_YN": "N",
			"FNCG_AMT_AUTO_RDPT_YN": "N",
			"PRCS_DVSN":    "00",
			"CTX_AREA_FK100": "",
			"CTX_AREA_NK100": "",
		}).
		SetResult(&BalanceResponse{}).
		Get(c.BaseURL + "/uapi/domestic-stock/v1/trading/inquire-balance")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*BalanceResponse), nil
}
