//go:build !solution

// 練習 05：併發。題目說明請看 README.md。
// 寫完務必用 go test -race 執行！
package concurrency

import (
	"context"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

// 5.1 ★★ WorkerPool 用固定 workers 個 goroutine 處理 inputs，結果順序與輸入相同。
func WorkerPool(inputs []int, workers int, fn func(int) int) []int {
	panic(todo.NotImplemented)
}

// 5.2 ★★ WithTimeout 執行 fn；超過 d 還沒完成就回傳 ErrTimeout。
func WithTimeout(fn func() int, d time.Duration) (int, error) {
	panic(todo.NotImplemented)
}

// 5.3 ★★ FanIn 把多個 channel 合併成一個；所有輸入都關閉後，輸出才關閉。
func FanIn(chs ...<-chan int) <-chan int {
	panic(todo.NotImplemented)
}

// 5.4 ★★★ Generate 從 start 開始不斷送出 start, start+1, ...，ctx 取消後關閉 channel。
func Generate(ctx context.Context, start int) <-chan int {
	panic(todo.NotImplemented)
}

// 5.5 ★★★ Cache 是可安全併發使用的快取。
type Cache[K comparable, V any] struct {
	// TODO: 加入你需要的欄位
}

func NewCache[K comparable, V any]() *Cache[K, V] { panic(todo.NotImplemented) }
func (c *Cache[K, V]) Get(key K) (V, bool)        { panic(todo.NotImplemented) }
func (c *Cache[K, V]) Set(key K, value V)         { panic(todo.NotImplemented) }

// GetOrCompute 回傳 key 的值；不存在時呼叫 compute 計算並存起來。
// 即使很多 goroutine 同時要同一個 key，compute 對每個 key 也只能被呼叫一次。
func (c *Cache[K, V]) GetOrCompute(key K, compute func() V) V { panic(todo.NotImplemented) }
