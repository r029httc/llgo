// Package collections 介紹陣列（array）、切片（slice）與映射（map）。
package collections

import "sort"

// Reverse 回傳一個反轉後的新切片，不修改原本的切片。
func Reverse(s []int) []int {
	out := make([]int, len(s)) // make 建立指定長度的切片
	for i, v := range s {
		out[len(s)-1-i] = v
	}
	return out
}

// Filter 保留讓 keep 回傳 true 的元素。
// 這裡用到泛型（Go 1.18+）：[T any] 表示 T 可以是任何型別。
func Filter[T any](s []T, keep func(T) bool) []T {
	var out []T // nil 切片也可以直接 append
	for _, v := range s {
		if keep(v) {
			out = append(out, v)
		}
	}
	return out
}

// WordCount 計算每個單字出現的次數，示範 map 的使用。
func WordCount(words []string) map[string]int {
	counts := make(map[string]int)
	for _, w := range words {
		counts[w]++ // 不存在的 key 會得到零值 0
	}
	return counts
}

// SortedKeys 回傳 map 的 key 並排序。
// 注意：Go 的 map 迭代順序是隨機的，需要固定順序時要自己排序。
func SortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m)) // 長度 0、容量 len(m)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
