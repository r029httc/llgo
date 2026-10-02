package concurrency

import (
	"context"
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

func TestWorkerPool(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	inputs := make([]int, 50)
	for i := range inputs {
		inputs[i] = i
	}

	var running, peak atomic.Int32
	square := func(n int) int {
		cur := running.Add(1)
		for {
			p := peak.Load()
			if cur <= p || peak.CompareAndSwap(p, cur) {
				break
			}
		}
		time.Sleep(time.Millisecond)
		running.Add(-1)
		return n * n
	}

	got := WorkerPool(inputs, 4, square)
	for i, v := range got {
		if v != i*i {
			t.Fatalf("got[%d] = %d, want %d（順序要與輸入相同）", i, v, i*i)
		}
	}
	if p := peak.Load(); p > 4 {
		t.Errorf("同時執行的 worker 最多 %d 個，超過 workers=4", p)
	}
	if p := peak.Load(); p < 2 {
		t.Errorf("同時執行的 worker 只有 %d 個，看起來沒有並行", p)
	}
	if got := WorkerPool(nil, 3, square); len(got) != 0 {
		t.Errorf("空輸入應回傳空切片，得到 %v", got)
	}
	if got := WorkerPool([]int{3}, 0, square); !slices.Equal(got, []int{9}) {
		t.Errorf("workers=0 應視為 1，得到 %v", got)
	}
}

func TestWithTimeout(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	v, err := WithTimeout(func() int { return 42 }, time.Second)
	if v != 42 || err != nil {
		t.Errorf("快速函式 = (%d, %v), want (42, nil)", v, err)
	}
	start := time.Now()
	_, err = WithTimeout(func() int { time.Sleep(500 * time.Millisecond); return 1 }, 20*time.Millisecond)
	if !errors.Is(err, ErrTimeout) {
		t.Errorf("慢速函式應回傳 ErrTimeout，得到 %v", err)
	}
	if el := time.Since(start); el > 300*time.Millisecond {
		t.Errorf("逾時後應立即返回，卻等了 %v", el)
	}
}

func TestFanIn(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	src := func(nums ...int) <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			for _, n := range nums {
				ch <- n
			}
		}()
		return ch
	}
	var got []int
	for v := range FanIn(src(1, 2, 3), src(10, 20), src(), src(100)) {
		got = append(got, v)
	}
	slices.Sort(got)
	if want := []int{1, 2, 3, 10, 20, 100}; !slices.Equal(got, want) {
		t.Errorf("FanIn = %v, want %v", got, want)
	}
	if _, ok := <-FanIn(); ok {
		t.Error("沒有輸入時，輸出 channel 應直接關閉")
	}
}

func TestGenerate(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	ctx, cancel := context.WithCancel(context.Background())
	ch := Generate(ctx, 10)
	var got []int
	for range 5 {
		got = append(got, <-ch)
	}
	if want := []int{10, 11, 12, 13, 14}; !slices.Equal(got, want) {
		t.Errorf("Generate = %v, want %v", got, want)
	}
	cancel()

	timeout := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return // 正確：取消後 channel 被關閉
			}
		case <-timeout:
			t.Fatal("cancel 後 1 秒內 channel 沒有關閉（goroutine 洩漏）")
		}
	}
}

func TestCache(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	c := NewCache[string, int]()
	if _, ok := c.Get("x"); ok {
		t.Error("空快取 Get 應回傳 ok=false")
	}
	c.Set("x", 1)
	if v, ok := c.Get("x"); !ok || v != 1 {
		t.Errorf("Get(x) = (%d, %v), want (1, true)", v, ok)
	}

	var calls atomic.Int32
	var wg sync.WaitGroup
	results := make([]int, 100)
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = c.GetOrCompute("slow", func() int {
				calls.Add(1)
				time.Sleep(10 * time.Millisecond)
				return 99
			})
		}()
	}
	wg.Wait()
	if n := calls.Load(); n != 1 {
		t.Errorf("compute 被呼叫 %d 次，應只有 1 次", n)
	}
	for i, v := range results {
		if v != 99 {
			t.Fatalf("results[%d] = %d, want 99", i, v)
		}
	}
	if v := c.GetOrCompute("x", func() int { return -1 }); v != 1 {
		t.Errorf("已存在的 key 不應重新計算，得到 %d", v)
	}
}
