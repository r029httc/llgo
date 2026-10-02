package collections

import (
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

func TestMap(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	got := Map([]int{1, 2, 3}, strconv.Itoa)
	if want := []string{"1", "2", "3"}; !slices.Equal(got, want) {
		t.Errorf("Map = %v, want %v", got, want)
	}
	lens := Map([]string{"go", "語言"}, func(s string) int { return len([]rune(s)) })
	if want := []int{2, 2}; !slices.Equal(lens, want) {
		t.Errorf("Map = %v, want %v", lens, want)
	}
	if got := Map(nil, strconv.Itoa); len(got) != 0 {
		t.Errorf("Map(nil) = %v, want empty", got)
	}
}

func TestUnique(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	got := Unique([]string{"b", "a", "b", "c", "a"})
	if want := []string{"b", "a", "c"}; !slices.Equal(got, want) {
		t.Errorf("Unique = %v, want %v", got, want)
	}
}

func TestChunk(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	got := Chunk([]int{1, 2, 3, 4, 5}, 2)
	if want := [][]int{{1, 2}, {3, 4}, {5}}; !reflect.DeepEqual(got, want) {
		t.Errorf("Chunk(.., 2) = %v, want %v", got, want)
	}
	if got := Chunk([]int{1, 2}, 5); !reflect.DeepEqual(got, [][]int{{1, 2}}) {
		t.Errorf("Chunk(.., 5) = %v", got)
	}
	if got := Chunk([]int{1, 2}, 0); got != nil {
		t.Errorf("Chunk(.., 0) = %v, want nil", got)
	}
	// 進階：對第一段 append 不應該覆蓋到第二段
	in := []int{1, 2, 3, 4}
	parts := Chunk(in, 2)
	_ = append(parts[0], 99)
	if parts[1][0] != 3 {
		t.Errorf("append 到第一段時覆蓋了第二段：%v（提示：三索引切片 s[a:b:c]）", parts)
	}
}

func TestGroupBy(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	words := []string{"apple", "avocado", "banana", "blueberry", "cherry"}
	got := GroupBy(words, func(s string) byte { return s[0] })
	want := map[byte][]string{
		'a': {"apple", "avocado"},
		'b': {"banana", "blueberry"},
		'c': {"cherry"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("GroupBy = %v, want %v", got, want)
	}
	byLen := GroupBy([]int{1, 22, 333, 4, 55}, func(n int) int { return len(strconv.Itoa(n)) })
	if !slices.Equal(byLen[2], []int{22, 55}) {
		t.Errorf("GroupBy by length = %v", byLen)
	}
}

func TestTopN(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	counts := map[string]int{"go": 5, "rust": 3, "java": 3, "c": 1, "zig": 5}
	tests := []struct {
		n    int
		want []string
	}{
		{1, []string{"go"}},
		{3, []string{"go", "zig", "java"}},
		{10, []string{"go", "zig", "java", "rust", "c"}},
		{0, []string{}},
	}
	for _, tt := range tests {
		if got := TopN(counts, tt.n); !slices.Equal(got, tt.want) {
			t.Errorf("TopN(n=%d) = %v, want %v", tt.n, got, tt.want)
		}
	}
}

func TestInsert(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	got := Insert([]int{1, 2, 3}, 1, 99)
	if want := []int{1, 99, 2, 3}; !slices.Equal(got, want) {
		t.Errorf("Insert middle = %v, want %v", got, want)
	}
	if got := Insert([]int{1, 2}, 0, 0); !slices.Equal(got, []int{0, 1, 2}) {
		t.Errorf("Insert front = %v", got)
	}
	if got := Insert([]int{1, 2}, 2, 3); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("Insert end = %v", got)
	}

	// 陷阱：s 的容量大於長度時，天真的 append 會改到原本的底層陣列
	s := make([]int, 3, 10)
	copy(s, []int{1, 2, 3})
	Insert(s, 1, 99)
	if full := s[:cap(s)]; !slices.Equal(full[:4], []int{1, 2, 3, 0}) {
		t.Errorf("Insert 修改了原切片的底層陣列：%v", full[:4])
	}
}

func TestStack(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	var s Stack[string] // 零值即可使用
	if _, ok := s.Pop(); ok {
		t.Error("空堆疊 Pop 應回傳 ok=false")
	}
	if _, ok := s.Peek(); ok {
		t.Error("空堆疊 Peek 應回傳 ok=false")
	}
	for _, w := range strings.Fields("a b c") {
		s.Push(w)
	}
	if s.Len() != 3 {
		t.Errorf("Len = %d, want 3", s.Len())
	}
	if v, _ := s.Peek(); v != "c" {
		t.Errorf("Peek = %q, want c", v)
	}
	var popped []string
	for s.Len() > 0 {
		v, _ := s.Pop()
		popped = append(popped, v)
	}
	if want := []string{"c", "b", "a"}; !slices.Equal(popped, want) {
		t.Errorf("Pop 順序 = %v, want %v", popped, want)
	}
}
