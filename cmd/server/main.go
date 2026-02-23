package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	s := server.NewMCPServer("KIS Trading MCP Server (Go)", "1.0.0", server.WithLogging())

	// [1] 국내 주식 도구 (잔고 포함)
	s.AddTool(mcp.NewTool("domestic_stock",
		mcp.WithDescription("국내주식 조회/주문/뉴스/공시/잔고 도구"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		apiType, _ := args["api_type"].(string)
		params, _ := args["params"].(map[string]any)

		switch apiType {
		case "find_stock_code":
			keyword, _ := params["stock_name"].(string)
			res, _ := repo.FindStock(keyword)
			data, _ := json.MarshalIndent(res, "", "  ")
			return mcp.NewToolResultText(string(data)), nil

		case "inquire_price":
			iscd, _ := params["fid_input_iscd"].(string)
			res, err := kisClient.InquirePrice("J", iscd)
			if err != nil { return mcp.NewToolResultError(err.Error()), nil }
			data, _ := json.MarshalIndent(res.Output, "", "  ")
			return mcp.NewToolResultText(string(data)), nil

		case "inquire_balance":
			res, err := kisClient.InquireBalance(cano)
			if err != nil { return mcp.NewToolResultError(err.Error()), nil }
			data, _ := json.MarshalIndent(res, "", "  ")
			return mcp.NewToolResultText(string(data)), nil

		case "news_title":
			iscd, _ := params["fid_input_iscd"].(string)
			res, err := kisClient.InquireNewsTitle(iscd)
			if err != nil { return mcp.NewToolResultError(err.Error()), nil }
			data, _ := json.MarshalIndent(res.Output, "", "  ")
			return mcp.NewToolResultText(string(data)), nil

		case "order_cash":
			trID := "TTTC0802U"
			if isPaper { trID = "VTTC0802U" }
			orderReq := kis.OrderRequest{
				CANO: cano,
				ACNT_PRDT_CD: "01",
				PDNO: params["fid_input_iscd"].(string),
				ORD_DVSN: "01",
				ORD_QTY: fmt.Sprintf("%v", params["ord_qty"]),
				ORD_UNPR: "0",
			}
			res, err := kisClient.OrderCash(trID, orderReq)
			if err != nil { return mcp.NewToolResultError(err.Error()), nil }
			data, _ := json.MarshalIndent(res.Output, "", "  ")
			return mcp.NewToolResultText(string(data)), nil

		default:
			return mcp.NewToolResultError("Unsupported api_type"), nil
		}
	})

	// [2] 해외 주식 도구 (미국 등)
	s.AddTool(mcp.NewTool("overseas_stock",
		mcp.WithDescription("해외주식(미국 등) 조회 도구"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, _ := request.Params.Arguments.(map[string]any)
		apiType, _ := args["api_type"].(string)
		params, _ := args["params"].(map[string]any)

		if apiType == "inquire_price" {
			res, err := kisClient.InquireOverseasPrice(params["excd"].(string), params["symb"].(string))
			if err != nil { return mcp.NewToolResultError(err.Error()), nil }
			data, _ := json.MarshalIndent(res.Output, "", "  ")
			return mcp.NewToolResultText(string(data)), nil
		}
		return mcp.NewToolResultError("Unsupported api_type"), nil
	})

	log.Printf("🚀 KIS MCP Server (Go) with Balance starting...")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
