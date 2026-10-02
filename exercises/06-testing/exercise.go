//go:build !solution

// 練習 06：測試、基準測試（benchmark）與模糊測試（fuzz）。題目說明請看 README.md。
package testpractice

import "github.com/r029httc/llgo/exercises/internal/todo"

// 6.1 ★★ Encode 做游程編碼（run-length encoding）："aaabcc" → "3a1b2c"。
func Encode(s string) string {
	panic(todo.NotImplemented)
}

// 6.1 ★★ Decode 是 Encode 的反函式；格式錯誤時回傳 error。
func Decode(s string) (string, error) {
	panic(todo.NotImplemented)
}

// 6.2 ★ JoinPlus 用 + 串接字串（不可使用 strings.Join）。
func JoinPlus(parts []string, sep string) string {
	panic(todo.NotImplemented)
}

// 6.2 ★ JoinBuilder 用 strings.Builder 串接字串（不可使用 strings.Join）。
func JoinBuilder(parts []string, sep string) string {
	panic(todo.NotImplemented)
}

// 6.3 ★★ Median 回傳中位數；空切片回傳 0，且不能修改呼叫端的切片。
// 這個實作「故意」有 bug！現有的測試抓不到它們——請先寫測試把 bug 抓出來，再修正。
func Median(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	return nums[len(nums)/2]
}
