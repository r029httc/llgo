//go:build !solution

// 練習 01：基礎語法。把每個 panic(todo.NotImplemented) 換成你的實作。
// 題目說明請看同目錄的 README.md。
package basics

import "github.com/r029httc/llgo/exercises/internal/todo"

// 1.1 ★ Max 回傳 nums 中的最大值；nums 為空時 ok 為 false。
func Max(nums ...int) (m int, ok bool) {
	panic(todo.NotImplemented)
}

// 1.2 ★ IsPrime 判斷 n 是否為質數。
func IsPrime(n int) bool {
	panic(todo.NotImplemented)
}

// 1.3 ★★ Fibonacci 回傳前 n 個費氏數：0, 1, 1, 2, 3, 5, ...
func Fibonacci(n int) []int {
	panic(todo.NotImplemented)
}

// 1.4 ★★ ReverseString 反轉字串，必須正確處理中文等多位元組字元。
func ReverseString(s string) string {
	panic(todo.NotImplemented)
}

// 1.5 ★★ IsPalindrome 判斷是否為迴文，忽略大小寫、空白與標點。
func IsPalindrome(s string) bool {
	panic(todo.NotImplemented)
}

// 1.6 ★★★ ToRoman 把 1~3999 轉成羅馬數字；超出範圍回傳空字串。
func ToRoman(n int) string {
	panic(todo.NotImplemented)
}
