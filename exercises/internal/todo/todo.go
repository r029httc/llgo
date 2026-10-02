// Package todo 提供練習題使用的輔助工具。
//
// 每個尚未完成的練習函式都會 panic(todo.NotImplemented)；
// 測試開頭的 defer todo.SkipIfNotImplemented(t) 會把這個 panic
// 轉成「跳過（SKIP）」，所以還沒寫的題目不會讓 go test 失敗。
// 你把 panic 換成自己的實作後，測試就會真的開始檢查。
package todo

import "testing"

type notImplemented struct{}

func (notImplemented) Error() string { return "TODO: 尚未實作" }

// NotImplemented 放在尚未完成的函式中：panic(todo.NotImplemented)。
var NotImplemented error = notImplemented{}

// SkipIfNotImplemented 必須用 defer 呼叫，放在測試函式的第一行。
func SkipIfNotImplemented(t testing.TB) {
	t.Helper()
	if r := recover(); r != nil {
		if r == NotImplemented {
			t.Skip("尚未實作：把 panic(todo.NotImplemented) 換成你的程式碼")
		}
		panic(r)
	}
}
