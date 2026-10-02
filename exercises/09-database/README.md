# 練習 09：資料庫（database/sql + SQLite）

> 檔案：`exercise.go`（你要寫的）、`types.go`（資料表結構、`OpenDB`、`Store`）

幾乎所有後端服務都需要資料庫。本章用一個**圖書館借閱系統**練習 Go 標準的資料庫介面 `database/sql`。

## 為什麼用 SQLite？

- 資料庫就是一個檔案，**不需要安裝或啟動任何伺服器**，`go test` 就能跑。
- 使用 `modernc.org/sqlite`：純 Go 實作，不需要 C 編譯器（cgo）。
- `database/sql` 是**統一介面**：換成 PostgreSQL 或 MySQL 時，只要換驅動程式與連線字串，你寫的程式碼幾乎不用改（差別主要是佔位符號：SQLite/MySQL 用 `?`，PostgreSQL 用 `$1`）。

## 資料表

```
books                                   loans
┌────┬───────┬────────┬──────┬───────┐  ┌────┬─────────┬────────┬─────────────┬─────────────┐
│ id │ title │ author │ year │ stock │  │ id │ book_id │ member │ borrowed_at │ returned_at │
└────┴───────┴────────┴──────┴───────┘  └────┴────┬────┴────────┴─────────────┴─────────────┘
  ▲                                               │                              NULL = 未歸還
  └──────────────── 外鍵 REFERENCES ──────────────┘
```

## database/sql 速查表

| 用途 | 方法 | 注意 |
| --- | --- | --- |
| 執行 INSERT / UPDATE / DELETE | `db.ExecContext(ctx, sql, args...)` | 回傳 `sql.Result`：`LastInsertId()`、`RowsAffected()` |
| 查詢**一列** | `db.QueryRowContext(...).Scan(&a, &b)` | 沒資料時 `Scan` 回傳 `sql.ErrNoRows` |
| 查詢**多列** | `rows, err := db.QueryContext(...)` | 必須 `defer rows.Close()`，迴圈後檢查 `rows.Err()` |
| 交易 | `tx, err := db.BeginTx(ctx, nil)` | `defer tx.Rollback()` + 最後 `tx.Commit()` |
| 可能為 NULL 的欄位 | `sql.NullString`、`sql.NullTime`、或指標 `*time.Time` | 直接 Scan 進 `string` 遇到 NULL 會出錯 |

> 💡 **永遠用 `?` 佔位符號傳參數，絕對不要用 `fmt.Sprintf` 或 `+` 把使用者輸入拼進 SQL。**
> 否則就會有 SQL injection 漏洞——這是 OWASP 十大網站安全風險之一。9.4 的測試會實際攻擊你的程式。

---

## 9.1 ★★ `Migrate(ctx, db) error` — 資料庫版本遷移

**規格**
- 依序執行 `Migrations` 中**尚未套用**的項目，並記錄目前版本。
- 可以重複執行：第二次執行時不做任何事、也不出錯。
- 從中間版本繼續：如果資料庫已經在 v1，只執行 v2、v3。

**提示**
- SQLite 每個資料庫檔都有一個整數 `PRAGMA user_version`（預設 0），很適合存版本號：
  ```go
  db.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version)
  ```
- `for i := version; i < len(Migrations); i++`，每一個版本放在**自己的交易**裡：執行 SQL + 更新 `user_version`，成功才 `Commit`。
- `PRAGMA user_version = ?` 不支援參數，要用 `fmt.Sprintf`——這裡可以，因為數字來自程式內部，不是使用者輸入。

