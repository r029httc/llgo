package webapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// do 對 handler 發出一個假請求，回傳狀態碼與 body（不需要真的開 port）。
func do(h http.Handler, method, target, body string) (int, string) {
	req := httptest.NewRequest(method, target, strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func TestHello(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := NewServer()
	if code, body := do(h, "GET", "/hello?name=Gopher", ""); code != 200 || body != "Hello, Gopher!" {
		t.Errorf("GET /hello?name=Gopher = %d %q", code, body)
	}
	if code, body := do(h, "GET", "/hello", ""); code != 200 || body != "Hello, World!" {
		t.Errorf("GET /hello = %d %q", code, body)
	}
}

func TestTodoAPI(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := NewServer()

	code, body := do(h, "GET", "/todos", "")
	if code != 200 || strings.TrimSpace(body) != "[]" {
		t.Fatalf("初始 GET /todos = %d %q, want 200 []（注意 nil 切片會編碼成 null）", code, body)
	}

	code, body = do(h, "POST", "/todos", `{"title": "學 HTTP"}`)
	if code != http.StatusCreated {
		t.Fatalf("POST /todos = %d %q, want 201", code, body)
	}
	var created Todo
	if err := json.Unmarshal([]byte(body), &created); err != nil || created != (Todo{ID: 1, Title: "學 HTTP"}) {
		t.Fatalf("POST 回應 = %q (%v), want {id:1 title:學 HTTP}", body, err)
	}
	do(h, "POST", "/todos", `{"title": "寫測試"}`)

	code, body = do(h, "GET", "/todos/2", "")
	var got Todo
	json.Unmarshal([]byte(body), &got)
	if code != 200 || got != (Todo{ID: 2, Title: "寫測試"}) {
		t.Errorf("GET /todos/2 = %d %q", code, body)
	}

	var list []Todo
	_, body = do(h, "GET", "/todos", "")
	if json.Unmarshal([]byte(body), &list); len(list) != 2 {
		t.Errorf("GET /todos 應有 2 筆，得到 %q", body)
	}

	tests := []struct {
		method, target, body string
		want                 int
	}{
		{"GET", "/todos/99", "", http.StatusNotFound},
		{"GET", "/todos/abc", "", http.StatusBadRequest},
		{"POST", "/todos", `{"title": ""}`, http.StatusBadRequest},
		{"POST", "/todos", `not json`, http.StatusBadRequest},
		{"DELETE", "/todos", "", http.StatusMethodNotAllowed},
	}
	for _, tt := range tests {
		if code, body := do(h, tt.method, tt.target, tt.body); code != tt.want {
			t.Errorf("%s %s %s = %d %q, want %d", tt.method, tt.target, tt.body, code, body, tt.want)
		}
	}
}

func TestTodoAPIConcurrent(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	h := NewServer()
	done := make(chan struct{})
	for range 50 {
		go func() {
			defer func() { done <- struct{}{} }()
			do(h, "POST", "/todos", `{"title": "x"}`)
			do(h, "GET", "/todos", "")
		}()
	}
	for range 50 {
		<-done
	}
	var list []Todo
	_, body := do(h, "GET", "/todos", "")
	json.Unmarshal([]byte(body), &list)
	seen := map[int]bool{}
	for _, td := range list {
		seen[td.ID] = true
	}
	if len(list) != 50 || len(seen) != 50 {
		t.Errorf("同時新增 50 筆後有 %d 筆、%d 個不同 ID，want 50/50（記得加鎖，並用 -race 測試）", len(list), len(seen))
	}
}

func TestRequireAPIKey(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "secret") })
	h := RequireAPIKey("s3cr3t", inner)

	for _, key := range []string{"", "wrong"} {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("X-API-Key", key)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized || strings.Contains(rec.Body.String(), "secret") {
			t.Errorf("key=%q: %d %q, want 401 且不洩漏內容", key, rec.Code, rec.Body.String())
		}
	}
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "s3cr3t")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != 200 || rec.Body.String() != "secret" {
		t.Errorf("正確 key: %d %q", rec.Code, rec.Body.String())
	}
}

func TestFetchTodo(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	// httptest.NewServer 會在本機開一個真正的 HTTP 伺服器
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/todos/1":
			io.WriteString(w, `{"id": 1, "title": "遠端資料", "done": true}`)
		case "/todos/500":
			http.Error(w, "boom", http.StatusInternalServerError)
		case "/todos/7":
			io.WriteString(w, `{broken`)
		case "/todos/8":
			select { // 模擬很慢的伺服器
			case <-time.After(2 * time.Second):
			case <-r.Context().Done():
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	ctx := context.Background()

	got, err := FetchTodo(ctx, srv.Client(), srv.URL, 1)
	if err != nil || got != (Todo{ID: 1, Title: "遠端資料", Done: true}) {
		t.Errorf("FetchTodo(1) = %+v, %v", got, err)
	}
	if _, err := FetchTodo(ctx, srv.Client(), srv.URL, 2); !errors.Is(err, ErrNotFound) {
		t.Errorf("404 應 errors.Is ErrNotFound，得到 %v", err)
	}
	if _, err := FetchTodo(ctx, srv.Client(), srv.URL, 500); err == nil || errors.Is(err, ErrNotFound) {
		t.Errorf("500 應回傳（非 ErrNotFound 的）錯誤，得到 %v", err)
	}
	if _, err := FetchTodo(ctx, srv.Client(), srv.URL, 7); err == nil {
		t.Error("無效 JSON 應回傳錯誤")
	}

	tctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	start := time.Now()
	_, err = FetchTodo(tctx, srv.Client(), srv.URL, 8)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("逾時應 errors.Is context.DeadlineExceeded，得到 %v", err)
	}
	if el := time.Since(start); el > time.Second {
		t.Errorf("應在 context 逾時後立即返回，卻花了 %v（用 http.NewRequestWithContext）", el)
	}
}
