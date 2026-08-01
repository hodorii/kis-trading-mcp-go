package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"kis-trading-mcp-go/internal/kis"
)

func main() {
	if err := godotenv.Load("../.env.prd"); err != nil {
		_ = godotenv.Load(".env")
	}

	appKey := os.Getenv("KIS_APP_KEY")
	appSecret := os.Getenv("KIS_APP_SECRET")
	isPaper := os.Getenv("KIS_SIMULATION") == "true"
	cano := os.Getenv("KIS_ACCT_STOCK")

	if appKey == "" || appSecret == "" {
		log.Fatal("KIS_APP_KEY or KIS_APP_SECRET not set")
	}
	if cano == "" {
		log.Fatal("KIS_ACCT_STOCK not set")
	}

	baseURL := "https://openapi.koreainvestment.com:9443"
	if isPaper {
		baseURL = "https://openapivts.koreainvestment.com:29443"
		fmt.Println("🚀 Mode: PAPER TRADING")
	} else {
		fmt.Println("⚠️  Mode: REAL TRADING")
	}

	client := kis.NewClient(appKey, appSecret, baseURL)

	fmt.Printf("📊 Testing account balance inquiry...\n")
	fmt.Printf("   Account: %s\n\n", cano)

	resp, err := client.InquireBalance(cano)
	if err != nil {
		log.Fatalf("❌ Balance inquiry failed: %v", err)
	}

	if resp.RtCd != "" && resp.RtCd != "0" {
		log.Fatalf("❌ API Error: rt_cd=%s, msg=%s", resp.RtCd, resp.Msg1)
	}

	data, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Printf("✅ Account Balance Result:\n%s\n", string(data))
}
