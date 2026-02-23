package main

import (
	"fmt"
	"log"
	"path/filepath"

	"kis-trading-mcp-go/internal/db"
)

func main() {
	// DB 경로 (절대 경로)
	dbPath, err := filepath.Abs("configs/master/master.db")
	if err != nil {
		log.Fatalf("❌ 경로 설정 실패: %v", err)
	}

	// 저장소 생성
	repo, err := db.NewRepository(dbPath)
	if err != nil {
		log.Fatalf("❌ DB 연결 실패: %v", err)
	}
	defer repo.Close()

	// 검색 테스트 (삼성전자)
	keyword := "삼성전자"
	results, err := repo.FindStock(keyword)
	if err != nil {
		log.Fatalf("❌ 검색 실패: %v", err)
	}

	fmt.Printf("🔍 '%s' 검색 결과 (%d건):\n", keyword, len(results))
	for _, s := range results {
		fmt.Printf("- [%s] %s (거래소: %s)\n", s.Code, s.Name, s.Ex)
	}
}
