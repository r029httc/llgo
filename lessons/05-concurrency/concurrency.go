// Package concurrency 介紹 goroutine、channel、sync.WaitGroup 與 sync.Mutex。
package concurrency

import "sync"

// ParallelSquare 為每個數字啟動一個 goroutine 計算平方，
// 並用 WaitGroup 等待全部完成。結果順序與輸入相同。
func ParallelSquare(nums []int) []int {
	out := make([]int, len(nums))
	var wg sync.WaitGroup
	for i, n := range nums {
		wg.Add(1)
		go func() { // go 關鍵字啟動一個 goroutine（輕量級執行緒）
			defer wg.Done()
			out[i] = n * n // 每個 goroutine 寫不同的索引，所以不會互相衝突
		}()
	}
	wg.Wait()
	return out
}

// Pipeline 示範 channel：generate -> square -> 收集。
// 「不要用共享記憶體來溝通；用溝通來共享記憶體。」
func Pipeline(nums ...int) []int {
	gen := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch) // 送完要 close，接收端的 range 才會結束
			for _, n := range nums {
				ch <- n
			}
		}()
		return ch
	}
	square := func(in <-chan int) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for n := range in {
				ch <- n * n
			}
		}()
		return ch
	}

	var result []int
	for v := range square(gen()) {
		result = append(result, v)
	}
	return result
}

// SafeCounter 用 Mutex 保護共享資料，避免資料競爭（data race）。
type SafeCounter struct {
	mu sync.Mutex
	m  map[string]int
}

func NewSafeCounter() *SafeCounter {
	return &SafeCounter{m: make(map[string]int)}
}

func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.m[key]
}
