package db

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// StockMaster 국내 주식 마스터 데이터 모델
type StockMaster struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Ex   string `json:"ex"`
}

// Repository 데이터베이스 저장소
type Repository struct {
	db *sql.DB
}

// NewRepository DB 연결 생성 (SQLite)
func NewRepository(dbPath string) (*Repository, error) {
	// modernc.org/sqlite 드라이버 이름은 "sqlite"
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db error: %w", err)
	}

	// 연결 확인
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db error: %w", err)
	}

	return &Repository{db: db}, nil
}

// Close DB 연결 종료
func (r *Repository) Close() error {
	return r.db.Close()
}

// FindStock 종목명 또는 코드로 검색 (LIKE)
func (r *Repository) FindStock(keyword string) ([]StockMaster, error) {
	// 국내 주식 마스터 테이블: domestic_stock_master (name, code, ex)
	query := `
		SELECT code, name, ex 
		FROM domestic_stock_master 
		WHERE name LIKE ? OR code LIKE ? 
		ORDER BY CASE WHEN name = ? THEN 0 ELSE 1 END, name
		LIMIT 10
	`
	param := "%" + keyword + "%"

	rows, err := r.db.Query(query, param, param, keyword)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var results []StockMaster
	for rows.Next() {
		var s StockMaster
		var ex sql.NullString // ex 컬럼 NULL 처리
		if err := rows.Scan(&s.Code, &s.Name, &ex); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		s.Ex = ex.String
		results = append(results, s)
	}
	return results, nil
}
