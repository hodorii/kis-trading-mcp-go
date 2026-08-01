package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"kis-trading-mcp-go/internal/db"
	"kis-trading-mcp-go/internal/kis"
)

func main() {
	_ = godotenv.Load(".env")

	appKey := os.Getenv("KIS_APP_KEY")
	appSecret := os.Getenv("KIS_APP_SECRET")
	isPaper := os.Getenv("KIS_PAPER_TRADING") == "true"
	cano := os.Getenv("KIS_ACCT_STOCK")

	baseURL := "https://openapi.koreainvestment.com:9443"
	if isPaper {
		baseURL = "https://openapivts.koreainvestment.com:29443"
	}

	kisClient := kis.NewClient(appKey, appSecret, baseURL)
	dbPath, _ := filepath.Abs("configs/master/master.db")
	repo, _ := db.NewRepository(dbPath)

	s := server.NewMCPServer("KIS Trading MCP Server (Korea Investment & Securities)", "1.0.0", server.WithLogging())

	// --- 국내 주식 도구 ---
	s.AddTool(mcp.NewTool("find_stock_code",
		mcp.WithDescription("종목명으로 주식 코드 검색"),
		mcp.WithString("stock_name", mcp.Required(), mcp.Description("검색할 종목명")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments.(map[string]any)
		name := args["stock_name"].(string)
		res, _ := repo.FindStock(name)
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	s.AddTool(mcp.NewTool("get_domestic_price",
		mcp.WithDescription("국내 주식 현재가 조회"),
		mcp.WithString("iscd", mcp.Required(), mcp.Description("종목 코드 (6자리)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments.(map[string]any)
		iscd := args["iscd"].(string)
		res, err := kisClient.InquirePrice("J", iscd)
		if err != nil { return mcp.NewToolResultError(err.Error()), nil }
		data, _ := json.MarshalIndent(res.Output, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	s.AddTool(mcp.NewTool("get_balance",
		mcp.WithDescription("국내 주식 계좌 잔고 조회"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := kisClient.InquireBalance(cano)
		if err != nil { return mcp.NewToolResultError(err.Error()), nil }
		data, _ := json.MarshalIndent(res, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	s.AddTool(mcp.NewTool("get_ohlcv",
		mcp.WithDescription("국내 주식 기간별 시세(OHLCV) 조회"),
		mcp.WithString("iscd", mcp.Required(), mcp.Description("종목 코드 (6자리)")),
		mcp.WithString("period", mcp.Required(), mcp.Description("조회 구분: 'D' (일봉), 'W' (주봉), 'M' (월봉)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments.(map[string]any)
		iscd := args["iscd"].(string)
		period := args["period"].(string)
		res, err := kisClient.InquireDailyPrice(iscd, period)
		if err != nil { return mcp.NewToolResultError(err.Error()), nil }
		data, _ := json.MarshalIndent(res.Output, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	s.AddTool(mcp.NewTool("get_minute_ohlcv",
		mcp.WithDescription("국내 주식 분봉(시간대별) 시세 조회 — 1회 호출당 약 120개 분봉 반환"),
		mcp.WithString("iscd", mcp.Required(), mcp.Description("종목 코드 (6자리)")),
		mcp.WithString("date", mcp.Description("조회 기준일 YYYYMMDD (생략 시 당일)")),
		mcp.WithString("hour", mcp.Description("조회 종료 시각 HHMMSS, 이 시각부터 과거로 조회 (예: 153000, 생략 시 최신)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments.(map[string]any)
		iscd := args["iscd"].(string)
		date, _ := args["date"].(string)
		hour, _ := args["hour"].(string)
		res, err := kisClient.InquireTimeDailyChartPrice(iscd, date, hour)
		if err != nil { return mcp.NewToolResultError(err.Error()), nil }
		data, _ := json.MarshalIndent(res.Output2, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})

	s.AddTool(mcp.NewTool("get_etf_constituents",
		mcp.WithDescription("ETF 구성종목 및 비중 조회"),
		mcp.WithString("iscd", mcp.Required(), mcp.Description("ETF 종목 코드 (6자리)")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments.(map[string]any)
		iscd := args["iscd"].(string)
		res, err := kisClient.InquireETFConstituents(iscd)
		if err != nil { return mcp.NewToolResultError(err.Error()), nil }
		data, _ := json.MarshalIndent(res.Output, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})


	s.AddTool(mcp.NewTool("order_domestic_stock",
		mcp.WithDescription("국내 주식 현금 주문"),
		mcp.WithString("iscd", mcp.Required(), mcp.Description("종목 코드 (6자리)")),
		mcp.WithString("side", mcp.Required(), mcp.Description("주문 구분 ('buy' 또는 'sell')")),
		mcp.WithString("qty", mcp.Required(), mcp.Description("주문 수량")),
		mcp.WithString("price", mcp.Required(), mcp.Description("주문 단가")),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.Params.Arguments.(map[string]any)
		iscd := args["iscd"].(string)
		side := args["side"].(string)
		qty := args["qty"].(string)
		price := args["price"].(string)

		trID := "TTTC0802U" // 기본 매수
		if side == "sell" {
			trID = "TTTC0801U"
		}
		if isPaper {
			trID = "V" + trID[1:]
		}

		orderReq := kis.OrderRequest{
			CANO:         cano,
			ACNT_PRDT_CD: "01",
			PDNO:         iscd,
			ORD_DVSN:     "00", // 지정가
			ORD_QTY:      qty,
			ORD_UNPR:     price,
		}
		res, err := kisClient.OrderCash(trID, orderReq)
		if err != nil { return mcp.NewToolResultError(err.Error()), nil }
		data, _ := json.MarshalIndent(res.Output, "", "  ")
		return mcp.NewToolResultText(string(data)), nil
	})


	log.Printf("🚀 KIS MCP Server (Go) with Balance starting...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
