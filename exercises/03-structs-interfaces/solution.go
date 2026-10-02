//go:build solution

// 參考解答。建議自己寫完再看！
package shapes

import (
	"cmp"
	"fmt"
	"math"
	"slices"
)

func (t Triangle) Perimeter() float64 { return t.A + t.B + t.C }

func (t Triangle) Area() float64 {
	s := t.Perimeter() / 2
	return math.Sqrt(s * (s - t.A) * (s - t.B) * (s - t.C))
}

func Largest(shapes []Shape) (Shape, bool) {
	if len(shapes) == 0 {
		return nil, false
	}
	best := shapes[0]
	for _, s := range shapes[1:] {
		if s.Area() > best.Area() {
			best = s
		}
	}
	return best, true
}

func SortByArea(shapes []Shape) {
	slices.SortFunc(shapes, func(a, b Shape) int { return cmp.Compare(a.Area(), b.Area()) })
}

func NewAccount(owner string) *Account { return &Account{owner: owner} }

func (a *Account) Deposit(amount int) bool {
	if amount <= 0 {
		return false
	}
	a.balance += amount
	return true
}

func (a *Account) Withdraw(amount int) bool {
	if amount <= 0 || amount > a.balance {
		return false
	}
	a.balance -= amount
	return true
}

func (a *Account) Balance() int { return a.balance }

func Transfer(from, to *Account, amount int) bool {
	if from == to || !from.Withdraw(amount) {
		return false
	}
	to.Deposit(amount)
	return true
}

func (d Dog) Sound() string { return "汪汪" }
func (c Cat) Sound() string { return "喵" }

func Introduce(p Pet) string { return p.Hello() + "，" + p.Sound() }

func Describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int: %d", x)
	case string:
		return fmt.Sprintf("string: %q", x)
	case []int:
		return fmt.Sprintf("[]int len=%d", len(x))
	case error: // error 是介面，要放在更一般的 case 之前
		return "error: " + x.Error()
	case Shape:
		return fmt.Sprintf("shape area=%.2f", x.Area())
	default:
		return fmt.Sprintf("unknown: %T", x)
	}
}

func (p Person) String() string { return fmt.Sprintf("%s (%d)", p.Name, p.Age) }

func (a ByAge) Len() int      { return len(a) }
func (a ByAge) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool {
	if a[i].Age != a[j].Age {
		return a[i].Age < a[j].Age
	}
	return a[i].Name < a[j].Name
}
