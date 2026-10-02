//go:build !solution

// 練習 02：切片、映射與泛型。題目說明請看 README.md。
package collections

import "github.com/r029httc/llgo/exercises/internal/todo"

// 2.1 ★ Map 對每個元素套用 f，回傳新切片。
func Map[T, U any](s []T, f func(T) U) []U {
	panic(todo.NotImplemented)
}

// 2.2 ★ Unique 移除重複字串，保留第一次出現的順序。
func Unique(s []string) []string {
	panic(todo.NotImplemented)
}

// 2.3 ★★ Chunk 把切片切成每段 size 個元素；size <= 0 回傳 nil。
func Chunk[T any](s []T, size int) [][]T {
	panic(todo.NotImplemented)
}

// 2.4 ★★ GroupBy 依 key 函式把元素分組。
func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	panic(todo.NotImplemented)
}

// 2.5 ★★ TopN 回傳出現次數最多的 n 個字：次數由大到小，同次數依字母排序。
func TopN(counts map[string]int, n int) []string {
	panic(todo.NotImplemented)
}

// 2.6 ★★★ Insert 在索引 i 插入 v，回傳新切片；絕對不能修改 s 的底層陣列。
func Insert(s []int, i int, v int) []int {
	panic(todo.NotImplemented)
}

// 2.7 ★★★ Stack 是泛型堆疊，零值即可使用。
type Stack[T any] struct {
	// TODO: 加入你需要的欄位
}

func (s *Stack[T]) Push(v T)        { panic(todo.NotImplemented) }
func (s *Stack[T]) Pop() (T, bool)  { panic(todo.NotImplemented) }
func (s *Stack[T]) Peek() (T, bool) { panic(todo.NotImplemented) }
func (s *Stack[T]) Len() int        { panic(todo.NotImplemented) }
