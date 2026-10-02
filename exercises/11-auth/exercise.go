//go:build !solution

// 練習 11：認證與安全（密碼雜湊、安全亂數、簽章 token、限流、認證中介層）。
// 題目說明請看 README.md。
package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 11.1 ★ HashPassword 用 bcrypt 雜湊密碼；少於 8 個字元回傳 ErrWeakPassword。
func HashPassword(password string) (string, error) {
	panic(todo.NotImplemented)
}

// 11.1 ★ CheckPassword 檢查密碼是否與雜湊相符。
func CheckPassword(hash, password string) bool {
	panic(todo.NotImplemented)
}

// 11.2 ★ RandomToken 產生 32 bytes 的密碼學安全亂數，以 base64 URL-safe（無補齊）編碼。
func RandomToken() (string, error) {
	panic(todo.NotImplemented)
}

// 11.3 ★★★ Sign 產生「payload.signature」格式的 token（HMAC-SHA256 簽章）。
func Sign(c Claims, secret []byte) (string, error) {
	panic(todo.NotImplemented)
}

// 11.3 ★★★ Verify 驗證 token 的簽章與期限，回傳其中的 Claims。
func Verify(token string, secret []byte, now time.Time) (Claims, error) {
	panic(todo.NotImplemented)
}

// 11.4 ★★ RateLimiter 以「令牌桶」演算法對每個 key（例如 IP）分別限流。
type RateLimiter struct {
	// TODO: 加入你需要的欄位
}

// NewRateLimiter：每秒補充 rate 個令牌，桶子最多 burst 個；now 是可替換的時鐘（方便測試）。
func NewRateLimiter(rate float64, burst int, now func() time.Time) *RateLimiter {
	panic(todo.NotImplemented)
}

func (l *RateLimiter) Allow(key string) bool {
	panic(todo.NotImplemented)
}

// 11.5 ★★★ Authenticate 驗證 Authorization: Bearer <token>，並把 Claims 放進 request 的 context。
func Authenticate(secret []byte, next http.Handler) http.Handler {
	panic(todo.NotImplemented)
}

// ClaimsFrom 從 context 取出 Authenticate 放入的 Claims。
func ClaimsFrom(ctx context.Context) (Claims, bool) {
	panic(todo.NotImplemented)
}

// RequireRole 只允許特定角色；必須包在 Authenticate 裡面使用。
func RequireRole(role string, next http.Handler) http.Handler {
	panic(todo.NotImplemented)
}

// 11.6 ★★ RateLimit 依客戶端 IP 限流，超過時回應 429。
func RateLimit(l *RateLimiter, next http.Handler) http.Handler {
	panic(todo.NotImplemented)
}
