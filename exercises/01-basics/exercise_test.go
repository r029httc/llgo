package basics

import (
	"slices"
	"testing"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

func TestMax(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	tests := []struct {
		in     []int
		want   int
		wantOK bool
	}{
		{[]int{3, 7, 2}, 7, true},
		{[]int{-5, -2, -9}, -2, true}, // 全是負數：初始值不能設 0！
		{[]int{42}, 42, true},
		{nil, 0, false},
	}
	for _, tt := range tests {
		got, ok := Max(tt.in...)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("Max(%v) = (%d, %v), want (%d, %v)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}

func TestIsPrime(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	primes := []int{2, 3, 5, 7, 11, 13, 97, 7919}
	notPrimes := []int{-7, 0, 1, 4, 9, 15, 49, 7917}
	for _, n := range primes {
		if !IsPrime(n) {
			t.Errorf("IsPrime(%d) = false, want true", n)
		}
	}
	for _, n := range notPrimes {
		if IsPrime(n) {
			t.Errorf("IsPrime(%d) = true, want false", n)
		}
	}
}

func TestFibonacci(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if got := Fibonacci(0); len(got) != 0 {
		t.Errorf("Fibonacci(0) = %v, want empty", got)
	}
	if got := Fibonacci(-3); len(got) != 0 {
		t.Errorf("Fibonacci(-3) = %v, want empty", got)
	}
	if got, want := Fibonacci(1), []int{0}; !slices.Equal(got, want) {
		t.Errorf("Fibonacci(1) = %v, want %v", got, want)
	}
	if got, want := Fibonacci(10), []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}; !slices.Equal(got, want) {
		t.Errorf("Fibonacci(10) = %v, want %v", got, want)
	}
}

func TestReverseString(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	tests := map[string]string{
		"":      "",
		"a":     "a",
		"Hello": "olleH",
		"你好，世界": "界世，好你",
		"Go語言":  "言語oG",
		"🙂👍":    "👍🙂",
	}
	for in, want := range tests {
		if got := ReverseString(in); got != want {
			t.Errorf("ReverseString(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsPalindrome(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	yes := []string{"", "a", "racecar", "A man, a plan, a canal: Panama", "上海自來水來自海上", "No 'x' in Nixon"}
	no := []string{"ab", "Hello", "你好"}
	for _, s := range yes {
		if !IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = false, want true", s)
		}
	}
	for _, s := range no {
		if IsPalindrome(s) {
			t.Errorf("IsPalindrome(%q) = true, want false", s)
		}
	}
}

func TestToRoman(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	tests := map[int]string{
		0: "", 4000: "", -1: "",
		1: "I", 3: "III", 4: "IV", 9: "IX", 14: "XIV", 40: "XL",
		90: "XC", 400: "CD", 1994: "MCMXCIV", 2026: "MMXXVI", 3999: "MMMCMXCIX",
	}
	for in, want := range tests {
		if got := ToRoman(in); got != want {
			t.Errorf("ToRoman(%d) = %q, want %q", in, got, want)
		}
	}
}
