package shapes

import (
	"fmt"
	"math"
	"testing"
)

func TestTotalArea(t *testing.T) {
	got := TotalArea(Rect{Width: 2, Height: 3}, Circle{Radius: 1})
	want := 6 + math.Pi
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("TotalArea = %v, want %v", got, want)
	}
}

func TestCounter(t *testing.T) {
	var c Counter // 零值即可使用
	c.Inc()
	c.Inc()
	if c.Value() != 2 {
		t.Errorf("Value = %d, want 2", c.Value())
	}
}

func ExampleCircle_String() {
	fmt.Println(Circle{Radius: 2})
	// Output: Circle(r=2.0)
}
