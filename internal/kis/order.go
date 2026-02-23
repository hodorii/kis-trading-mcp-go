package kis

import (
	"fmt"
)

// OrderRequest 주문 요청 구조체
type OrderRequest struct {
	CANO         string `json:"CANO"`           // 계좌번호 앞 8자리
	ACNT_PRDT_CD string `json:"ACNT_PRDT_CD"`   // 계좌상품코드 (보통 01)
	PDNO         string `json:"PDNO"`           // 종목번호
	ORD_DVSN     string `json:"ORD_DVSN"`       // 주문구분 (00: 지정가)
	ORD_QTY      string `json:"ORD_QTY"`        // 주문수량
	ORD_UNPR     string `json:"ORD_UNPR"`       // 주문단가
}

// OrderResponse 주문 응답 구조체
type OrderResponse struct {
	Output struct {
		KRX_FWDG_ORD_ORG_NO string `json:"KRX_FWDG_ORD_ORG_NO"` // 한국거래소전송주문조직번호
		ODNO                string `json:"ODNO"`                // 주문번호
		ORD_TMD             string `json:"ORD_TMD"`             // 주문시각
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	MsgCd string `json:"msg_cd"`
	Msg1  string `json:"msg1"`
}

// OrderCash 현금 주문 (매수/매도)
func (c *Client) OrderCash(trID string, req OrderRequest) (*OrderResponse, error) {
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
			"tr_id":         trID, // TTTC0802U: 매수, TTTC0801U: 매도
		}).
		SetBody(req).
		SetResult(&OrderResponse{}).
		Post(c.BaseURL + "/uapi/domestic-stock/v1/trading/order-cash")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API 에러: %s", resp.String())
	}

	return resp.Result().(*OrderResponse), nil
}
