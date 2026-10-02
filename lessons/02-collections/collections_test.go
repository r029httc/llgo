package collections

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	in := []int{1, 2, 3}
	got := Reverse(in)
	if want := []int{3, 2, 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("Reverse = %v, want %v", got, want)
	}
	if in[0] != 1 {
		t.Error("Reverse 不應修改原切片")
	}
}

func TestFilter(t *testing.T) {
	even := Filter([]int{1, 2, 3, 4, 5, 6}, func(n int) bool { return n%2 == 0 })
	if want := []int{2, 4, 6}; !reflect.DeepEqual(even, want) {
		t.Errorf("Filter = %v, want %v", even, want)
	}
}

func TestWordCount(t *testing.T) {
	got := WordCount([]string{"go", "is", "fun", "go"})
	if got["go"] != 2 || got["fun"] != 1 {
		t.Errorf("WordCount = %v", got)
	}
	if keys := SortedKeys(got); !reflect.DeepEqual(keys, []string{"fun", "go", "is"}) {
		t.Errorf("SortedKeys = %v", keys)
	}
}
