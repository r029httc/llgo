# 練習 01：基礎語法

> 對應課程：[`lessons/01-basics`](../../lessons/01-basics)　｜　檔案：`exercise.go`

```bash
go test -v ./exercises/01-basics            # 執行全部
go test -v -run TestMax ./exercises/01-basics   # 只跑一題
```

---

## 1.1 ★ `Max(nums ...int) (m int, ok bool)`

**規格**
- 回傳 `nums` 中最大的數，`ok` 為 `true`。
- `nums` 為空（包括沒傳參數、傳 `nil`）時回傳 `(0, false)`。

| 輸入 | 輸出 |
| --- | --- |
| `Max(3, 7, 2)` | `7, true` |
| `Max(-5, -2, -9)` | `-2, true` |
| `Max()` | `0, false` |

**提示**
- 常見錯誤：把初始最大值設為 `0`，結果全是負數時答案錯了。用 `nums[0]` 當初始值。
- 可變參數 `nums ...int` 在函式內就是 `[]int`；呼叫端可以用 `Max(s...)` 展開切片。

**學到的概念**：可變參數、多回傳值、`comma ok` 慣例、`for range`

**延伸挑戰**：用泛型改寫成 `Max[T cmp.Ordered](nums ...T) (T, bool)`，讓它也能比較 `float64` 和 `string`。

---

## 1.2 ★ `IsPrime(n int) bool`

**規格**：`n` 是質數回傳 `true`。小於 2 的數（含負數、0、1）都不是質數。

**提示**
- 只需要檢查到 √n：`for i := 2; i*i <= n; i++`。
- 為什麼用 `i*i <= n` 而不是 `math.Sqrt`？可以避免浮點數誤差與型別轉換。

**學到的概念**：三段式 `for` 迴圈、`%` 取餘數、提早 `return`

**延伸挑戰**：寫 `PrimesUpTo(n int) []int`，用「埃拉托斯特尼篩法」（Sieve of Eratosthenes）找出所有 ≤ n 的質數，並和逐一呼叫 `IsPrime` 比較速度（學完練習 06 再用 benchmark 量）。

---

## 1.3 ★★ `Fibonacci(n int) []int`

**規格**：回傳前 `n` 個費氏數 `0, 1, 1, 2, 3, 5, 8, ...`；`n <= 0` 回傳空切片。

**提示**
- 先用 `make([]int, n)` 配置好長度，再用索引填值，比一直 `append` 更有效率。
- `out[i] = out[i-1] + out[i-2]`，前兩項特別處理。

**學到的概念**：`make`、切片索引、邊界條件

**延伸挑戰**
1. `Fibonacci(93)` 的最後一項會發生什麼事？（提示：`int64` 溢位）改用 `math/big` 處理大數。
2. 寫一個「閉包產生器」`FibGen() func() int`，每次呼叫回傳下一個費氏數。

---

## 1.4 ★★ `ReverseString(s string) string`

**規格**：反轉字串，必須正確處理中文、emoji 等多位元組字元。

| 輸入 | 輸出 |
| --- | --- |
| `"Hello"` | `"olleH"` |
| `"你好，世界"` | `"界世，好你"` |
| `"🙂👍"` | `"👍🙂"` |

**提示**
- Go 的 `string` 是**唯讀的 byte 序列**（UTF-8 編碼）。`len("你")` 是 3，不是 1！
- 逐 byte 反轉會把中文字拆壞，請先轉成 `[]rune`（一個 rune = 一個 Unicode 字元）。
- 雙指標交換：`r[i], r[j] = r[j], r[i]`。

**學到的概念**：`string` vs `[]byte` vs `[]rune`、UTF-8、多重賦值

**思考題**：`for i, c := range "你好"` 印出的 `i` 是 `0, 1` 還是 `0, 3`？為什麼？

---

## 1.5 ★★ `IsPalindrome(s string) bool`

**規格**：判斷是否為迴文，**只看字母與數字**，忽略大小寫、空白、標點。

| 輸入 | 輸出 |
| --- | --- |
| `"A man, a plan, a canal: Panama"` | `true` |
| `"上海自來水來自海上"` | `true` |
| `"Hello"` | `false` |

**提示**：`unicode.IsLetter`、`unicode.IsDigit`、`unicode.ToLower` 對中文也有效。

**學到的概念**：`unicode` 套件、以 `range` 逐字元走訪字串

**延伸挑戰**：不建立新切片，直接用兩個索引從字串兩端往中間走（需要用 `utf8.DecodeRuneInString` / `utf8.DecodeLastRuneInString`）。

---

## 1.6 ★★★ `ToRoman(n int) string`

**規格**：把 1~3999 轉成羅馬數字；超出範圍回傳 `""`。

| 1 | 4 | 9 | 14 | 40 | 90 | 400 | 1994 | 3999 |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| I | IV | IX | XIV | XL | XC | CD | MCMXCIV | MMMCMXCIX |

**提示**
- 準備兩個對應的切片：`1000, 900, 500, 400, 100, ...` 與 `"M", "CM", "D", "CD", "C", ...`（把 900、400 等「減法組合」也當作一個符號，問題就變簡單了）。
- 從大到小，能減就減、並附加符號。
- 用 `strings.Builder` 串接字串。

**學到的概念**：貪婪演算法、平行切片、`strings.Builder`

**延伸挑戰**：寫反向的 `FromRoman(s string) (int, error)`，並寫一個測試確認 1~3999 全部都能 `FromRoman(ToRoman(n)) == n`。
