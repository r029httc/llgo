//go:build !solution

// 練習 04：錯誤處理。題目說明請看 README.md。
package errs

import "github.com/r029httc/llgo/exercises/internal/todo"

// 4.1 ★ ParseAge 把字串轉成 0~150 的年齡。
func ParseAge(s string) (int, error) {
	panic(todo.NotImplemented)
}

// 4.2 ★★ ReadConfig 讀取 key=value 格式的設定檔。
func ReadConfig(path string) (map[string]string, error) {
	panic(todo.NotImplemented)
}

// 4.3 ★★ ValidateUser 一次回報所有驗證錯誤（errors.Join）。
func ValidateUser(u User) error {
	panic(todo.NotImplemented)
}

// 4.4 ★★★ Retry 最多呼叫 fn attempts 次，直到成功或遇到 ErrPermanent。
func Retry(attempts int, fn func() error) error {
	panic(todo.NotImplemented)
}

// 4.5 ★★★ Safely 執行 fn，把 fn 裡發生的 panic 轉成 error 回傳。
func Safely(fn func()) (err error) {
	panic(todo.NotImplemented)
}
