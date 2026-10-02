package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
	"golang.org/x/crypto/bcrypt"
)

func init() { BcryptCost = bcrypt.MinCost } // 測試時用最低成本，速度才快

var secret = []byte("test-secret-key-at-least-32-bytes!!")

// ---------- 11.1 ----------

func TestPassword(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h1, err := HashPassword("correct horse")
	if err != nil {
		t.Fatal(err)
	}
	h2, _ := HashPassword("correct horse")
	if h1 == h2 {
		t.Error("同一個密碼兩次雜湊結果應不同（bcrypt 會加隨機 salt）")
	}
	if strings.Contains(h1, "correct horse") {
		t.Error("雜湊結果不應包含明文密碼！")
	}
	if !CheckPassword(h1, "correct horse") || !CheckPassword(h2, "correct horse") {
		t.Error("正確密碼應通過驗證")
	}
	if CheckPassword(h1, "wrong horse") || CheckPassword(h1, "") {
		t.Error("錯誤密碼不應通過驗證")
	}
	if CheckPassword("not-a-hash", "correct horse") {
		t.Error("無效的雜湊不應通過驗證")
	}
	if _, err := HashPassword("短密碼"); !errors.Is(err, ErrWeakPassword) {
		t.Errorf("太短的密碼應 ErrWeakPassword，得到 %v", err)
	}
	if _, err := HashPassword("八個中文字的密碼"); err != nil {
		t.Errorf("8 個中文字應該合法（算的是字元數，不是 bytes），得到 %v", err)
	}
	if _, err := HashPassword(strings.Repeat("a", 100)); err == nil {
		t.Error("bcrypt 最多支援 72 bytes，超過應回傳錯誤")
	}
}

// ---------- 11.2 ----------

func TestRandomToken(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	seen := map[string]bool{}
	for range 100 {
		tok, err := RandomToken()
		if err != nil {
			t.Fatal(err)
		}
		raw, err := base64.RawURLEncoding.DecodeString(tok)
		if err != nil || len(raw) != 32 {
			t.Fatalf("token %q 應為 32 bytes 的 RawURLEncoding，解碼得到 %d bytes, %v", tok, len(raw), err)
		}
		if seen[tok] {
			t.Fatal("產生了重複的 token")
		}
		seen[tok] = true
	}
}

// ---------- 11.3 ----------

func TestSignVerify(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	c := Claims{UserID: 42, Role: "admin", ExpiresAt: now.Add(time.Hour).Unix()}
	tok, err := Sign(c, secret)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(tok, ".") != 1 {
		t.Fatalf("token 格式應為 payload.signature，得到 %q", tok)
	}
	got, err := Verify(tok, secret, now)
	if err != nil || got != c {
		t.Errorf("Verify = %+v, %v; want %+v", got, err, c)
	}

	if _, err := Verify(tok, secret, now.Add(2*time.Hour)); !errors.Is(err, ErrExpired) {
		t.Errorf("過期 token 應 ErrExpired，得到 %v", err)
	}
	if _, err := Verify(tok, []byte("other-secret"), now); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("用錯誤密鑰驗證應 ErrInvalidToken，得到 %v", err)
	}

	// 攻擊：把 payload 改成 role=admin、userID=1，但沿用原本的簽章
	user, _ := Sign(Claims{UserID: 7, Role: "user", ExpiresAt: now.Add(time.Hour).Unix()}, secret)
	_, sig, _ := strings.Cut(user, ".")
	forged := base64.RawURLEncoding.EncodeToString([]byte(`{"uid":1,"role":"admin","exp":9999999999}`)) + "." + sig
	if _, err := Verify(forged, secret, now); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("竄改過的 token 應 ErrInvalidToken，得到 %v", err)
	}

	for _, bad := range []string{"", "no-dot", ".", "a.b.c", "!!!." + sig} {
		if _, err := Verify(bad, secret, now); !errors.Is(err, ErrInvalidToken) {
			t.Errorf("Verify(%q) 應 ErrInvalidToken，得到 %v", bad, err)
		}
	}
}

// ---------- 11.4 ----------

type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time          { return c.t }
func (c *fakeClock) Advance(d time.Duration) { c.t = c.t.Add(d) }

