package testpractice

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

func TestEncodeDecode(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	tests := map[string]string{
		"":             "",
		"a":            "1a",
		"aaabcc":       "3a1b2c",
		"aaaaaaaaaaaa": "12a",
		"好好好學":         "3好1學",
		"abab":         "1a1b1a1b",
	}
	for in, want := range tests {
		if got := Encode(in); got != want {
			t.Errorf("Encode(%q) = %q, want %q", in, got, want)
		}
		if got, err := Decode(want); err != nil || got != in {
			t.Errorf("Decode(%q) = (%q, %v), want (%q, nil)", want, got, err, in)
		}
	}
	for _, bad := range []string{"a", "3", "2a3", "x2y"} {
		if _, err := Decode(bad); err == nil {
			t.Errorf("Decode(%q) 應回傳錯誤", bad)
		}
	}
}

// FuzzRoundTrip 是模糊測試：
//   - go test 時只會跑下面的種子（seed corpus），當作一般測試。
//   - go test -fuzz=FuzzRoundTrip 會自動產生大量隨機輸入找 bug。
func FuzzRoundTrip(f *testing.F) {
	for _, seed := range []string{"", "a", "aaabcc", "好好學", "🙂🙂x"} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, s string) {
		defer todo.SkipIfNotImplemented(t)
		if !utf8.ValidString(s) || strings.IndexFunc(s, unicode.IsDigit) >= 0 {
			t.Skip("此格式無法表示含數字或無效 UTF-8 的字串")
		}
		got, err := Decode(Encode(s))
		if err != nil || got != s {
			t.Errorf("Decode(Encode(%q)) = (%q, %v)", s, got, err)
		}
	})
}

func TestJoin(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	cases := [][]string{nil, {"a"}, {"a", "b", "c"}, {"", ""}}
	for _, parts := range cases {
		want := strings.Join(parts, ", ")
		if got := JoinPlus(parts, ", "); got != want {
			t.Errorf("JoinPlus(%q) = %q, want %q", parts, got, want)
		}
		if got := JoinBuilder(parts, ", "); got != want {
			t.Errorf("JoinBuilder(%q) = %q, want %q", parts, got, want)
		}
	}
}

var benchParts = strings.Fields(strings.Repeat("go is fun ", 300))

// 執行：go test -bench=Join -benchmem ./exercises/06-testing
func BenchmarkJoinPlus(b *testing.B) {
	defer todo.SkipIfNotImplemented(b)
	for b.Loop() {
		JoinPlus(benchParts, ",")
	}
}

func BenchmarkJoinBuilder(b *testing.B) {
	defer todo.SkipIfNotImplemented(b)
	for b.Loop() {
		JoinBuilder(benchParts, ",")
	}
}

// 這些測試太弱，抓不到 Median 的 bug。
// 你的任務：在 median_test.go 新增測試讓它失敗，然後修好 Median。
func TestMedianWeak(t *testing.T) {
	if got := Median(nil); got != 0 {
		t.Errorf("Median(nil) = %v, want 0", got)
	}
	if got := Median([]float64{1, 2, 3}); got != 2 {
		t.Errorf("Median([1 2 3]) = %v, want 2", got)
	}
}
