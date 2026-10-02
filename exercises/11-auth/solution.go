//go:build solution

// 參考解答。建議自己寫完再看！
package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

// ---------- 11.1 ----------

func HashPassword(password string) (string, error) {
	if utf8.RuneCountInString(password) < 8 {
		return "", ErrWeakPassword
	}
	// bcrypt 會自動產生隨機 salt 並存進結果字串，所以同一個密碼每次雜湊結果都不同
	h, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("雜湊密碼: %w", err) // 例如超過 72 bytes 的 bcrypt.ErrPasswordTooLong
	}
	return string(h), nil
}

func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// ---------- 11.2 ----------

func RandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil { // crypto/rand，絕對不要用 math/rand 產生密鑰或 token
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// ---------- 11.3 ----------

var b64 = base64.RawURLEncoding

func sign(payload string, secret []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(payload))
	return b64.EncodeToString(mac.Sum(nil))
}

func Sign(c Claims, secret []byte) (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	payload := b64.EncodeToString(data)
	return payload + "." + sign(payload, secret), nil
}

func Verify(token string, secret []byte, now time.Time) (Claims, error) {
	payload, sig, ok := strings.Cut(token, ".")
	if !ok {
		return Claims{}, fmt.Errorf("格式錯誤: %w", ErrInvalidToken)
	}
	// hmac.Equal 是常數時間比較：用 == 比較的話，攻擊者可以從回應時間猜出簽章
	if !hmac.Equal([]byte(sig), []byte(sign(payload, secret))) {
		return Claims{}, fmt.Errorf("簽章不符: %w", ErrInvalidToken)
	}
	// 簽章正確之後才解析內容
	data, err := b64.DecodeString(payload)
	if err != nil {
		return Claims{}, fmt.Errorf("payload 編碼: %w", ErrInvalidToken)
	}
	var c Claims
	if err := json.Unmarshal(data, &c); err != nil {
		return Claims{}, fmt.Errorf("payload 內容: %w", ErrInvalidToken)
	}
	if now.Unix() >= c.ExpiresAt {
		return Claims{}, ErrExpired
	}
	return c, nil
}

// ---------- 11.4 ----------

type bucket struct {
	tokens float64
	last   time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	rate    float64
	burst   float64
	now     func() time.Time
	buckets map[string]*bucket
}

func NewRateLimiter(rate float64, burst int, now func() time.Time) *RateLimiter {
	if now == nil {
		now = time.Now
	}
	return &RateLimiter{rate: rate, burst: float64(burst), now: now, buckets: make(map[string]*bucket)}
}

func (l *RateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	b, ok := l.buckets[key]
	if !ok {
		b = &bucket{tokens: l.burst, last: now} // 新的 key 從滿桶開始
		l.buckets[key] = b
	}
	// 依照經過的時間補充令牌，但不超過桶子容量
	b.tokens = min(l.burst, b.tokens+now.Sub(b.last).Seconds()*l.rate)
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

// ---------- 11.5 ----------

// 自訂一個未匯出的型別當 context key，避免和其他套件的 key 衝突
type ctxKey struct{}

func Authenticate(secret []byte, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "需要登入", http.StatusUnauthorized)
			return
		}
		claims, err := Verify(token, secret, time.Now())
		if err != nil {
			w.Header().Set("WWW-Authenticate", "Bearer")
			http.Error(w, "token 無效或已過期", http.StatusUnauthorized) // 不要把詳細錯誤原因告訴客戶端
			return
		}
		ctx := context.WithValue(r.Context(), ctxKey{}, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFrom(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(ctxKey{}).(Claims)
	return c, ok
}

func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, ok := ClaimsFrom(r.Context())
		if !ok {
			http.Error(w, "需要登入", http.StatusUnauthorized) // 401：你是誰？
			return
		}
		if c.Role != role {
			http.Error(w, "權限不足", http.StatusForbidden) // 403：知道你是誰，但你不能做這件事
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ---------- 11.6 ----------

func RateLimit(l *RateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr) // RemoteAddr 是 "ip:port"
		if err != nil {
			ip = r.RemoteAddr
		}
		if !l.Allow(ip) {
			w.Header().Set("Retry-After", "1")
			http.Error(w, "請求太頻繁", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
