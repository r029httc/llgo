//go:build solution

// 參考解答。建議自己寫完再看！
package webapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type server struct {
	mu     sync.Mutex // HTTP handler 會被多個 goroutine 同時呼叫，必須加鎖
	todos  []Todo
	nextID int
}

func NewServer() http.Handler {
	s := &server{nextID: 1}
	mux := http.NewServeMux()
	// Go 1.22+ 的路由可以指定 HTTP 方法與路徑參數 {id}
	mux.HandleFunc("GET /hello", s.hello)
	mux.HandleFunc("GET /todos", s.list)
	mux.HandleFunc("POST /todos", s.create)
	mux.HandleFunc("GET /todos/{id}", s.get)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *server) hello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "World"
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Hello, %s!", name)
}

func (s *server) list(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	todos := append([]Todo{}, s.todos...) // 複製一份，解鎖後再編碼
	s.mu.Unlock()
	writeJSON(w, http.StatusOK, todos)
}

func (s *server) create(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "無效的 JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Title) == "" {
		http.Error(w, "title 不可為空", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	t := Todo{ID: s.nextID, Title: in.Title}
	s.nextID++
	s.todos = append(s.todos, t)
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, t)
}

func (s *server) get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id 必須是整數", http.StatusBadRequest)
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.todos {
		if t.ID == id {
			writeJSON(w, http.StatusOK, t)
			return
		}
	}
	http.NotFound(w, r)
}

func RequireAPIKey(key string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != key {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func FetchTodo(ctx context.Context, client *http.Client, baseURL string, id int) (Todo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/todos/%d", baseURL, id), nil)
	if err != nil {
		return Todo{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return Todo{}, fmt.Errorf("取得 todo %d: %w", id, err)
	}
	defer resp.Body.Close() // 一定要關閉 Body，否則連線無法重用

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return Todo{}, fmt.Errorf("todo %d: %w", id, ErrNotFound)
	default:
		return Todo{}, fmt.Errorf("todo %d: 非預期的狀態碼 %d", id, resp.StatusCode)
	}

	var t Todo
	if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
		return Todo{}, fmt.Errorf("解析 todo %d: %w", id, err)
	}
	return t, nil
}
