# 練習 06：測試、基準測試與模糊測試

> 檔案：`exercise.go`　｜　這一章沒有對應課程，README 本身就是教材。

Go 內建了完整的測試工具，不需要額外的框架：

| 函式簽名 | 種類 | 執行方式 |
| --- | --- | --- |
| `func TestXxx(t *testing.T)` | 單元測試 | `go test` |
| `func BenchmarkXxx(b *testing.B)` | 效能基準 | `go test -bench=.` |
| `func FuzzXxx(f *testing.F)` | 模糊測試 | `go test -fuzz=FuzzXxx` |
| `func ExampleXxx()` | 範例（同時是文件） | `go test` |

常用 flag：`-v` 詳細輸出、`-run 正規式` 篩選、`-count=1` 不使用快取、`-cover` 覆蓋率、`-coverprofile=c.out` + `go tool cover -html=c.out` 用瀏覽器看哪幾行沒被測到。

---

## 6.1 ★★ 游程編碼 `Encode` / `Decode`

**規格**
- `Encode("aaabcc")` → `"3a1b2c"`：每一段連續相同的字元寫成「次數 + 字元」。必須支援中文：`"好好好學"` → `"3好1學"`。
- `Decode` 是反函式。格式錯誤要回傳 error：`"a"`（缺次數）、`"3"`（結尾只有次數）。
- 次數可以超過一位數：`"12a"` → 12 個 a。

**模糊測試**：測試檔裡的 `FuzzRoundTrip` 檢查的是一個**性質**：對任何字串 s，`Decode(Encode(s)) == s`。

```bash
go test -run=^$ -fuzz=FuzzRoundTrip -fuzztime=30s ./exercises/06-testing
```

Go 會自動產生大量奇怪的輸入（空字串、無效 UTF-8、超長重複字元……）。發現失敗時，它會把那個輸入存到 `testdata/fuzz/FuzzRoundTrip/`，之後每次 `go test` 都會自動重跑——bug 就變成了回歸測試。

**思考題**：為什麼 fuzz 測試要跳過含有數字的字串？這說明了這種編碼格式有什麼限制？你能設計一種沒有這個限制的格式嗎？

---

## 6.2 ★ `JoinPlus` / `JoinBuilder` 與基準測試

**規格**：兩個函式都要做到和 `strings.Join(parts, sep)` 一樣的結果，但**不准使用 `strings.Join`**：
- `JoinPlus` 用 `out += ...` 串接。
- `JoinBuilder` 用 `strings.Builder`。

寫完後執行：

```bash
go test -bench=Join -benchmem ./exercises/06-testing
```

會看到類似這樣的結果（數字因機器而異）：

```
BenchmarkJoinPlus-4         1297     925693 ns/op   2881810 B/op   1798 allocs/op
BenchmarkJoinBuilder-4    152476       8036 ns/op      8440 B/op     11 allocs/op
```

**你的任務**：解釋為什麼差了 100 倍以上。（提示：Go 的字串不可變，每次 `+` 都要配置新記憶體並複製**全部**內容，總成本是 O(n²)。）

**延伸挑戰**：在 `JoinBuilder` 開頭先算出總長度並呼叫 `b.Grow(total)`，再跑一次 benchmark，`allocs/op` 變成多少？

---

## 6.3 ★★ 找出 `Median` 的 bug（測試驅動除錯）

`Median` 已經寫好了，而且現有的 `TestMedianWeak` 也通過了。**但它有兩個 bug。**

**你的任務**
1. 新增一個 `median_test.go`，用表格驅動測試寫出更好的測試案例，讓測試**失敗**。
2. 修正 `Median`，讓測試通過。
3. 確保修正後的 `Median` **不會修改**呼叫端傳入的切片（也為這點寫一個測試）。

**提示**
- 輸入沒有排序時會怎樣？`[]float64{3, 1, 2}`
- 偶數個元素時中位數怎麼算？`[]float64{1, 2, 3, 4}` 應該是 2.5。
- 如果你用 `sort.Float64s(nums)` 修 bug，呼叫端的資料會被打亂嗎？

**觀念**：「測試通過」不代表程式正確，只代表**你測的那些情況**是對的。先寫出會失敗的測試、再修 bug，是最可靠的除錯流程。

---

## 6.4 自由練習：替前面的練習補強測試

- 用 `go test -cover ./exercises/...` 看看各練習的覆蓋率。
- 替練習 01 的 `ToRoman` 寫一個 `ExampleToRoman()`，讓它出現在 `go doc` 裡。
- 替練習 02 的 `Unique` 寫一個 fuzz 測試，檢查三個性質：結果沒有重複、結果的每個元素都在輸入中、輸入的每個元素都在結果中。
- 讀標準函式庫 `strings` 套件的測試（`go env GOROOT` 底下的 `src/strings/strings_test.go`），看看 Go 團隊怎麼寫測試。
