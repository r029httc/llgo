//go:build !solution

// 練習 09：資料庫（database/sql + SQLite）。題目說明請看 README.md。
package library

import (
	"context"
	"database/sql"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 9.1 ★★ Migrate 執行尚未套用的 Migrations，並記錄目前版本。可以重複執行。
func Migrate(ctx context.Context, db *sql.DB) error {
	panic(todo.NotImplemented)
}

// 9.2 ★ CreateBook 新增一本書，回傳新的 ID。
func (s *Store) CreateBook(ctx context.Context, b Book) (int64, error) {
	panic(todo.NotImplemented)
}

// 9.3 ★★ GetBook 依 ID 查詢；不存在時回傳包裝了 ErrNotFound 的錯誤。
func (s *Store) GetBook(ctx context.Context, id int64) (Book, error) {
	panic(todo.NotImplemented)
}

// 9.4 ★★ SearchBooks 以書名或作者模糊搜尋，依 ID 排序並支援分頁。
func (s *Store) SearchBooks(ctx context.Context, keyword string, limit, offset int) ([]Book, error) {
	panic(todo.NotImplemented)
}

// 9.5 ★★★ Borrow 在一個交易中：扣庫存 + 新增借閱紀錄。回傳借閱紀錄 ID。
func (s *Store) Borrow(ctx context.Context, bookID int64, member string) (int64, error) {
	panic(todo.NotImplemented)
}

// 9.6 ★★★ Return 在一個交易中：標記歸還時間 + 庫存加回。
func (s *Store) Return(ctx context.Context, loanID int64) error {
	panic(todo.NotImplemented)
}

// 9.7 ★★ ActiveLoans 回傳某會員尚未歸還的借閱（含書名），依借閱 ID 排序。
func (s *Store) ActiveLoans(ctx context.Context, member string) ([]Loan, error) {
	panic(todo.NotImplemented)
}
