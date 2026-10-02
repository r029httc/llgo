//go:build solution

// 參考解答。建議自己寫完再看！
package iojson

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func WordFrequency(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)
	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords) // 以空白切割單字
	for sc.Scan() {
		w := strings.TrimFunc(sc.Text(), func(c rune) bool {
			return !unicode.IsLetter(c) && !unicode.IsDigit(c)
		})
		if w != "" {
			counts[strings.ToLower(w)]++
		}
	}
	return counts, sc.Err()
}

func (cw *CountingWriter) Write(p []byte) (int, error) {
	n, err := cw.W.Write(p)
	cw.N += int64(n) // 只算真正寫出去的位元組
	return n, err
}

func SaveTodos(w io.Writer, todos []Todo) error {
	if todos == nil {
		todos = []Todo{} // 讓空清單輸出 [] 而不是 null
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(todos)
}

func LoadTodos(r io.Reader) ([]Todo, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var todos []Todo
	if err := dec.Decode(&todos); err != nil {
		return nil, fmt.Errorf("解析 todos: %w", err)
	}
	return todos, nil
}

func SumColumn(r io.Reader, column string) (float64, error) {
	cr := csv.NewReader(r)
	header, err := cr.Read()
	if err != nil {
		return 0, fmt.Errorf("讀取標題列: %w", err)
	}
	idx := slices.Index(header, column)
	if idx < 0 {
		return 0, fmt.Errorf("%q: %w", column, ErrColumnNotFound)
	}
	var sum float64
	for row := 2; ; row++ {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			return sum, nil
		}
		if err != nil {
			return 0, err
		}
		v, err := strconv.ParseFloat(strings.TrimSpace(rec[idx]), 64)
		if err != nil {
			return 0, fmt.Errorf("第 %d 列: %w", row, err)
		}
		sum += v
	}
}
