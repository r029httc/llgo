package library

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// newDB 在暫存目錄建立一個新的資料庫檔案，測試結束自動關閉並刪除。
func newDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := OpenDB(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// newStore 直接執行所有 Migrations（不依賴你的 Migrate），讓各題可以獨立練習。
func newStore(t *testing.T) *Store {
	t.Helper()
	db := newDB(t)
	for _, m := range Migrations {
		if _, err := db.Exec(m); err != nil {
			t.Fatal(err)
		}
	}
	return NewStore(db)
}

// seed 直接用 SQL 塞測試資料，回傳 ID。
func seed(t *testing.T, s *Store, title, author string, year, stock int) int64 {
	t.Helper()
	res, err := s.db.Exec(`INSERT INTO books (title, author, year, stock) VALUES (?, ?, ?, ?)`, title, author, year, stock)
	if err != nil {
		t.Fatal(err)
	}
	id, _ := res.LastInsertId()
	return id
}

func TestMigrate(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx := context.Background()
	db := newDB(t)

	for i := range 2 { // 執行兩次：第二次不應出錯，也不應重複建表
		if err := Migrate(ctx, db); err != nil {
			t.Fatalf("第 %d 次 Migrate: %v", i+1, err)
		}
	}
	var version int
	db.QueryRow("PRAGMA user_version").Scan(&version)
	if version != len(Migrations) {
		t.Errorf("user_version = %d, want %d", version, len(Migrations))
	}
	for _, table := range []string{"books", "loans"} {
		var n int
		db.QueryRow(`SELECT count(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&n)
		if n != 1 {
			t.Errorf("資料表 %s 不存在", table)
		}
	}
}

func TestMigrateResumesFromVersion(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	db := newDB(t)
	// 模擬「舊版程式已經套用了 v1」
	db.Exec(Migrations[0])
	db.Exec("PRAGMA user_version = 1")
	if err := Migrate(context.Background(), db); err != nil {
		t.Fatalf("應從 v2 繼續，而不是重新執行 v1：%v", err)
	}
}

func TestCreateAndGetBook(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx := context.Background()
	s := newStore(t)

	id, err := s.CreateBook(ctx, Book{Title: "Go 程式設計", Author: "Donovan", Year: 2015, Stock: 2})
	if err != nil || id <= 0 {
		t.Fatalf("CreateBook = %d, %v", id, err)
	}
	got, err := s.GetBook(ctx, id)
	want := Book{ID: id, Title: "Go 程式設計", Author: "Donovan", Year: 2015, Stock: 2}
	if err != nil || got != want {
		t.Errorf("GetBook = %+v, %v; want %+v", got, err, want)
	}
	if _, err := s.GetBook(ctx, 9999); !errors.Is(err, ErrNotFound) {
		t.Errorf("不存在的書應 errors.Is ErrNotFound，得到 %v", err)
	}
	if _, err := s.CreateBook(ctx, Book{Title: "負庫存", Author: "x", Stock: -1}); err == nil {
		t.Error("stock < 0 違反 CHECK 約束，應回傳錯誤")
	}
}

func TestSearchBooks(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx := context.Background()
	s := newStore(t)
	seed(t, s, "Go 語言聖經", "Donovan", 2015, 1)
	seed(t, s, "Rust 程式設計", "Klabnik", 2018, 1)
	seed(t, s, "Go 併發實戰", "Cox-Buday", 2017, 1)
	seed(t, s, "Learning Go", "Bodner", 2021, 1)

	got, err := s.SearchBooks(ctx, "Go", 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	if titles := titlesOf(got); !equal(titles, []string{"Go 語言聖經", "Go 併發實戰", "Learning Go"}) {
		t.Errorf("SearchBooks(Go) = %v", titles)
	}
	if got, _ := s.SearchBooks(ctx, "Klabnik", 10, 0); len(got) != 1 {
		t.Errorf("應能用作者搜尋，得到 %v", titlesOf(got))
	}

	page2, _ := s.SearchBooks(ctx, "Go", 2, 2)
	if titles := titlesOf(page2); !equal(titles, []string{"Learning Go"}) {
		t.Errorf("分頁 limit=2 offset=2 = %v, want [Learning Go]", titles)
	}

	// SQL injection：如果你用字串拼接 SQL，這個關鍵字會讓 WHERE 永遠成立
	evil, err := s.SearchBooks(ctx, "' OR '1'='1", 10, 0)
	if err != nil || len(evil) != 0 {
		t.Errorf("SQL injection 測試：得到 %d 筆、err=%v；必須用 ? 參數而不是字串拼接！", len(evil), err)
	}
	if none, err := s.SearchBooks(ctx, "不存在", 10, 0); err != nil || len(none) != 0 {
		t.Errorf("沒有結果時應回傳空切片與 nil，得到 %v, %v", none, err)
	}
}

func TestBorrowAndReturn(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx := context.Background()
	s := newStore(t)
	bookID := seed(t, s, "Go", "Pike", 2009, 1)

	loanID, err := s.Borrow(ctx, bookID, "alice")
	if err != nil || loanID <= 0 {
		t.Fatalf("Borrow = %d, %v", loanID, err)
	}
	if stock := stockOf(t, s, bookID); stock != 0 {
		t.Errorf("借出後庫存 = %d, want 0", stock)
	}
	if _, err := s.Borrow(ctx, bookID, "bob"); !errors.Is(err, ErrOutOfStock) {
		t.Errorf("沒庫存應 ErrOutOfStock，得到 %v", err)
	}
	if _, err := s.Borrow(ctx, 9999, "bob"); !errors.Is(err, ErrNotFound) {
		t.Errorf("不存在的書應 ErrNotFound，得到 %v", err)
	}
	if n := loanCount(t, s); n != 1 {
		t.Errorf("失敗的借閱不應留下紀錄（交易要回滾），loans 有 %d 筆", n)
	}

	if err := s.Return(ctx, loanID); err != nil {
		t.Fatalf("Return: %v", err)
	}
	if stock := stockOf(t, s, bookID); stock != 1 {
		t.Errorf("歸還後庫存 = %d, want 1", stock)
	}
	if err := s.Return(ctx, loanID); !errors.Is(err, ErrAlreadyReturned) {
		t.Errorf("重複歸還應 ErrAlreadyReturned，得到 %v", err)
	}
	if stock := stockOf(t, s, bookID); stock != 1 {
		t.Errorf("重複歸還不應再增加庫存，得到 %d", stock)
	}
	if err := s.Return(ctx, 9999); !errors.Is(err, ErrNotFound) {
		t.Errorf("不存在的借閱應 ErrNotFound，得到 %v", err)
	}
}

func TestBorrowConcurrent(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx := context.Background()
	s := newStore(t)
	bookID := seed(t, s, "熱門書", "x", 2024, 4)
	// 先同步借走 1 本（剩 3 本）。尚未實作時會在這裡 panic 並被 SKIP，而不是在 goroutine 裡讓整個測試當掉
	if _, err := s.Borrow(ctx, bookID, "first"); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	ok, outOfStock := 0, 0
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.Borrow(ctx, bookID, "someone")
			mu.Lock()
			defer mu.Unlock()
			switch {
			case err == nil:
				ok++
			case errors.Is(err, ErrOutOfStock):
				outOfStock++
			default:
				t.Errorf("非預期的錯誤: %v", err)
			}
		}()
	}
	wg.Wait()
	if ok != 3 || outOfStock != 17 || stockOf(t, s, bookID) != 0 || loanCount(t, s) != 4 {
		t.Errorf("20 人搶 3 本：成功 %d、庫存不足 %d、剩餘庫存 %d、紀錄 %d 筆；want 3/17/0/4",
			ok, outOfStock, stockOf(t, s, bookID), loanCount(t, s))
	}
}

func TestActiveLoans(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx := context.Background()
	s := newStore(t)
	b1 := seed(t, s, "書一", "a", 2020, 5)
	b2 := seed(t, s, "書二", "b", 2021, 5)

	l1, _ := s.Borrow(ctx, b1, "alice")
	l2, _ := s.Borrow(ctx, b2, "alice")
	s.Borrow(ctx, b1, "bob")
	s.Return(ctx, l1)

	before := time.Now().Add(-time.Minute)
	loans, err := s.ActiveLoans(ctx, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(loans) != 1 {
		t.Fatalf("alice 應有 1 筆未歸還，得到 %+v", loans)
	}
	l := loans[0]
	if l.ID != l2 || l.BookID != b2 || l.BookTitle != "書二" || l.Member != "alice" || l.ReturnedAt != nil {
		t.Errorf("ActiveLoans[0] = %+v", l)
	}
	if l.BorrowedAt.Before(before) || l.BorrowedAt.After(time.Now().Add(time.Minute)) {
		t.Errorf("BorrowedAt = %v，應該是剛剛的時間", l.BorrowedAt)
	}
	if none, _ := s.ActiveLoans(ctx, "nobody"); len(none) != 0 {
		t.Errorf("沒有借閱的會員應回傳空，得到 %+v", none)
	}
}

func stockOf(t *testing.T, s *Store, id int64) int {
	t.Helper()
	var n int
	if err := s.db.QueryRow(`SELECT stock FROM books WHERE id = ?`, id).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func loanCount(t *testing.T, s *Store) int {
	t.Helper()
	var n int
	s.db.QueryRow(`SELECT count(*) FROM loans`).Scan(&n)
	return n
}

func titlesOf(books []Book) []string {
	var out []string
	for _, b := range books {
		out = append(out, b.Title)
	}
	return out
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
