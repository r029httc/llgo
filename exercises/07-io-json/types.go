// 練習 07 共用的型別（題目已提供，不需修改）。
package iojson

import (
	"errors"
	"io"
	"time"
)

// Todo 的 struct tag 決定 JSON 欄位名稱。
type Todo struct {
	ID    int        `json:"id"`
	Title string     `json:"title"`
	Done  bool       `json:"done"`
	Due   *time.Time `json:"due,omitempty"` // 指標 + omitempty：nil 時不輸出
}

// CountingWriter 包裝另一個 io.Writer，並記錄寫入的總位元組數。（7.2）
type CountingWriter struct {
	W io.Writer
	N int64
}

// ErrColumnNotFound 由 SumColumn 在找不到欄位時回傳。（7.4）
var ErrColumnNotFound = errors.New("找不到欄位")
