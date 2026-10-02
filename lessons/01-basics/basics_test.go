package basics

import (
	"fmt"
	"testing"
)

func TestGreet(t *testing.T) {
	if got := Greet("Go"); got != "Hello, Go!" {
		t.Errorf("Greet(\"Go\") = %q, want %q", got, "Hello, Go!")
	}
	if got := Greet(""); got != "Hello, World!" {
		t.Errorf("Greet(\"\") = %q, want %q", got, "Hello, World!")
	}
}

func TestDivmod(t *testing.T) {
	q, r := Divmod(17, 5)
	if q != 3 || r != 2 {
		t.Errorf("Divmod(17, 5) = (%d, %d), want (3, 2)", q, r)
	}
}

// 表格驅動測試（table-driven test）是 Go 社群最常見的測試寫法。
func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		in   int
		want string
	}{
		{1, "1"},
		{3, "Fizz"},
		{5, "Buzz"},
		{15, "FizzBuzz"},
		{22, "22"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprint(tt.in), func(t *testing.T) {
			if got := FizzBuzz(tt.in); got != tt.want {
				t.Errorf("FizzBuzz(%d) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// Example 函式既是文件也是測試：go test 會比對 Output 註解。
func ExampleSum() {
	fmt.Println(Sum(1, 2, 3, 4))
	// Output: 10
}
