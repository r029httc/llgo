# 練習 09：期末專案

前面的練習都是「填空題」；這一章是**從零開始**做出一個完整的程式。沒有自動測試——你要自己寫。

選一個專題（建議依序做），每個都拆成幾個里程碑，**每完成一個里程碑就 commit 一次**。

---

## 專題 A：命令列待辦清單 `cmd/todo`（綜合 02、04、07）

```bash
go run ./cmd/todo add "買牛奶"
go run ./cmd/todo add "寫 Go 練習" -due 2026-12-31
go run ./cmd/todo list
#  1  [ ]  買牛奶
#  2  [ ]  寫 Go 練習   (到期：2026-12-31)
go run ./cmd/todo done 1
go run ./cmd/todo list -all
go run ./cmd/todo rm 2
```

**里程碑**
1. **子命令解析**：用 `os.Args[1]` 判斷子命令，每個子命令用自己的 `flag.NewFlagSet` 解析參數。未知命令印出使用說明並 `os.Exit(2)`。
2. **資料儲存**：存在 `~/.todo.json`（`os.UserHomeDir()`）。重用練習 07 的 `SaveTodos` / `LoadTodos`。檔案不存在時視為空清單（`errors.Is(err, os.ErrNotExist)`）。
3. **安全寫檔**：先寫到暫存檔再 `os.Rename` 取代原檔，避免寫到一半當機導致資料毀損。
4. **分層設計**：把邏輯拆成 `internal/todo`（純邏輯，不碰檔案和終端機）與 `cmd/todo`（處理輸入輸出）。純邏輯部分要有單元測試。
5. **可測試的 main**：把 `main` 的內容移到 `run(args []string, stdout io.Writer) error`，`main` 只負責呼叫它並處理錯誤。這樣就能在測試中模擬命令列。
6. **加分**：`-file` 參數指定資料檔；用 `text/tabwriter` 對齊輸出；`list -json` 輸出 JSON。

**自我檢查**
- [ ] `go vet ./...` 沒有警告
- [ ] `internal/todo` 覆蓋率 ≥ 80%（`go test -cover`）
- [ ] 錯誤訊息清楚（「找不到 ID 5 的待辦事項」而不是 panic）

---

## 專題 B：短網址服務 `cmd/shorturl`（綜合 05、07、08）

**規格**
- `POST /shorten` body `{"url": "https://go.dev"}` → `{"code": "aZ3k9Q", "short": "http://localhost:8080/aZ3k9Q"}`
- `GET /{code}` → 302 重新導向到原網址；不存在則 404
- `GET /stats/{code}` → `{"url": "...", "hits": 12, "created": "..."}`

**里程碑**
1. 以 `map` + `sync.RWMutex` 實作記憶體儲存；用 `crypto/rand` 產生 6 碼的 code，碰撞時重新產生。
2. 用 `net/url.Parse` 驗證網址，只接受 `http` / `https`。
3. 點擊數統計：用 `sync/atomic` 或鎖。用 `-race` 加上 100 個併發請求的測試證明正確。
4. 定期（每 30 秒）把資料存到 JSON 檔；啟動時讀回。用 `time.Ticker` + `context` 控制背景 goroutine 的生命週期。
5. 優雅關機：Ctrl+C 時存檔再結束。
6. **加分**：限流（每個 IP 每分鐘最多 10 次 POST），使用 `golang.org/x/time/rate`；加上 `log/slog` 結構化日誌。

---

## 專題 C：併發網站健康檢查器 `cmd/healthcheck`（綜合 05、08）

```bash
go run ./cmd/healthcheck -c 5 -timeout 3s urls.txt
# ✔ 200  https://go.dev         123ms
# ✘ ---  https://example.invalid  dial tcp: lookup ... no such host
# 摘要：9 成功、1 失敗，平均 210ms
```

**里程碑**
1. 從檔案讀取 URL 清單（每行一個，略過空行與 `#`）。
2. 用練習 5.1 的 worker pool 模式，最多 `-c` 個同時請求；每個請求有 `-timeout` 逾時。
3. 結果依輸入順序輸出；最後印出摘要。
4. 有任何失敗時程式以結束碼 1 結束（方便放進 CI 或 cron）。
5. 測試：用 `httptest.NewServer` 模擬快、慢、失敗的網站。
6. **加分**：`-watch 30s` 參數每 30 秒重新檢查一次；狀態從成功變失敗時才印出通知。

---

## 完成之後

- 把專案放上 GitHub，寫一份好的 README（安裝方式、使用範例、設計說明）。
- 用 `go build` 交叉編譯：`GOOS=windows GOARCH=amd64 go build -o todo.exe ./cmd/todo`。
- 設定 GitHub Actions 在每次 push 時跑 `go vet` 與 `go test -race`（本 repo 的 `.github/workflows/go.yml` 可以參考）。
- 安裝 [`golangci-lint`](https://golangci-lint.run/) 或 [`staticcheck`](https://staticcheck.dev/)，看看它們會抓到什麼問題。
