//go:build solution

// 參考解答。建議自己寫完再看！
package concurrency

import (
	"context"
	"sync"
	"time"
)

func WorkerPool(inputs []int, workers int, fn func(int) int) []int {
	workers = max(workers, 1)
	out := make([]int, len(inputs))
	jobs := make(chan int) // 傳送「索引」，worker 依索引寫回結果，順序自然正確

	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				out[i] = fn(inputs[i])
			}
		}()
	}
	for i := range inputs {
		jobs <- i
	}
	close(jobs) // 讓 worker 的 range 結束
	wg.Wait()
	return out
}

func WithTimeout(fn func() int, d time.Duration) (int, error) {
	// 緩衝為 1：逾時後就算沒人接收，goroutine 也能送出並結束，不會洩漏
	done := make(chan int, 1)
	go func() { done <- fn() }()

	select {
	case v := <-done:
		return v, nil
	case <-time.After(d):
		return 0, ErrTimeout
	}
}

func FanIn(chs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range chs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func Generate(ctx context.Context, start int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := start; ; n++ {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

type entry[V any] struct {
	value V
	ready chan struct{} // 計算完成後 close
}

type Cache[K comparable, V any] struct {
	mu sync.Mutex
	m  map[K]*entry[V]
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{m: make(map[K]*entry[V])}
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	e, ok := c.m[key]
	c.mu.Unlock()
	if !ok {
		var zero V
		return zero, false
	}
	<-e.ready
	return e.value, true
}

func (c *Cache[K, V]) Set(key K, value V) {
	e := &entry[V]{value: value, ready: make(chan struct{})}
	close(e.ready)
	c.mu.Lock()
	c.m[key] = e
	c.mu.Unlock()
}

func (c *Cache[K, V]) GetOrCompute(key K, compute func() V) V {
	c.mu.Lock()
	if e, ok := c.m[key]; ok {
		c.mu.Unlock()
		<-e.ready // 別人正在算：等它算完
		return e.value
	}
	e := &entry[V]{ready: make(chan struct{})}
	c.m[key] = e
	c.mu.Unlock() // 計算時不持有鎖，其他 key 不會被卡住

	e.value = compute()
	close(e.ready)
	return e.value
}
