# 練習 03：結構、方法與介面

> 對應課程：[`lessons/03-structs-interfaces`](../../lessons/03-structs-interfaces)
> 檔案：`exercise.go`（你要寫的）、`types.go`（題目提供的型別，不用改）

---

## 3.1 ★ `Triangle` 實作 `Shape`

**規格**：`Triangle{A, B, C}` 是三邊長。實作 `Area()` 與 `Perimeter()`。

**提示**：海龍公式——令 `s = (A+B+C)/2`，面積 = `√(s(s-A)(s-B)(s-C))`。

**觀念**：測試檔裡有一行 `var _ Shape = Triangle{}`，這是 Go 慣用的**編譯期介面檢查**：只要少實作一個方法，就會直接編譯失敗。

**延伸挑戰**：加上 `NewTriangle(a, b, c float64) (Triangle, error)`，在三邊無法構成三角形時回傳錯誤（任兩邊之和必須大於第三邊）。

---

## 3.2 ★★ `Largest` 與 `SortByArea`

**規格**
- `Largest(shapes []Shape) (Shape, bool)`：面積最大的圖形；空切片回傳 `(nil, false)`。
- `SortByArea(shapes []Shape)`：依面積由小到大**原地**排序（不回傳新切片）。

**提示**
- 函式收到的切片和呼叫端共用底層陣列，所以原地排序後呼叫端看得到結果。
- `slices.SortFunc` + `cmp.Compare(a.Area(), b.Area())`。

**思考題**：測試中用 `got != (Circle{10})` 直接比較兩個介面值。介面值什麼時候可以用 `==` 比較？如果放進去的是含有切片欄位的結構會怎樣？

---

## 3.3 ★★ 銀行帳戶 `Account`

**規格**
| 方法 | 行為 |
| --- | --- |
| `NewAccount(owner) *Account` | 建立餘額 0 的帳戶 |
| `Deposit(amount) bool` | `amount <= 0` 回傳 `false`；否則加到餘額 |
| `Withdraw(amount) bool` | `amount <= 0` 或餘額不足回傳 `false`；否則扣款 |
| `Balance() int` | 目前餘額 |
| `Transfer(from, to, amount) bool` | 從 from 轉到 to；失敗時**兩邊都不能改變**；轉給自己回傳 `false` |

**提示**
- `balance` 是小寫欄位：其他套件無法直接改它，只能透過方法——這就是 Go 的**封裝**。
- 失敗時不能改變狀態：先檢查，再修改。
- `from == to` 比較的是兩個指標是否指向同一個帳戶。

**學到的概念**：建構函式慣例 `NewXxx`、指標接收者、封裝、不變量（invariant）

**延伸挑戰**：如果兩個 goroutine 同時對同一帳戶 `Withdraw`，會發生什麼事？學完練習 05 後，用 `sync.Mutex` 讓 `Account` 可以安全併發使用，並寫一個 `-race` 測試證明它。（`Transfer` 同時鎖兩個帳戶時要小心**死結**！）

---

## 3.4 ★★ 嵌入（embedding）：`Dog`、`Cat` 與 `Pet`

**規格**
- `Dog.Sound()` 回傳 `"汪汪"`，`Cat.Sound()` 回傳 `"喵"`。
- `Introduce(p Pet) string` 回傳 `p.Hello() + "，" + p.Sound()`。

```go
Introduce(Dog{Animal: Animal{Name: "小黑"}})   // "我是小黑，汪汪"
```

**觀念**
- `Dog` 嵌入了 `Animal`，所以 `Dog` 自動擁有 `Name` 欄位和 `Hello()` 方法（「提升」promotion）。
- 加上你寫的 `Sound()` 後，`Dog` 就同時擁有 `Hello` 與 `Sound`，因此**自動**滿足 `Pet` 介面。

**思考題（很重要）**：如果在 `Animal` 上寫 `func (a Animal) Intro() string { return a.Hello() + a.Sound() }` 並期待它呼叫 `Dog.Sound()`——這行得通嗎？嵌入和其他語言的「繼承」差在哪裡？（提示：Go 沒有虛擬方法。）

---

## 3.5 ★★★ `Describe(v any) string`（type switch）

**規格**

| 輸入 | 輸出 |
| --- | --- |
| `nil` | `nil` |
| `42` | `int: 42` |
| `"hi"` | `string: "hi"`（含引號，用 `%q`） |
| `[]int{1,2,3}` | `[]int len=3` |
| 任何 `error` | `error: <訊息>` |
| 任何 `Shape` | `shape area=6.00`（小數兩位） |
| 其他 | `unknown: <型別>`（用 `%T`） |

**提示**
```go
switch x := v.(type) {
case int:
	// 這裡 x 的型別是 int
case error:
	// 這裡 x 的型別是 error
}
```
- `case` 可以是具體型別，也可以是**介面**。
- 順序很重要：一個值可能同時滿足多個介面，第一個符合的 `case` 勝出。

**學到的概念**：`any`、type switch、型別斷言、`%q` `%T` 格式化動詞

---

## 3.6 ★★★ `fmt.Stringer` 與 `sort.Interface`

**規格**
- `Person.String()` 回傳 `"Alice (30)"` 格式。實作後 `fmt.Println(p)` 會自動使用它。
- `ByAge` 實作 `sort.Interface`（`Len`、`Less`、`Swap`）：依年齡由小到大，同年齡依名字。

```go
people := []Person{{"Carol", 35}, {"Bob", 25}, {"Alice", 35}}
sort.Sort(ByAge(people))
fmt.Println(people)   // [Bob (25) Alice (35) Carol (35)]
```

**觀念**：`ByAge(people)` 是**型別轉換**，不會複製資料；`ByAge` 和 `[]Person` 共用同一個底層陣列，所以排序後 `people` 也變了。

**學到的概念**：標準函式庫的小介面設計、具名切片型別、型別轉換

**延伸挑戰**：現代 Go 更常用 `slices.SortFunc`。用它改寫一次，比較兩種寫法的優缺點。
