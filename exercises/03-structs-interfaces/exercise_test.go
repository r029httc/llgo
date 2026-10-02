package shapes

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 編譯期檢查：若 Triangle 沒有實作 Shape，這行會編譯失敗。
var _ Shape = Triangle{}

func almostEqual(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestTriangle(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	tri := Triangle{3, 4, 5}
	if !almostEqual(tri.Area(), 6) {
		t.Errorf("Area = %v, want 6", tri.Area())
	}
	if !almostEqual(tri.Perimeter(), 12) {
		t.Errorf("Perimeter = %v, want 12", tri.Perimeter())
	}
	eq := Triangle{2, 2, 2}
	if !almostEqual(eq.Area(), math.Sqrt(3)) {
		t.Errorf("正三角形 Area = %v, want %v", eq.Area(), math.Sqrt(3))
	}
}

func TestLargestAndSort(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if _, ok := Largest(nil); ok {
		t.Error("Largest(nil) 應回傳 ok=false")
	}
	shapes := []Shape{Rect{10, 10}, Circle{1}, Rect{1, 2}, Circle{10}}
	if got, _ := Largest(shapes); got != (Circle{10}) {
		t.Errorf("Largest = %v, want Circle{10}", got)
	}
	SortByArea(shapes)
	want := []Shape{Rect{1, 2}, Circle{1}, Rect{10, 10}, Circle{10}}
	for i := range want {
		if shapes[i] != want[i] {
			t.Fatalf("SortByArea = %v, want %v", shapes, want)
		}
	}
}

func TestAccount(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	a := NewAccount("小明")
	if a.Balance() != 0 {
		t.Fatalf("新帳戶餘額 = %d, want 0", a.Balance())
	}
	steps := []struct {
		op      string
		amount  int
		wantOK  bool
		balance int
	}{
		{"deposit", 100, true, 100},
		{"deposit", -5, false, 100},
		{"deposit", 0, false, 100},
		{"withdraw", 30, true, 70},
		{"withdraw", 100, false, 70},
		{"withdraw", -1, false, 70},
		{"withdraw", 70, true, 0},
	}
	for i, s := range steps {
		var ok bool
		if s.op == "deposit" {
			ok = a.Deposit(s.amount)
		} else {
			ok = a.Withdraw(s.amount)
		}
		if ok != s.wantOK || a.Balance() != s.balance {
			t.Errorf("步驟 %d %s(%d) = %v, 餘額 %d; want %v, %d", i, s.op, s.amount, ok, a.Balance(), s.wantOK, s.balance)
		}
	}

	x, y := NewAccount("x"), NewAccount("y")
	x.Deposit(50)
	if !Transfer(x, y, 20) || x.Balance() != 30 || y.Balance() != 20 {
		t.Errorf("Transfer 成功後 x=%d y=%d, want 30 20", x.Balance(), y.Balance())
	}
	if Transfer(x, y, 999) || x.Balance() != 30 || y.Balance() != 20 {
		t.Errorf("餘額不足的 Transfer 不應改變任何餘額：x=%d y=%d", x.Balance(), y.Balance())
	}
	if Transfer(x, x, 10) || x.Balance() != 30 {
		t.Errorf("轉給自己應回傳 false")
	}
}

func TestIntroduce(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	d := Dog{Animal: Animal{Name: "小黑"}, Breed: "柴犬"}
	c := Cat{Animal{Name: "咪咪"}}
	if got, want := Introduce(d), "我是小黑，汪汪"; got != want {
		t.Errorf("Introduce(dog) = %q, want %q", got, want)
	}
	if got, want := Introduce(c), "我是咪咪，喵"; got != want {
		t.Errorf("Introduce(cat) = %q, want %q", got, want)
	}
}

func TestDescribe(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	tests := []struct {
		in   any
		want string
	}{
		{nil, "nil"},
		{42, "int: 42"},
		{"hi", `string: "hi"`},
		{[]int{1, 2, 3}, "[]int len=3"},
		{errors.New("boom"), "error: boom"},
		{Rect{2, 3}, "shape area=6.00"},
		{3.14, "unknown: float64"},
		{map[string]int{}, "unknown: map[string]int"},
	}
	for _, tt := range tests {
		if got := Describe(tt.in); got != tt.want {
			t.Errorf("Describe(%#v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestPersonSort(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if got := (Person{"Alice", 30}).String(); got != "Alice (30)" {
		t.Errorf("String() = %q, want %q", got, "Alice (30)")
	}
	people := []Person{{"Carol", 35}, {"Bob", 25}, {"Alice", 35}, {"Dave", 20}}
	sort.Sort(ByAge(people))
	if got, want := fmt.Sprint(people), "[Dave (20) Bob (25) Alice (35) Carol (35)]"; got != want {
		t.Errorf("排序後 = %s, want %s", got, want)
	}
}
