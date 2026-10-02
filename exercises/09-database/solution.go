//go:build solution

// 參考解答。建議自己寫完再看！
package library

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func Migrate(ctx context.Context, db *sql.DB) error {
	var version int
	// SQLite 內建的 user_version 很適合記錄結構版本；PostgreSQL 通常用一張 schema_migrations 表
	if err := db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return fmt.Errorf("讀取版本: %w", err)
	}
	for i := version; i < len(Migrations); i++ {
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		// 每個版本在自己的交易裡執行：失敗就整個回滾，不會留下做一半的結構
		if _, err := tx.ExecContext(ctx, Migrations[i]); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration v%d: %w", i+1, err)
		}
		// PRAGMA 不支援 ? 參數，這裡的 i 是程式內部的整數，不是使用者輸入，所以可以用 Sprintf
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", i+1)); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) CreateBook(ctx context.Context, b Book) (int64, error) {
	res, err := s.db.ExecContext(ctx,
		`INSERT INTO books (title, author, year, stock) VALUES (?, ?, ?, ?)`,
		b.Title, b.Author, b.Year, b.Stock)
	if err != nil {
		return 0, fmt.Errorf("新增書籍: %w", err)
	}
	return res.LastInsertId()
}

func (s *Store) GetBook(ctx context.Context, id int64) (Book, error) {
	var b Book
	err := s.db.QueryRowContext(ctx,
		`SELECT id, title, author, year, stock FROM books WHERE id = ?`, id,
	).Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Stock)
	if errors.Is(err, sql.ErrNoRows) {
		return Book{}, fmt.Errorf("書籍 %d: %w", id, ErrNotFound)
	}
	if err != nil {
		return Book{}, fmt.Errorf("查詢書籍 %d: %w", id, err)
	}
	return b, nil
}

func (s *Store) SearchBooks(ctx context.Context, keyword string, limit, offset int) ([]Book, error) {
	pattern := "%" + keyword + "%" // 萬用字元加在「值」裡，SQL 本身仍然用 ? 參數
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, title, author, year, stock FROM books
		 WHERE title LIKE ? OR author LIKE ?
		 ORDER BY id LIMIT ? OFFSET ?`,
		pattern, pattern, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("搜尋書籍: %w", err)
	}
	defer rows.Close() // 沒關閉的話，連線會一直被佔用

	var books []Book
	for rows.Next() {
		var b Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Stock); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, rows.Err() // 迴圈中途發生的錯誤要從這裡取得
}

func (s *Store) Borrow(ctx context.Context, bookID int64, member string) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback() // Commit 之後再 Rollback 是無害的 no-op；確保任何錯誤路徑都會回滾

	// 「檢查 + 扣除」用同一條 UPDATE 完成，條件 stock > 0 保證不會超借（不會有競態條件）
	res, err := tx.ExecContext(ctx, `UPDATE books SET stock = stock - 1 WHERE id = ? AND stock > 0`, bookID)
	if err != nil {
		return 0, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		// 沒有更新到任何列：可能是書不存在，也可能是沒庫存，要分辨清楚
		var exists bool
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM books WHERE id = ?)`, bookID).Scan(&exists); err != nil {
			return 0, err
		}
		if !exists {
			return 0, fmt.Errorf("書籍 %d: %w", bookID, ErrNotFound)
		}
		return 0, fmt.Errorf("書籍 %d: %w", bookID, ErrOutOfStock)
	}

	res, err = tx.ExecContext(ctx,
		`INSERT INTO loans (book_id, member, borrowed_at) VALUES (?, ?, ?)`,
		bookID, member, time.Now().UTC())
	if err != nil {
		return 0, err
	}
	loanID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return loanID, tx.Commit()
}

func (s *Store) Return(ctx context.Context, loanID int64) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var bookID int64
	var returnedAt sql.NullTime // 可能為 NULL 的欄位要用 sql.NullXxx 或指標接收
	err = tx.QueryRowContext(ctx, `SELECT book_id, returned_at FROM loans WHERE id = ?`, loanID).Scan(&bookID, &returnedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("借閱 %d: %w", loanID, ErrNotFound)
	}
	if err != nil {
		return err
	}
	if returnedAt.Valid {
		return fmt.Errorf("借閱 %d: %w", loanID, ErrAlreadyReturned)
	}

	if _, err := tx.ExecContext(ctx, `UPDATE loans SET returned_at = ? WHERE id = ?`, time.Now().UTC(), loanID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE books SET stock = stock + 1 WHERE id = ?`, bookID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) ActiveLoans(ctx context.Context, member string) ([]Loan, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT l.id, l.book_id, b.title, l.member, l.borrowed_at, l.returned_at
		 FROM loans l JOIN books b ON b.id = l.book_id
		 WHERE l.member = ? AND l.returned_at IS NULL
		 ORDER BY l.id`, member)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []Loan
	for rows.Next() {
		var l Loan
		if err := rows.Scan(&l.ID, &l.BookID, &l.BookTitle, &l.Member, &l.BorrowedAt, &l.ReturnedAt); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}
