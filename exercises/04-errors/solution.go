//go:build solution

// 參考解答。建議自己寫完再看！
package errs

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseAge(s string) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, fmt.Errorf("解析年齡 %q: %w", s, err)
	}
	switch {
	case n < 0:
		return 0, fmt.Errorf("年齡 %d: %w", n, ErrNegative)
	case n > 150:
		return 0, fmt.Errorf("年齡 %d: %w", n, ErrTooOld)
	}
	return n, nil
}

func ReadConfig(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("讀取設定檔: %w", err)
	}
	defer f.Close() // 確保函式結束時關閉檔案

	cfg := make(map[string]string)
	sc := bufio.NewScanner(f)
	for line := 1; sc.Scan(); line++ {
		text := strings.TrimSpace(sc.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		key, value, ok := strings.Cut(text, "=")
		key = strings.TrimSpace(key)
		if !ok || key == "" {
			return nil, &ParseError{Line: line, Text: sc.Text()}
		}
		cfg[key] = strings.TrimSpace(value)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("讀取設定檔: %w", err)
	}
	return cfg, nil
}

func ValidateUser(u User) error {
	var errs []error
	if strings.TrimSpace(u.Name) == "" {
		errs = append(errs, ErrEmptyName)
	}
	if !strings.Contains(u.Email, "@") {
		errs = append(errs, fmt.Errorf("%q: %w", u.Email, ErrInvalidEmail))
	}
	if u.Age < 0 || u.Age > 150 {
		errs = append(errs, fmt.Errorf("%d: %w", u.Age, ErrInvalidAge))
	}
	return errors.Join(errs...) // 沒有錯誤時回傳 nil
}

func Retry(attempts int, fn func() error) error {
	attempts = max(attempts, 1)
	var err error
	for i := 1; i <= attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		if errors.Is(err, ErrPermanent) {
			return err
		}
	}
	return fmt.Errorf("重試 %d 次後仍失敗: %w", attempts, err)
}

func Safely(fn func()) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		// 具名回傳值 err 可以在 defer 中修改
		if e, ok := r.(error); ok {
			err = fmt.Errorf("panic: %w", e)
		} else {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	fn()
	return nil
}
