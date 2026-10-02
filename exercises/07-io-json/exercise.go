//go:build !solution

// 練習 07：標準函式庫 I/O、JSON、CSV。題目說明請看 README.md。
package iojson

import (
	"io"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 7.1 ★★ WordFrequency 讀取 r 的所有文字，統計每個單字（轉小寫、去除前後標點）出現次數。
func WordFrequency(r io.Reader) (map[string]int, error) {
	panic(todo.NotImplemented)
}

// 7.2 ★★ Write 讓 *CountingWriter 實作 io.Writer。
func (cw *CountingWriter) Write(p []byte) (int, error) {
	panic(todo.NotImplemented)
}

// 7.3 ★★ SaveTodos 把 todos 以縮排 2 個空白的 JSON 陣列寫入 w。
func SaveTodos(w io.Writer, todos []Todo) error {
	panic(todo.NotImplemented)
}

// 7.3 ★★★ LoadTodos 從 r 讀取 JSON 陣列；遇到未知欄位要回傳錯誤。
func LoadTodos(r io.Reader) ([]Todo, error) {
	panic(todo.NotImplemented)
}

// 7.4 ★★★ SumColumn 讀取含標題列的 CSV，加總指定欄位的數值。
func SumColumn(r io.Reader, column string) (float64, error) {
	panic(todo.NotImplemented)
}
