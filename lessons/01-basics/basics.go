// Package basics 介紹 Go 的基本語法：變數、常數、型別、函式與流程控制。
package basics

import "fmt"

// Pi 是常數。常數在編譯期就決定，不能被修改。
const Pi = 3.14159

// Greet 回傳打招呼字串。
// 名稱以大寫開頭的識別字（Greet）會被「匯出」，其他套件可以使用；
// 小寫開頭（例如 helper）只能在本套件內使用。
func Greet(name string) string {
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

// Add 示範有多個參數的函式。相同型別的參數可以合併寫成 a, b int。
func Add(a, b int) int {
	return a + b
}

// Divmod 示範 Go 的「多回傳值」。
func Divmod(a, b int) (quotient, remainder int) {
	quotient = a / b
	remainder = a % b
	return // 具名回傳值可以直接 return
}

// Sum 示範 for 迴圈與可變參數（...int）。
// Go 只有 for 一種迴圈關鍵字。
func Sum(nums ...int) int {
	total := 0 // := 是「宣告並賦值」的簡寫，型別由右邊推斷
	for _, n := range nums {
		total += n
	}
	return total
}

// FizzBuzz 示範 switch。Go 的 switch 不需要 break，預設不會往下穿透。
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0:
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return fmt.Sprint(n)
	}
}
