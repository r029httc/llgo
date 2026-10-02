// 練習 03 共用的型別（題目已提供，不需修改）。
package shapes

import "math"

// Shape 是所有圖形的介面。
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rect struct{ Width, Height float64 }

func (r Rect) Area() float64      { return r.Width * r.Height }
func (r Rect) Perimeter() float64 { return 2 * (r.Width + r.Height) }

type Circle struct{ Radius float64 }

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// Triangle 三邊長為 A、B、C。（3.1 要幫它實作 Shape）
type Triangle struct{ A, B, C float64 }

// Account 是銀行帳戶。欄位小寫，外部套件無法直接修改餘額。（3.3）
type Account struct {
	owner   string
	balance int
}

// Animal 會被 Dog、Cat 嵌入（embedding）。（3.4）
type Animal struct{ Name string }

func (a Animal) Hello() string { return "我是" + a.Name }

type Dog struct {
	Animal // 嵌入：Dog 自動擁有 Name 欄位與 Hello 方法
	Breed  string
}

type Cat struct{ Animal }

// Pet 需要 Hello 與 Sound 兩個方法。
type Pet interface {
	Hello() string
	Sound() string
}

// Person 用於 3.6。
type Person struct {
	Name string
	Age  int
}

// ByAge 依年齡排序 Person（同年齡再依名字）。（3.6 要實作 sort.Interface）
type ByAge []Person
