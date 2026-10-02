//go:build solution

// 參考解答。建議自己寫完再看！
package testpractice

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func Encode(s string) string {
	r := []rune(s)
	var b strings.Builder
	for i := 0; i < len(r); {
		j := i
		for j < len(r) && r[j] == r[i] {
			j++
		}
		b.WriteString(strconv.Itoa(j - i))
		b.WriteRune(r[i])
		i = j
	}
	return b.String()
}

func Decode(s string) (string, error) {
	var b strings.Builder
	count := ""
	for _, c := range s {
		if unicode.IsDigit(c) {
			count += string(c)
			continue
		}
		if count == "" {
			return "", fmt.Errorf("字元 %q 前面缺少次數", c)
		}
		n, err := strconv.Atoi(count)
		if err != nil {
			return "", fmt.Errorf("次數 %q: %w", count, err)
		}
		b.WriteString(strings.Repeat(string(c), n))
		count = ""
	}
	if count != "" {
		return "", errors.New("結尾只有次數，缺少字元")
	}
	return b.String(), nil
}

func JoinPlus(parts []string, sep string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p // 每次 + 都會配置新字串並複製全部內容：O(n²)
	}
	return out
}

func JoinBuilder(parts []string, sep string) string {
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			b.WriteString(sep)
		}
		b.WriteString(p)
	}
	return b.String()
}

func Median(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	s := slices.Clone(nums) // 複製一份再排序，避免修改呼叫端資料
	slices.Sort(s)
	mid := len(s) / 2
	if len(s)%2 == 0 {
		return (s[mid-1] + s[mid]) / 2
	}
	return s[mid]
}
