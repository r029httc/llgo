# 練習 04：錯誤處理

> 對應課程：[`lessons/04-errors`](../../lessons/04-errors)
> 檔案：`exercise.go`（你要寫的）、`types.go`（已定義好的錯誤與型別）

Go 的錯誤處理哲學：**錯誤就是一般的值**。函式把 `error` 放在最後一個回傳值，呼叫端立刻檢查。

三個一定要分清楚的工具：

| 工具 | 用途 | 範例 |
| --- | --- | --- |
| `fmt.Errorf("...: %w", err)` | 包裝錯誤，附加上下文 | `fmt.Errorf("讀取 %s: %w", path, err)` |
| `errors.Is(err, target)` | 錯誤鏈中有沒有**某個值** | `errors.Is(err, os.ErrNotExist)` |
| `errors.As(err, &target)` | 錯誤鏈中有沒有**某個型別**，並取出來 | `var pe *ParseError; errors.As(err, &pe)` |

---

## 4.1 ★ `ParseAge(s string) (int, error)`

**規格**
- 合法範圍 0~150。
- 不是數字 → 包裝 `strconv` 的錯誤（測試用 `errors.As` 檢查 `*strconv.NumError`）。
- 負數 → 包裝 `ErrNegative`；大於 150 → 包裝 `ErrTooOld`。

**提示**：用 `%w` 而不是 `%v`。`%v` 只會把訊息變成字串，錯誤鏈就斷了，`errors.Is` 會找不到。

**思考題**：錯誤訊息要寫「解析年齡失敗」還是「解析年齡」？Go 慣例是錯誤訊息**小寫開頭、不加句號**，而且因為會被層層包裝，最好只描述「正在做什麼」。

---

## 4.2 ★★ `ReadConfig(path string) (map[string]string, error)`

**規格**：讀取這種格式的設定檔：

```ini
# 註解
host = localhost
port=8080
```

- 空白行與 `#` 開頭的行跳過。
- key 與 value 前後的空白要去掉；value 中可以再出現 `=`（只切第一個）。
- 沒有 `=` 或 key 為空 → 回傳 `*ParseError{Line: 行號, Text: 原始內容}`，行號從 1 開始，**空白行和註解也算行數**。
- 檔案不存在 → 回傳的錯誤要滿足 `errors.Is(err, os.ErrNotExist)`。

**提示**
- `os.Open` + `defer f.Close()`：確保無論怎麼 return 都會關閉檔案。
- `bufio.NewScanner(f)` 逐行讀取，結束後**一定要檢查** `sc.Err()`。
- `strings.Cut(s, "=")` 回傳 `before, after, found`，比 `strings.SplitN` 更好用。
- 測試用 `t.TempDir()` 建立暫存目錄，測試結束會自動刪除——你自己寫測試時也可以這樣做。

**學到的概念**：`defer`、檔案 I/O、自訂錯誤型別、`errors.As`

**延伸挑戰**：支援 `[section]` 區段，回傳 `map[string]map[string]string`。

---

## 4.3 ★★ `ValidateUser(u User) error`

**規格**：**一次回報所有錯誤**，而不是遇到第一個就停：
- `Name` 去掉空白後為空 → `ErrEmptyName`
- `Email` 不含 `@` → `ErrInvalidEmail`
- `Age` 小於 0 或大於 150 → `ErrInvalidAge`
- 全部合法 → `nil`

**提示**：`errors.Join(errs...)`（Go 1.20+）把多個錯誤合成一個，`errors.Is` 對其中任何一個都會成立。而且 `errors.Join()` 沒有參數時回傳 `nil`。

**學到的概念**：`errors.Join`、收集多個錯誤

---

## 4.4 ★★★ `Retry(attempts int, fn func() error) error`

**規格**
- 呼叫 `fn`，失敗就重試，最多 `attempts` 次；`attempts < 1` 時至少呼叫 1 次。
- 成功 → 立即回傳 `nil`。
- 錯誤鏈中有 `ErrPermanent` → **立即停止**並回傳該錯誤（重試也沒用的錯誤，例如密碼錯誤）。
- 次數用完 → 回傳包裝了最後一個錯誤的 error。

**學到的概念**：把函式當參數、閉包、錯誤分類（可重試 vs 不可重試）

**延伸挑戰**
1. 加入指數退避（exponential backoff）：每次失敗後等待 `base * 2^i`。
2. 接收 `context.Context`，在 ctx 取消時停止等待（學完練習 05 再做）。
3. 為了讓測試不必真的 sleep，把「等待」也做成一個可替換的參數——這叫**依賴注入**。

---

## 4.5 ★★★ `Safely(fn func()) (err error)`

**規格**：執行 `fn`；如果 `fn` panic，把 panic 轉成 error 回傳：
- panic 的值是 `error` → 包裝它（`errors.Is` 要能找到原錯誤）。
- 其他值（字串、runtime 錯誤等）→ 轉成描述它的 error。
- 沒有 panic → `nil`。

**提示**
```go
func Safely(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ... // 具名回傳值可以在 defer 中修改！
		}
	}()
	...
}
```
- `recover()` 只有在 `defer` 的函式中直接呼叫才有效。
- `r.(error)` 型別斷言判斷 panic 值是不是 error。

**觀念**：Go 的 panic **不是**用來做一般錯誤處理的。`recover` 只在少數場景使用，例如 HTTP 伺服器避免單一請求的 bug 把整個程式弄掛。`net/http` 內部就是這樣做的。

**學到的概念**：`panic`、`recover`、`defer` 與具名回傳值的互動
