# llgo — 用 Go 學習軟體開發

這個 repo 是一份循序漸進的 Go 學習專案。每一課都是一個獨立的套件，包含**有註解的範例程式**與**測試**，你可以邊讀、邊改、邊跑測試來學習。

## 環境準備

1. 安裝 Go（建議 1.24 以上）：<https://go.dev/dl/>
2. 確認安裝成功：
   ```bash
   go version
   ```
3. 取得專案：
   ```bash
   git clone https://github.com/r029httc/llgo.git
   cd llgo
   ```
4. 下載相依套件（資料庫驅動 `modernc.org/sqlite`、`golang.org/x/crypto`）：
   ```bash
   go mod download
   ```
5. 編輯器建議使用 VS Code + 官方 Go 擴充套件（會自動格式化、補全、跑測試）。

## 第一個程式

```bash
go run ./cmd/hello                 # Hello, Gopher!
go run ./cmd/hello -name 你的名字   # Hello, 你的名字!
```

## 專案結構

```
llgo/
├── go.mod                          # 模組定義（模組路徑 github.com/r029httc/llgo）
├── cmd/hello/main.go               # 可執行程式（package main）
└── lessons/
    ├── 01-basics/                  # 變數、常數、函式、多回傳值、for、switch
    ├── 02-collections/             # slice、map、泛型
    ├── 03-structs-interfaces/      # struct、方法、介面、指標接收者
    ├── 04-errors/                  # error、包裝錯誤、errors.Is / errors.As
    └── 05-concurrency/             # goroutine、channel、WaitGroup、Mutex
└── exercises/                      # 62 題練習（附測試與參考解答），見 exercises/README.md
```

## 常用指令

| 指令 | 用途 |
| --- | --- |
| `go run ./cmd/hello` | 編譯並執行 |
| `go build -o bin/hello ./cmd/hello` | 編譯成執行檔 |
| `go test ./...` | 執行所有測試 |
| `go test -v ./lessons/01-basics` | 執行單一課程測試並顯示細節 |
| `go test -run TestFizzBuzz ./...` | 只跑名稱符合的測試 |
| `go test -race ./...` | 偵測資料競爭（併發課程必用） |
| `go test -cover ./...` | 顯示測試覆蓋率 |
| `go fmt ./...` | 自動格式化程式碼 |
| `go vet ./...` | 靜態檢查常見錯誤 |
| `go doc ./lessons/01-basics` | 查看套件文件 |
| `go test -tags solution ./exercises/...` | 用參考解答跑練習測試 |

## 學習路線

建議的學習方式：**先讀課程程式碼與註解 → 跑課程測試 → 做對應的練習題 → 全部測試通過**。

| 順序 | 課程（範例） | 練習題 | 主題 |
| --- | --- | --- | --- |
| 1 | [`lessons/01-basics`](lessons/01-basics) | [`exercises/01-basics`](exercises/01-basics) | 基礎語法、字串與 rune |
| 2 | [`lessons/02-collections`](lessons/02-collections) | [`exercises/02-collections`](exercises/02-collections) | slice、map、泛型 |
| 3 | [`lessons/03-structs-interfaces`](lessons/03-structs-interfaces) | [`exercises/03-structs-interfaces`](exercises/03-structs-interfaces) | struct、介面、嵌入 |
| 4 | [`lessons/04-errors`](lessons/04-errors) | [`exercises/04-errors`](exercises/04-errors) | 錯誤處理 |
| 5 | [`lessons/05-concurrency`](lessons/05-concurrency) | [`exercises/05-concurrency`](exercises/05-concurrency) | goroutine、channel、context |
| 6 | — | [`exercises/06-testing`](exercises/06-testing) | 測試、benchmark、fuzz |
| 7 | — | [`exercises/07-io-json`](exercises/07-io-json) | io.Reader/Writer、JSON、CSV |
| 8 | — | [`exercises/08-http`](exercises/08-http) | HTTP 伺服器、中介層、用戶端 |
| 9 | — | [`exercises/09-database`](exercises/09-database) | 資料庫：SQL、交易、migration、防 SQL injection |
| 10 | — | [`exercises/10-network`](exercises/10-network) | 網路通訊：TCP、聊天室、pub/sub、SSE 即時推送 |
| 11 | — | [`exercises/11-auth`](exercises/11-auth) | 認證與安全：bcrypt、簽章 token、限流 |
| 12 | — | [`exercises/12-project`](exercises/12-project) | 期末專案（四選一） |

**👉 從 [`exercises/README.md`](exercises/README.md) 開始做練習**：裡面有做題流程、難度說明與完整進度表。
每一題都有詳細規格、提示、思考題與延伸挑戰，並附自動評分測試與參考解答。

## 接下來可以學什麼

- 標準函式庫：`net/http`（寫 Web API）、`encoding/json`、`os`、`io`、`context`
- 建立 CLI 工具或 REST API 小專案，並放進 `cmd/` 底下
- 推薦資源：
  - [A Tour of Go](https://go.dev/tour/)（有繁體中文版）
  - [Go by Example](https://gobyexample.com/)
  - [Effective Go](https://go.dev/doc/effective_go)
  - [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)
