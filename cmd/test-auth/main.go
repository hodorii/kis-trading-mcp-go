package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"kis-trading-mcp-go/internal/kis"
)

func main() {
	// .env 로드
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Note: .env file not found, using system environment variables")
	}

	appKey := os.Getenv("KIS_APP_KEY")
	appSecret := os.Getenv("KIS_APP_SECRET")
	isPaper := os.Getenv("KIS_PAPER_TRADING") == "true"

	baseURL := "https://openapi.koreainvestment.com:9443"
	modeName := "REAL TRADING"
	if isPaper {
		baseURL = "https://openapivts.koreainvestment.com:29443"
		modeName = "PAPER TRADING"
	}

	fmt.Printf("🚀 Testing in %s mode...\n", modeName)

	if appKey == "" || appSecret == "" {
		log.Fatal("KIS_APP_KEY or KIS_APP_SECRET is not set")
	}

	// 클라이언트 생성 및 토큰 테스트
	client := kis.NewClient(appKey, appSecret, baseURL)
	err := client.FetchToken()
	if err != nil {
		log.Fatalf("❌ Token Fetch Failed: %v", err)
	}

	fmt.Println("✅ Token Fetch Success!")
	if len(client.Token) > 20 {
		fmt.Printf("🔑 Access Token (Prefix): %s...\n", client.Token[:20])
	} else {
		fmt.Println("🔑 Access Token: (too short to show prefix)")
	}
}
