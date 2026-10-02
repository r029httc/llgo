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
4. 編輯器建議使用 VS Code + 官方 Go 擴充套件（會自動格式化、補全、跑測試）。

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

## 學習路線與練習

建議的學習方式：**先讀程式碼與註解 → 跑測試 → 動手做練習 → 為練習寫測試 → 跑 `go test` 確認通過**。

### 01 基礎語法 `lessons/01-basics`
- 重點：`:=` 與 `var`、匯出規則（大寫開頭）、多回傳值、`for range`、`switch`、可變參數
- 練習：
  1. 寫一個 `Max(nums ...int) int`，回傳最大值。
  2. 寫一個 `IsPrime(n int) bool`，並用表格驅動測試驗證。

### 02 集合 `lessons/02-collections`
- 重點：slice 的長度與容量、`append`、`make`、map 的零值、map 迭代無序、泛型
- 練習：
  1. 寫泛型函式 `Map[T, U any](s []T, f func(T) U) []U`。
  2. 寫 `Unique(s []string) []string`，移除重複且保持原順序。

### 03 結構與介面 `lessons/03-structs-interfaces`
- 重點：方法接收者（值 vs 指標）、隱式實作介面、`fmt.Stringer`
- 練習：
  1. 新增 `Triangle` 型別並實作 `Shape`。
  2. 想想看：如果把 `Counter.Inc` 改成值接收者 `(c Counter)`，測試會發生什麼事？為什麼？

### 04 錯誤處理 `lessons/04-errors`
- 重點：`if err != nil`、`fmt.Errorf("%w")` 包裝錯誤、`errors.Is`、`errors.As`、自訂錯誤型別
- 練習：
  1. 讓 `ParseAge` 拒絕大於 150 的數字，並新增哨兵錯誤 `ErrTooOld`。
  2. 寫一個讀取檔案的函式，檔案不存在時用 `errors.Is(err, os.ErrNotExist)` 判斷。

### 05 併發 `lessons/05-concurrency`
- 重點：`go` 關鍵字、unbuffered channel、`close` 與 `range`、`sync.WaitGroup`、`sync.Mutex`、`-race`
- 練習：
  1. 寫一個 worker pool：固定 3 個 goroutine 處理 10 個工作。
  2. 用 `select` 搭配 `time.After` 實作逾時。
  3. 把 `SafeCounter` 的 `Lock` 拿掉，跑 `go test -race` 看看會發生什麼。

## 接下來可以學什麼

- 標準函式庫：`net/http`（寫 Web API）、`encoding/json`、`os`、`io`、`context`
- 建立 CLI 工具或 REST API 小專案，並放進 `cmd/` 底下
- 推薦資源：
  - [A Tour of Go](https://go.dev/tour/)（有繁體中文版）
  - [Go by Example](https://gobyexample.com/)
  - [Effective Go](https://go.dev/doc/effective_go)
  - [Learn Go with Tests](https://quii.gitbook.io/learn-go-with-tests)
