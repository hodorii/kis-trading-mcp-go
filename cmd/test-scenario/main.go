package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"kis-trading-mcp-go/internal/db"
	"kis-trading-mcp-go/internal/kis"
)

func main() {
	// .env 로드
	_ = godotenv.Load(".env")

	appKey := os.Getenv("KIS_APP_KEY")
	appSecret := os.Getenv("KIS_APP_SECRET")
	isPaper := os.Getenv("KIS_PAPER_TRADING") == "true"
	cano := os.Getenv("KIS_ACCT_STOCK") // 계좌번호

	baseURL := "https://openapi.koreainvestment.com:9443"
	if isPaper {
		baseURL = "https://openapivts.koreainvestment.com:29443"
		fmt.Println("🚀 Mode: PAPER TRADING")
	} else {
		fmt.Println("⚠️ Mode: REAL TRADING")
	}

	// 1. 초기화
	client := kis.NewClient(appKey, appSecret, baseURL)
	dbPath, _ := filepath.Abs("configs/master/master.db")
	repo, err := db.NewRepository(dbPath)
	if err != nil {
		log.Fatalf("DB Open Error: %v", err)
	}
	defer repo.Close()

	// [시나리오 1단계: 종목 조회]
	stockName := "삼성전자"
	stocks, err := repo.FindStock(stockName)
	if err != nil || len(stocks) == 0 {
		log.Fatalf("Stock not found: %s", stockName)
	}
	iscd := stocks[0].Code
	fmt.Printf("✅ [Step 1] Stock Found: %s (%s)\n", stockName, iscd)

	// [시나리오 2단계: 현재가 확인]
	priceRes, err := client.InquirePrice("J", iscd)
	if err != nil {
		log.Fatalf("Price lookup failed: %v", err)
	}
	fmt.Printf("✅ [Step 2] Current Price: %s KRW\n", priceRes.Output.StckPrpr)

	// [시나리오 3단계: 주문 준비]
	if cano == "" {
		fmt.Println("ℹ️ Skip: Account (cano) not set")
		return
	}

	fmt.Printf("✅ [Step 3] Ready to order for account: %s\n", cano)
	fmt.Println("🎉 Scenario Test Success!")
}
