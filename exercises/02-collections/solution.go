//go:build solution

// 參考解答。建議自己寫完再看！
package collections

import (
	"cmp"
	"slices"
)

func Map[T, U any](s []T, f func(T) U) []U {
	out := make([]U, 0, len(s))
	for _, v := range s {
		out = append(out, f(v))
	}
	return out
}

func Unique(s []string) []string {
	seen := make(map[string]struct{}) // struct{} 不佔記憶體，常用來當 set
	var out []string
	for _, v := range s {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func Chunk[T any](s []T, size int) [][]T {
	if size <= 0 {
		return nil
	}
	var out [][]T
	for size < len(s) {
		// 三索引切片 s[low:high:max] 限制容量，避免呼叫端 append 時覆蓋到下一段
		out = append(out, s[:size:size])
		s = s[size:]
	}
	if len(s) > 0 {
		out = append(out, s)
	}
	return out
}

func GroupBy[T any, K comparable](s []T, key func(T) K) map[K][]T {
	out := make(map[K][]T)
	for _, v := range s {
		k := key(v)
		out[k] = append(out[k], v)
	}
	return out
}

func TopN(counts map[string]int, n int) []string {
	words := make([]string, 0, len(counts))
	for w := range counts {
		words = append(words, w)
	}
	slices.SortFunc(words, func(a, b string) int {
		if c := cmp.Compare(counts[b], counts[a]); c != 0 { // 次數大的在前
			return c
		}
		return cmp.Compare(a, b)
	})
	if n < len(words) {
		words = words[:max(n, 0)]
	}
	return words
}

func Insert(s []int, i int, v int) []int {
	out := make([]int, 0, len(s)+1) // 配置新的底層陣列
	out = append(out, s[:i]...)
	out = append(out, v)
	return append(out, s[i:]...)
}

type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
	v, ok := s.Peek()
	if ok {
		var zero T
		s.items[len(s.items)-1] = zero // 清掉引用，讓 GC 可以回收
		s.items = s.items[:len(s.items)-1]
	}
	return v, ok
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int { return len(s.items) }
