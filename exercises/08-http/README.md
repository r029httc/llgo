# 練習 08：HTTP 伺服器與用戶端

> 檔案：`exercise.go`（你要寫的）、`types.go`（已定義好的型別）

`net/http` 是 Go 最常被使用的套件之一，不需要任何框架就能寫出正式環境等級的 Web 服務。
本章測試使用 `net/http/httptest`：不需要真的開 port，就能測試 handler。

寫完後可以真的把它跑起來，見本頁最後一節。

---

## 8.1 ★★★ `NewServer() http.Handler` — Todo REST API

**規格**

| 方法 | 路徑 | 行為 | 狀態碼 |
| --- | --- | --- | --- |
| GET | `/hello?name=X` | 回傳純文字 `Hello, X!`；沒有 name 時用 `World` | 200 |
| GET | `/todos` | 回傳所有 todo 的 JSON 陣列；沒有資料時是 `[]` | 200 |
| POST | `/todos` | body `{"title": "..."}`，建立 todo（ID 從 1 開始遞增），回傳建立的 todo | 201 |
| POST | `/todos` | JSON 無效或 title 為空 | 400 |
| GET | `/todos/{id}` | 回傳該 todo | 200 |
| GET | `/todos/{id}` | 不存在 | 404 |
| GET | `/todos/{id}` | id 不是整數 | 400 |
| DELETE | `/todos` | 不支援的方法 | 405 |

另外，**50 個請求同時新增**時，不能遺失資料或產生重複 ID（測試會檢查，請用 `-race` 執行）。

**提示**
- Go 1.22+ 的 `http.ServeMux` 支援方法與路徑參數：
  ```go
  mux.HandleFunc("GET /todos/{id}", handler)
  id := r.PathValue("id")
  ```
  有註冊 `GET /todos` 和 `POST /todos` 後，對 `/todos` 發 DELETE 會**自動**得到 405。
- 把狀態放在一個 `server` 結構裡（`todos []Todo`、`nextID int`、`mu sync.Mutex`），handler 寫成它的方法。
- **HTTP handler 會在多個 goroutine 中同時執行**，所有共享狀態都必須加鎖。
- 回應 JSON 的順序：先 `w.Header().Set("Content-Type", "application/json")`，再 `w.WriteHeader(201)`，最後寫 body。順序錯了 header 會失效。
- 錯誤回應可以用 `http.Error(w, "訊息", http.StatusBadRequest)`。

**延伸挑戰**
1. 加上 `PATCH /todos/{id}`（更新 done）與 `DELETE /todos/{id}`（回 204）。
2. 把儲存層抽成介面 `type Store interface { List() []Todo; Add(title string) Todo; ... }`，讓 server 可以換成檔案或資料庫實作。
3. 加上 `GET /todos?done=true` 篩選。

---

## 8.2 ★★ `RequireAPIKey(key string, next http.Handler) http.Handler`

**規格**：請求標頭 `X-API-Key` 不等於 `key` 時回應 401，**而且不能呼叫** `next`；相等時把請求交給 `next`。

**觀念：中介層（middleware）**

```go
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 前處理
		next.ServeHTTP(w, r)
		// 後處理
	})
}
```

中介層是一個「接收 handler、回傳 handler」的函式，可以一層層包起來：`Logging(RequireAPIKey(key, NewServer()))`。

**延伸挑戰**
1. 寫 `Logging(next)`：印出方法、路徑、狀態碼與耗時。提示：要記錄狀態碼，需要包裝 `http.ResponseWriter`，攔截 `WriteHeader`。
2. 寫 `Recover(next)`：handler panic 時回傳 500 而不是讓連線斷掉（用到練習 4.5 的技巧）。
3. 安全性：比較金鑰時改用 `crypto/subtle.ConstantTimeCompare`，避免時序攻擊（timing attack）。

---

## 8.3 ★★★ `FetchTodo(ctx, client, baseURL, id) (Todo, error)`

**規格**
- 發送 `GET {baseURL}/todos/{id}`。
- 200 → 解析 JSON 回傳；404 → 包裝 `ErrNotFound`；其他狀態碼 → 回傳錯誤（不是 `ErrNotFound`）。
- JSON 無效 → 回傳錯誤。
- `ctx` 逾時 → 立即返回，錯誤要滿足 `errors.Is(err, context.DeadlineExceeded)`。

**提示**
- 用 `http.NewRequestWithContext(ctx, http.MethodGet, url, nil)` 建立請求，再 `client.Do(req)`。不要用 `http.Get`，它無法帶 context。
- **`defer resp.Body.Close()`**：忘了關的話，連線無法重用，久了會耗盡資源。
- 不要使用 `http.DefaultClient` 寫正式程式：它沒有逾時設定。正式環境應該建立 `&http.Client{Timeout: 10 * time.Second}`。

**學到的概念**：HTTP 用戶端、`context` 傳遞逾時、資源釋放、狀態碼處理

---

## 實際跑起來

寫完 8.1 和 8.2 後，建立 `cmd/todoserver/main.go`：

```go
package main

import (
	"log"
	"net/http"

	webapi "github.com/r029httc/llgo/exercises/08-http"
)

func main() {
	h := webapi.RequireAPIKey("dev-key", webapi.NewServer())
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", h))
}
```

```bash
go run ./cmd/todoserver
# 另開一個終端機：
curl -H 'X-API-Key: dev-key' localhost:8080/hello?name=Go
curl -H 'X-API-Key: dev-key' -X POST -d '{"title":"學 Go"}' localhost:8080/todos
curl -H 'X-API-Key: dev-key' localhost:8080/todos
```

**延伸挑戰**：用 `signal.NotifyContext` 和 `server.Shutdown(ctx)` 實作**優雅關機**：按 Ctrl+C 時，先等進行中的請求處理完再結束。
