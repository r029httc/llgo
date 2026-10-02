//go:build !solution

// 練習 03：結構、方法與介面。題目說明請看 README.md。
package shapes

import "github.com/r029httc/llgo/exercises/internal/todo"

// 3.1 ★ 讓 Triangle 實作 Shape。面積請用海龍公式。
func (t Triangle) Area() float64      { panic(todo.NotImplemented) }
func (t Triangle) Perimeter() float64 { panic(todo.NotImplemented) }

// 3.2 ★★ Largest 回傳面積最大的圖形；SortByArea 依面積由小到大「原地」排序。
func Largest(shapes []Shape) (Shape, bool) { panic(todo.NotImplemented) }
func SortByArea(shapes []Shape)            { panic(todo.NotImplemented) }

// 3.3 ★★ 銀行帳戶。金額 <= 0 或餘額不足時回傳 false 且不改變餘額。
func NewAccount(owner string) *Account            { panic(todo.NotImplemented) }
func (a *Account) Deposit(amount int) bool        { panic(todo.NotImplemented) }
func (a *Account) Withdraw(amount int) bool       { panic(todo.NotImplemented) }
func (a *Account) Balance() int                   { panic(todo.NotImplemented) }
func Transfer(from, to *Account, amount int) bool { panic(todo.NotImplemented) }

// 3.4 ★★ 讓 Dog 和 Cat 實作 Pet，並寫出 Introduce。
func (d Dog) Sound() string  { panic(todo.NotImplemented) }
func (c Cat) Sound() string  { panic(todo.NotImplemented) }
func Introduce(p Pet) string { panic(todo.NotImplemented) }

// 3.5 ★★★ Describe 用 type switch 描述任意值。
func Describe(v any) string { panic(todo.NotImplemented) }

// 3.6 ★★★ Person 實作 fmt.Stringer；ByAge 實作 sort.Interface。
func (p Person) String() string { panic(todo.NotImplemented) }

func (a ByAge) Len() int           { panic(todo.NotImplemented) }
func (a ByAge) Less(i, j int) bool { panic(todo.NotImplemented) }
func (a ByAge) Swap(i, j int)      { panic(todo.NotImplemented) }
