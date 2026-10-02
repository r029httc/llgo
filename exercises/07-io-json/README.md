# 練習 07：I/O、JSON 與 CSV

> 檔案：`exercise.go`（你要寫的）、`types.go`（已定義好的型別）

Go 標準函式庫最強大的設計之一，就是兩個只有一個方法的介面：

```go
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
```

檔案、網路連線、HTTP body、`bytes.Buffer`、`strings.Reader`、壓縮串流……全部都實作了它們。
所以本章的函式都接收 `io.Reader` / `io.Writer`，而不是檔名——**同一份程式碼可以處理檔案、網路、記憶體**，測試時也只要傳一個 `strings.NewReader("...")` 就好。

---

## 7.1 ★★ `WordFrequency(r io.Reader) (map[string]int, error)`

**規格**
- 以空白切出單字，轉成小寫，去掉**前後**的標點（`"(really)"` → `"really"`，`"--"` → 整個丟掉）。
- 讀取失敗時要回傳 error（測試會傳一個永遠失敗的 Reader）。

**提示**
- `bufio.NewScanner(r)` + `sc.Split(bufio.ScanWords)`。
- `strings.TrimFunc(w, func(c rune) bool { return !unicode.IsLetter(c) && !unicode.IsDigit(c) })`。
- `for sc.Scan() { ... }` 結束後，回傳 `sc.Err()`——這是最常被忘記的一步。

**延伸挑戰**：寫一個 `cmd/wordfreq` 程式，讀取命令列參數指定的檔案（或沒有參數時讀 `os.Stdin`），印出前 10 名。試試：`cat README.md | go run ./cmd/wordfreq`。

---

## 7.2 ★★ `CountingWriter`

**規格**：`CountingWriter{W: 底層 writer}` 把資料轉寫到 `W`，並把**實際寫入**的位元組數累加到 `N`。

**提示**：只有 3 行。重點是理解：只要實作了 `Write` 方法，`fmt.Fprintf`、`io.Copy`、`json.NewEncoder` 等所有接受 `io.Writer` 的函式都能使用它。這就是**裝飾器模式**（decorator）。

**延伸挑戰**
- 寫一個 `UpperWriter`，把寫入的內容轉成大寫再寫到底層。
- 用 `io.MultiWriter` 同時寫到檔案與螢幕；用 `io.TeeReader` 在讀取時同時計算 SHA-256（`crypto/sha256`）。

---

## 7.3 ★★ / ★★★ `SaveTodos` / `LoadTodos`

**規格**
- `SaveTodos`：輸出**縮排 2 個空白**的 JSON 陣列，結尾一個換行。空清單或 `nil` 要輸出 `[]` 而不是 `null`。
- `Due` 為 `nil` 時不輸出該欄位（`omitempty`）；有值時輸出 RFC 3339 格式（`time.Time` 預設就是）。
- `LoadTodos`：無效 JSON、型別錯誤（`"id": "one"`）、**未知欄位**（打錯字 `"titel"`）都要回傳錯誤。

**提示**
- `enc := json.NewEncoder(w); enc.SetIndent("", "  "); enc.Encode(v)`。
- `nil` 切片會編碼成 `null`，空切片 `[]Todo{}` 才會編碼成 `[]`。
- `dec := json.NewDecoder(r); dec.DisallowUnknownFields()`。

**觀念：struct tag**
```go
Due *time.Time `json:"due,omitempty"`
```
- `json:"due"`：JSON 欄位名稱。
- `omitempty`：零值時省略。為什麼用 `*time.Time` 而不是 `time.Time`？因為 `time.Time` 是結構，`omitempty` 對結構無效；指標的零值是 `nil`，就可以省略。
- **只有大寫開頭的欄位會被編碼**，小寫欄位會被 `encoding/json` 忽略。

---

## 7.4 ★★★ `SumColumn(r io.Reader, column string) (float64, error)`

**規格**
- 第一列是標題。找到名為 `column` 的欄，把每一列該欄的數值加總。
- 數值前後可能有空白，要先去掉。
- 找不到欄位 → 包裝 `ErrColumnNotFound`。
- 數值無法解析 → 包裝 `*strconv.NumError`，錯誤訊息要包含**列號**（標題是第 1 列）。

**提示**
- 用 `encoding/csv`，**不要**自己用 `strings.Split(line, ",")`——測試裡有 `"香蕉, 進口"` 這種含逗號的引號欄位。
- 逐列 `cr.Read()`，直到 `errors.Is(err, io.EOF)`。
- `slices.Index(header, column)` 找欄位位置。

**延伸挑戰**：寫 `ReadCSV[T any](r io.Reader, parse func(map[string]string) (T, error)) ([]T, error)`，把每一列轉成 `map[欄名]值` 交給 parse 函式。