**觀念**：正式專案中，資料表結構也需要版本控制。已經上線的 migration **絕對不能修改**，只能新增。常用工具有 [goose](https://github.com/pressly/goose)、[golang-migrate](https://github.com/golang-migrate/migrate)。

---

## 9.2 ★ `CreateBook(ctx, b Book) (int64, error)`

**規格**：新增一本書，回傳資料庫自動產生的 ID。`stock < 0` 會違反資料表的 `CHECK` 約束，要回傳錯誤。

**提示**：`ExecContext` 搭配 `INSERT ... VALUES (?, ?, ?, ?)`，再呼叫 `res.LastInsertId()`。

**觀念**：讓資料庫幫你守住規則（`NOT NULL`、`CHECK`、外鍵），就算程式有 bug，壞資料也進不去。

---

## 9.3 ★★ `GetBook(ctx, id) (Book, error)`

**規格**：依 ID 查詢。不存在時回傳的錯誤要滿足 `errors.Is(err, ErrNotFound)`。

**提示**
```go
err := s.db.QueryRowContext(ctx, `SELECT ... WHERE id = ?`, id).Scan(&b.ID, &b.Title, ...)
if errors.Is(err, sql.ErrNoRows) { ... }
```

**觀念：為什麼要把 `sql.ErrNoRows` 轉成自己的 `ErrNotFound`？**
呼叫端（例如 HTTP handler）只需要知道「找不到 → 回 404」，不應該依賴「底層是 SQL 資料庫」這個細節。這樣以後換成 Redis 或其他儲存方式時，呼叫端不用改。

---

## 9.4 ★★ `SearchBooks(ctx, keyword, limit, offset) ([]Book, error)`

**規格**
- 書名**或**作者包含 `keyword` 的書，依 ID 排序。
- `limit` / `offset` 做分頁。
- 沒有結果時回傳空切片與 `nil`（不是錯誤）。
- 必須擋住 SQL injection：測試會搜尋 `' OR '1'='1`。

**提示**
- `WHERE title LIKE ? OR author LIKE ? ORDER BY id LIMIT ? OFFSET ?`
- 萬用字元 `%` 加在**參數值**裡：`pattern := "%" + keyword + "%"`。
- 多列查詢的標準寫法：
  ```go
  rows, err := s.db.QueryContext(ctx, query, args...)
  if err != nil { return nil, err }
  defer rows.Close()
  for rows.Next() {
      var b Book
      if err := rows.Scan(...); err != nil { return nil, err }
      books = append(books, b)
  }
  return books, rows.Err()
  ```

**思考題**：如果使用者搜尋 `100%`，`%` 會被當成萬用字元。怎麼讓它被當成普通字元？（提示：`LIKE ? ESCAPE '\'`）

**延伸挑戰**：`OFFSET` 在資料量大時很慢（資料庫還是要掃過前面所有列）。研究「游標分頁」（keyset pagination）：`WHERE id > ? ORDER BY id LIMIT ?`。

---

## 9.5 ★★★ `Borrow(ctx, bookID, member) (int64, error)` — 交易

**規格**：在**同一個交易**中：
1. 庫存減 1；
2. 新增一筆借閱紀錄（`borrowed_at` 為現在時間）。

- 書不存在 → `ErrNotFound`；沒有庫存 → `ErrOutOfStock`。
- 失敗時不能留下任何變更（測試會檢查 loans 表沒有多出紀錄）。
- **20 個人同時搶 3 本書，恰好 3 個人成功**（測試會真的開 20 個 goroutine）。

**提示**
- 交易的標準寫法：
  ```go
  tx, err := s.db.BeginTx(ctx, nil)
  if err != nil { return 0, err }
  defer tx.Rollback() // Commit 成功後再 Rollback 不會有任何作用；任何提早 return 都會自動回滾
  ...
  return id, tx.Commit()
  ```
- 交易中要用 `tx.ExecContext`，不是 `s.db.ExecContext`（後者會在交易**外**執行）。
- 防止超借的關鍵：**檢查和扣除用同一條 SQL 完成**：
  ```sql
  UPDATE books SET stock = stock - 1 WHERE id = ? AND stock > 0
  ```
  然後看 `RowsAffected()`：0 代表「書不存在」或「沒庫存」，再查一次分辨是哪一種。

**思考題（很重要）**：如果寫成「先 `SELECT stock`，在 Go 裡判斷 `stock > 0`，再 `UPDATE stock = ?`」，在 PostgreSQL 這種多人同時寫入的資料庫上，兩個請求可能同時讀到 `stock = 1`，結果都借成功——這叫 **lost update**（更新遺失）。原子性的 `UPDATE ... WHERE stock > 0` 或 `SELECT ... FOR UPDATE` 可以避免它。

---

## 9.6 ★★★ `Return(ctx, loanID) error`

**規格**：在同一個交易中：設定 `returned_at` + 庫存加 1。
- 借閱紀錄不存在 → `ErrNotFound`；已經歸還過 → `ErrAlreadyReturned`（庫存不能再加）。

**提示**：`returned_at` 可能是 NULL，用 `sql.NullTime` 接收，`.Valid` 為 `false` 代表 NULL。

---

## 9.7 ★★ `ActiveLoans(ctx, member) ([]Loan, error)` — JOIN

**規格**：某會員**尚未歸還**的借閱，要包含書名（`BookTitle`），依借閱 ID 排序。

**提示**
```sql
SELECT l.id, l.book_id, b.title, l.member, l.borrowed_at, l.returned_at
FROM loans l JOIN books b ON b.id = l.book_id
WHERE l.member = ? AND l.returned_at IS NULL
ORDER BY l.id
```
- 注意 SQL 中判斷 NULL 要用 `IS NULL`，`= NULL` 永遠不成立。
- `returned_at` 可以直接 Scan 進 `*time.Time`：NULL 時會是 `nil`。

**延伸挑戰**：如果改成「先查出所有借閱，再對每一筆分別查書名」，會發生什麼事？這叫 **N+1 查詢問題**，是 ORM 最常見的效能陷阱。

---

## 綜合延伸挑戰

1. **接上 HTTP**：結合練習 08，做出 `GET /books?q=go&page=2`、`POST /books/{id}/borrow` 的 REST API。資料庫錯誤要轉成正確的狀態碼（`ErrNotFound` → 404、`ErrOutOfStock` → 409 Conflict）。
2. **用介面解耦**：定義 `type BookStore interface { GetBook(...); ... }`，HTTP handler 只依賴介面。測試 handler 時就可以傳入一個假的記憶體實作，不需要真的資料庫。
3. **換成 PostgreSQL**：用 Docker 啟動 PostgreSQL（`docker run -e POSTGRES_PASSWORD=pw -p 5432:5432 postgres`），改用 [`github.com/jackc/pgx/v5/stdlib`](https://github.com/jackc/pgx) 驅動，把 `?` 換成 `$1`，看看還需要改哪些地方。
4. **連線池設定**：研究 `db.SetMaxOpenConns`、`SetMaxIdleConns`、`SetConnMaxLifetime` 的意義，以及為什麼正式環境一定要設定。
5. **認識工具**：[sqlc](https://sqlc.dev/)（從 SQL 產生型別安全的 Go 程式碼）、[sqlx](https://github.com/jmoiron/sqlx)（Scan 到 struct）、[GORM](https://gorm.io/)（ORM）。比較它們和純 `database/sql` 的取捨。
