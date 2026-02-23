package kis

import "fmt"

// NewsTitleResponse 뉴스 제목 응답
type NewsTitleResponse struct {
	Output []struct {
		NewsPrepDate string `json:"news_prep_date"` // 작성일자
		NewsPrepTime string `json:"news_prep_time"` // 작성시간
		NewsTitl     string `json:"news_titl"`      // 뉴스제목
		NewsSrno     string `json:"news_srno"`      // 일련번호
	} `json:"output"`
	RtCd  string `json:"rt_cd"`
	Msg1  string `json:"msg1"`
}

// InquireNewsTitle 종합 시황/공시 제목 조회
func (c *Client) InquireNewsTitle(iscd string) (*NewsTitleResponse, error) {
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	// 뉴스/공시 조회 TR_ID: FHKST01010700 (실전 기준)
	trID := "FHKST01010700"
	
	resp, err := c.httpClient.R().
		SetHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + token,
			"appkey":        c.AppKey,
			"appsecret":     c.AppSecret,
			"tr_id":         trID,
		}).
		SetQueryParams(map[string]string{
			"fid_news_ofer_entp_code": "", // 전체
			"fid_cond_mrkt_cls_code":  "", // 전체
			"fid_input_iscd":         iscd,
			"fid_titl_cntt":          "",
			"fid_input_date_1":       "",
			"fid_input_hour_1":       "",
			"fid_rank_sort_cls_code": "",
			"fid_input_srno":         "",
		}).
		SetResult(&NewsTitleResponse{}).
		Get(c.BaseURL + "/uapi/domestic-stock/v1/quotations/news-title")

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("API error: %s", resp.String())
	}

	return resp.Result().(*NewsTitleResponse), nil
}
