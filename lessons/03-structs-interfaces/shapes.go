// Package shapes 介紹結構（struct）、方法（method）與介面（interface）。
package shapes

import (
	"fmt"
	"math"
)

// Shape 是一個介面：任何擁有 Area() 與 Perimeter() 方法的型別，
// 都「自動」實作了 Shape，不需要寫 implements。
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rect 是長方形。
type Rect struct {
	Width, Height float64
}

// Area 是 Rect 的方法。(r Rect) 稱為接收者（receiver）。
func (r Rect) Area() float64      { return r.Width * r.Height }
func (r Rect) Perimeter() float64 { return 2 * (r.Width + r.Height) }

// Circle 是圓形。
type Circle struct {
	Radius float64
}

func (c Circle) Area() float64      { return math.Pi * c.Radius * c.Radius }
func (c Circle) Perimeter() float64 { return 2 * math.Pi * c.Radius }

// String 讓 Circle 實作 fmt.Stringer，fmt.Println 時會使用它。
func (c Circle) String() string { return fmt.Sprintf("Circle(r=%.1f)", c.Radius) }

// TotalArea 接受任何 Shape，這就是介面帶來的多型。
func TotalArea(shapes ...Shape) float64 {
	total := 0.0
	for _, s := range shapes {
		total += s.Area()
	}
	return total
}

// Counter 示範「指標接收者」：要修改結構內容時必須用 *Counter。
type Counter struct {
	n int
}

func (c *Counter) Inc()       { c.n++ }
func (c *Counter) Value() int { return c.n }