func TestRateLimiter(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	clock := &fakeClock{time.Unix(0, 0)}
	l := NewRateLimiter(2, 3, clock.Now) // 每秒 2 個，最多累積 3 個

	for i := range 3 {
		if !l.Allow("1.1.1.1") {
			t.Fatalf("第 %d 次應允許（桶子一開始是滿的）", i+1)
		}
	}
	if l.Allow("1.1.1.1") {
		t.Error("用完 burst 後應拒絕")
	}
	if !l.Allow("2.2.2.2") {
		t.Error("不同 key 應該有各自的桶子")
	}

	clock.Advance(500 * time.Millisecond) // 補 1 個
	if !l.Allow("1.1.1.1") || l.Allow("1.1.1.1") {
		t.Error("0.5 秒後應恰好補充 1 個令牌")
	}

	clock.Advance(time.Hour) // 補很多，但最多到 burst
	allowed := 0
	for range 10 {
		if l.Allow("1.1.1.1") {
			allowed++
		}
	}
	if allowed != 3 {
		t.Errorf("長時間閒置後最多累積 burst=3 個，實際允許 %d 次", allowed)
	}
}

func TestRateLimiterConcurrent(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	clock := &fakeClock{time.Unix(0, 0)}
	l := NewRateLimiter(1, 50, clock.Now)
	l.Allow("warmup")

	var wg sync.WaitGroup
	var mu sync.Mutex
	allowed := 0
	for range 200 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("k") {
				mu.Lock()
				allowed++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if allowed != 50 {
		t.Errorf("200 個併發請求應恰好允許 50 個，實際 %d（記得加鎖，並用 -race 測試）", allowed)
	}
}

// ---------- 11.5 ----------

func whoami(w http.ResponseWriter, r *http.Request) {
	c, ok := ClaimsFrom(r.Context())
	if !ok {
		http.Error(w, "no claims", 500)
		return
	}
	io.WriteString(w, c.Role)
}

func call(h http.Handler, authHeader string) (int, string) {
	req := httptest.NewRequest("GET", "/", nil)
	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestAuthenticate(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := Authenticate(secret, http.HandlerFunc(whoami))
	valid, _ := Sign(Claims{UserID: 1, Role: "editor", ExpiresAt: time.Now().Add(time.Hour).Unix()}, secret)
	expired, _ := Sign(Claims{UserID: 1, Role: "editor", ExpiresAt: time.Now().Add(-time.Hour).Unix()}, secret)

	if code, body := call(h, "Bearer "+valid); code != 200 || body != "editor" {
		t.Errorf("有效 token = %d %q, want 200 editor", code, body)
	}
	for _, header := range []string{"", "Bearer ", "Basic abc", valid, "Bearer " + expired, "Bearer garbage"} {
		if code, body := call(h, header); code != http.StatusUnauthorized || body == "editor" {
			t.Errorf("Authorization=%q = %d %q, want 401", header, code, body)
		}
	}
	if _, ok := ClaimsFrom(context.Background()); ok {
		t.Error("空的 context 應回傳 ok=false")
	}
}

func TestRequireRole(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	admin, _ := Sign(Claims{UserID: 1, Role: "admin", ExpiresAt: time.Now().Add(time.Hour).Unix()}, secret)
	user, _ := Sign(Claims{UserID: 2, Role: "user", ExpiresAt: time.Now().Add(time.Hour).Unix()}, secret)
	h := Authenticate(secret, RequireRole("admin", http.HandlerFunc(whoami)))

	if code, _ := call(h, "Bearer "+admin); code != 200 {
		t.Errorf("admin = %d, want 200", code)
	}
	if code, _ := call(h, "Bearer "+user); code != http.StatusForbidden {
		t.Errorf("user = %d, want 403", code)
	}
	if code, _ := call(RequireRole("admin", http.HandlerFunc(whoami)), ""); code != http.StatusUnauthorized {
		t.Errorf("沒有經過 Authenticate = %d, want 401", code)
	}
}

// ---------- 11.6 ----------

func TestRateLimitMiddleware(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	clock := &fakeClock{time.Unix(0, 0)}
	h := RateLimit(NewRateLimiter(1, 2, clock.Now), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	}))
	from := func(addr string) int {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = addr
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		return rec.Code
	}
	// 同一個 IP、不同 port 要算同一個人
	codes := []int{from("10.0.0.1:1111"), from("10.0.0.1:2222"), from("10.0.0.1:3333")}
	if codes[0] != 200 || codes[1] != 200 || codes[2] != http.StatusTooManyRequests {
		t.Errorf("同 IP 連續 3 次 = %v, want [200 200 429]", codes)
	}
	if code := from("10.0.0.2:1111"); code != 200 {
		t.Errorf("其他 IP = %d, want 200", code)
	}
}
