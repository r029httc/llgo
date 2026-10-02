//go:build !solution

// 練習 08：HTTP 伺服器與用戶端。題目說明請看 README.md。
package webapi

import (
	"context"
	"net/http"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 8.1 ★★★ NewServer 回傳一個提供 Todo REST API 的 http.Handler（資料存在記憶體中）。
func NewServer() http.Handler {
	panic(todo.NotImplemented)
}

// 8.2 ★★ RequireAPIKey 是中介層（middleware）：X-API-Key 標頭不等於 key 時回應 401。
func RequireAPIKey(key string, next http.Handler) http.Handler {
	panic(todo.NotImplemented)
}

// 8.3 ★★★ FetchTodo 呼叫 GET {baseURL}/todos/{id} 並解析回應。
func FetchTodo(ctx context.Context, client *http.Client, baseURL string, id int) (Todo, error) {
	panic(todo.NotImplemented)
}
