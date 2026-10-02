package iojson

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/r029httc/llgo/exercises/internal/todo"
)

func TestWordFrequency(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	text := "Go is fun.\nGo, go, GO!   Is it (really) fun? -- yes"
	got, err := WordFrequency(strings.NewReader(text))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"go": 4, "is": 2, "fun": 2, "it": 1, "really": 1, "yes": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFrequency = %v, want %v", got, want)
	}
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("disk on fire") }

func TestWordFrequencyReadError(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	if _, err := WordFrequency(failingReader{}); err == nil {
		t.Error("讀取失敗時應回傳 error（提示：scanner.Err()）")
	}
}

func TestCountingWriter(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	var buf bytes.Buffer
	cw := &CountingWriter{W: &buf}
	var w io.Writer = cw // 編譯期確認 *CountingWriter 是 io.Writer
	fmt.Fprintf(w, "hello %s", "世界")
	io.WriteString(w, "!")
	if buf.String() != "hello 世界!" {
		t.Errorf("內容 = %q", buf.String())
	}
	if cw.N != int64(len("hello 世界!")) {
		t.Errorf("N = %d, want %d（中文一個字 3 bytes）", cw.N, len("hello 世界!"))
	}
}

func TestSaveTodos(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	var buf bytes.Buffer
	if err := SaveTodos(&buf, []Todo{{ID: 1, Title: "學 Go"}}); err != nil {
		t.Fatal(err)
	}
	want := `[
  {
    "id": 1,
    "title": "學 Go",
    "done": false
  }
]
`
	if buf.String() != want {
		t.Errorf("SaveTodos 輸出 =\n%s\nwant\n%s", buf.String(), want)
	}

	buf.Reset()
	SaveTodos(&buf, nil)
	if strings.TrimSpace(buf.String()) != "[]" {
		t.Errorf("空清單應輸出 []，得到 %q", buf.String())
	}
}

func TestLoadTodosRoundTrip(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	due := time.Date(2026, 12, 31, 18, 0, 0, 0, time.UTC)
	todos := []Todo{
		{ID: 1, Title: "寫練習", Done: true},
		{ID: 2, Title: "讀 Effective Go", Due: &due},
	}
	var buf bytes.Buffer
	if err := SaveTodos(&buf, todos); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"due": "2026-12-31T18:00:00Z"`) {
		t.Errorf("time.Time 應輸出 RFC 3339 格式：\n%s", buf.String())
	}
	got, err := LoadTodos(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != todos[0] || got[1].Due == nil || !got[1].Due.Equal(due) {
		t.Errorf("LoadTodos = %+v", got)
	}

	for _, bad := range []string{`not json`, `[{"id": "one"}]`, `[{"id": 1, "titel": "typo"}]`} {
		if _, err := LoadTodos(strings.NewReader(bad)); err == nil {
			t.Errorf("LoadTodos(%s) 應回傳錯誤", bad)
		}
	}
}

func TestSumColumn(t *testing.T) {
	defer todo.SkipIfNotImplemented(t)
	data := `name,price,qty
蘋果,30.5,2
"香蕉, 進口",20,5
芒果,  49.5 ,1
`
	got, err := SumColumn(strings.NewReader(data), "price")
	if err != nil || got != 100 {
		t.Errorf("SumColumn(price) = (%v, %v), want (100, nil)", got, err)
	}
	if got, _ := SumColumn(strings.NewReader(data), "qty"); got != 8 {
		t.Errorf("SumColumn(qty) = %v, want 8", got)
	}

	_, err = SumColumn(strings.NewReader(data), "weight")
	if !errors.Is(err, ErrColumnNotFound) {
		t.Errorf("找不到欄位應 errors.Is ErrColumnNotFound，得到 %v", err)
	}

	_, err = SumColumn(strings.NewReader("a,b\n1,2\nx,3\n"), "a")
	var numErr *strconv.NumError
	if !errors.As(err, &numErr) {
		t.Errorf("數值錯誤應包裝 *strconv.NumError，得到 %v", err)
	} else if !strings.Contains(err.Error(), "3") {
		t.Errorf("錯誤訊息應包含列號 3，得到 %v", err)
	}
}
