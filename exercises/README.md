# 練習題

共 **9 個主題、43 題**，每題都有自動化測試與參考解答。

## 怎麼做練習

每個練習目錄都有：

| 檔案 | 內容 |
| --- | --- |
| `README.md` | **詳細題目**：規格、範例、提示、觀念說明、思考題、延伸挑戰 |
| `exercise.go` | 你要編輯的檔案。每個函式目前都是 `panic(todo.NotImplemented)` |
| `types.go` | 題目提供的型別與錯誤（不用改） |
| `exercise_test.go` | 自動評分的測試 |
| `solution.go` | 參考解答（**先自己寫，卡關再看**） |

**流程**

1. 讀 `README.md` 的題目。
2. 打開 `exercise.go`，把 `panic(todo.NotImplemented)` 換成你的實作。
3. 執行測試：
   ```bash
   go test -v ./exercises/01-basics
   ```
   - `SKIP`：還沒寫（不算失敗）
   - `FAIL`：寫了但有錯，看錯誤訊息修正
   - `PASS`：完成！
4. 全部通過後，挑戰 README 裡的「延伸挑戰」。
5. 寫完一題就 commit 一次：`git commit -am "練習 1.3 Fibonacci"`。

**其他指令**

```bash
go test ./exercises/...                    # 全部練習（未完成的會 SKIP）
go test -v -run TestToRoman ./exercises/01-basics   # 只跑一題
go test -race ./exercises/05-concurrency   # 併發題必加 -race
go test -tags solution ./exercises/...     # 用參考解答跑測試
```

> `solution.go` 開頭的 `//go:build solution` 是**建置標籤**：只有加上 `-tags solution` 時才會編譯它，
> 同時 `exercise.go` 的 `//go:build !solution` 讓你的版本被排除。這也是 Go 處理不同平台程式碼（如 `_windows.go`）的機制。

## 難度

- ★ 熟悉語法，10~20 分鐘
- ★★ 需要理解觀念，20~60 分鐘
- ★★★ 有陷阱或需要設計，可能要 1 小時以上，卡住很正常

## 進度表

複製到你自己的筆記，完成一題打一個勾。

### [01 基礎語法](01-basics/README.md)
- [ ] 1.1 ★ `Max`：可變參數、多回傳值
- [ ] 1.2 ★ `IsPrime`：迴圈、邊界條件
- [ ] 1.3 ★★ `Fibonacci`：`make`、切片
- [ ] 1.4 ★★ `ReverseString`：`rune` 與 UTF-8
- [ ] 1.5 ★★ `IsPalindrome`：`unicode` 套件
- [ ] 1.6 ★★★ `ToRoman`：貪婪演算法、`strings.Builder`

### [02 切片、映射與泛型](02-collections/README.md)
- [ ] 2.1 ★ `Map`：泛型函式
- [ ] 2.2 ★ `Unique`：map 當集合
- [ ] 2.3 ★★ `Chunk`：共用底層陣列、三索引切片
- [ ] 2.4 ★★ `GroupBy`：`comparable` 約束
- [ ] 2.5 ★★ `TopN`：map 無序、多鍵排序
- [ ] 2.6 ★★★ `Insert`：`append` 的別名陷阱
- [ ] 2.7 ★★★ `Stack[T]`：泛型型別、零值可用

### [03 結構、方法與介面](03-structs-interfaces/README.md)
- [ ] 3.1 ★ `Triangle`：實作介面
- [ ] 3.2 ★★ `Largest` / `SortByArea`：介面切片
- [ ] 3.3 ★★ `Account` / `Transfer`：封裝、指標接收者
- [ ] 3.4 ★★ `Dog` / `Cat` / `Introduce`：嵌入
- [ ] 3.5 ★★★ `Describe`：type switch
- [ ] 3.6 ★★★ `Person` / `ByAge`：`fmt.Stringer`、`sort.Interface`

### [04 錯誤處理](04-errors/README.md)
- [ ] 4.1 ★ `ParseAge`：`%w` 包裝、哨兵錯誤
- [ ] 4.2 ★★ `ReadConfig`：檔案 I/O、自訂錯誤型別、`errors.As`
- [ ] 4.3 ★★ `ValidateUser`：`errors.Join`
- [ ] 4.4 ★★★ `Retry`：錯誤分類
- [ ] 4.5 ★★★ `Safely`：`panic` / `recover`

### [05 併發](05-concurrency/README.md)
- [ ] 5.1 ★★ `WorkerPool`：WaitGroup、channel
- [ ] 5.2 ★★ `WithTimeout`：`select`、goroutine 洩漏
- [ ] 5.3 ★★ `FanIn`：合併 channel
- [ ] 5.4 ★★★ `Generate`：`context` 取消
- [ ] 5.5 ★★★ `Cache`：細粒度鎖、只計算一次

### [06 測試、基準與模糊測試](06-testing/README.md)
- [ ] 6.1 ★★ `Encode` / `Decode`：fuzz 測試
- [ ] 6.2 ★ `JoinPlus` / `JoinBuilder`：benchmark
- [ ] 6.3 ★★ 找出 `Median` 的 bug：測試驅動除錯
- [ ] 6.4 自由練習：補強覆蓋率、Example、fuzz

### [07 I/O、JSON 與 CSV](07-io-json/README.md)
- [ ] 7.1 ★★ `WordFrequency`：`io.Reader`、`bufio.Scanner`
- [ ] 7.2 ★★ `CountingWriter`：實作 `io.Writer`
- [ ] 7.3 ★★ `SaveTodos` / ★★★ `LoadTodos`：JSON、struct tag
- [ ] 7.4 ★★★ `SumColumn`：`encoding/csv`

### [08 HTTP](08-http/README.md)
- [ ] 8.1 ★★★ `NewServer`：REST API、路由、併發安全
- [ ] 8.2 ★★ `RequireAPIKey`：中介層
- [ ] 8.3 ★★★ `FetchTodo`：HTTP 用戶端、`context` 逾時

### [09 期末專案](09-project/README.md)
- [ ] 專題 A：命令列待辦清單
- [ ] 專題 B：短網址服務
- [ ] 專題 C：併發網站健康檢查器

## 卡關時

1. 重讀錯誤訊息——Go 的測試訊息會告訴你輸入、實際結果、預期結果。
2. 在程式中加 `fmt.Printf("debug: %+v\n", x)`，用 `go test -v` 執行就看得到輸出。
3. 查官方文件：`go doc strings.Cut`、`go doc -all slices`，或 <https://pkg.go.dev>。
4. 看 README 的「提示」。
5. 還是不行再看 `solution.go`——看完後**刪掉自己的版本，憑記憶重寫一次**，這樣才真的學會。
