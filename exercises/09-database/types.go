// 練習 09 共用的型別、錯誤與資料表結構（題目已提供，不需修改）。
package library

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite" // 匿名匯入：只為了執行驅動程式的 init()，把 "sqlite" 註冊到 database/sql
)

// Book 對應 books 資料表的一列。
type Book struct {
	ID     int64
	Title  string
	Author string
	Year   int
	Stock  int // 可借出的庫存數
}

// Loan 是一筆借閱紀錄。BookTitle 來自 JOIN books。
type Loan struct {
	ID         int64
	BookID     int64
	BookTitle  string
	Member     string
	BorrowedAt time.Time
	ReturnedAt *time.Time // 資料庫中可能是 NULL（尚未歸還）
}

var (
	ErrNotFound        = errors.New("找不到資料")
	ErrOutOfStock      = errors.New("庫存不足")
	ErrAlreadyReturned = errors.New("已經歸還過了")
)

// Migrations 是依序執行的資料庫結構變更。
// 正式專案中「資料表結構」也要做版本控制：只能在最後面「新增」，絕對不要修改已經上線的項目。
var Migrations = []string{
	// v1：書籍
	`CREATE TABLE books (
		id     INTEGER PRIMARY KEY AUTOINCREMENT,
		title  TEXT    NOT NULL,
		author TEXT    NOT NULL,
		year   INTEGER NOT NULL DEFAULT 0,
		stock  INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0)
	)`,
	// v2：借閱紀錄
	`CREATE TABLE loans (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		book_id     INTEGER  NOT NULL REFERENCES books(id),
		member      TEXT     NOT NULL,
		borrowed_at DATETIME NOT NULL,
		returned_at DATETIME
	)`,
	// v3：索引，加速「查某會員的借閱」
	`CREATE INDEX idx_loans_member ON loans(member)`,
}

// OpenDB 開啟（或建立）SQLite 資料庫檔案。
//
//   - foreign_keys(1)：SQLite 預設不檢查外鍵，要手動開啟
//   - busy_timeout(5000)：資料庫被鎖住時最多等 5 秒，而不是立刻失敗
//   - _txlock=immediate：交易一開始就取得寫入鎖，避免併發交易互相死結
//
// 換成 PostgreSQL 時，只需要改驅動程式與連線字串，database/sql 的用法完全相同。
func OpenDB(path string) (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_txlock=immediate", path)
	db, err := sql.Open("sqlite", dsn) // sql.Open 不會真的連線
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil { // Ping 才會真的連線，確認設定正確
		db.Close()
		return nil, fmt.Errorf("連線資料庫: %w", err)
	}
	return db, nil
}

// Store 封裝所有資料庫操作（Repository 模式）。
// 其他程式只透過 Store 的方法存取資料，不直接寫 SQL。
type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store { return &Store{db: db} }
