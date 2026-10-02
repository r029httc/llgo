package concurrency

import (
	"reflect"
	"sync"
	"testing"
)

func TestParallelSquare(t *testing.T) {
	got := ParallelSquare([]int{1, 2, 3, 4})
	if want := []int{1, 4, 9, 16}; !reflect.DeepEqual(got, want) {
		t.Errorf("ParallelSquare = %v, want %v", got, want)
	}
}

func TestPipeline(t *testing.T) {
	got := Pipeline(2, 3, 4)
	if want := []int{4, 9, 16}; !reflect.DeepEqual(got, want) {
		t.Errorf("Pipeline = %v, want %v", got, want)
	}
}

// 用 go test -race ./... 執行，可以偵測資料競爭。
func TestSafeCounter(t *testing.T) {
	c := NewSafeCounter()
	var wg sync.WaitGroup
	for range 1000 { // Go 1.22+ 可以直接 range 整數
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("k")
		}()
	}
	wg.Wait()
	if got := c.Value("k"); got != 1000 {
		t.Errorf("Value = %d, want 1000", got)
	}
}
