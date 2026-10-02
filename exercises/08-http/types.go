// 練習 08 共用的型別（題目已提供，不需修改）。
package webapi

import "errors"

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// ErrNotFound 由 FetchTodo 在伺服器回應 404 時回傳。（8.3）
var ErrNotFound = errors.New("找不到資源")
