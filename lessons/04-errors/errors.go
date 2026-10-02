// Package errs 介紹 Go 的錯誤處理：error 是一般的回傳值，而不是例外（exception）。
package errs

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrNegative 是一個「哨兵錯誤」（sentinel error），呼叫端可以用 errors.Is 比對。
var ErrNegative = errors.New("數值不可為負")

// ParseAge 把字串轉成年齡。慣例：error 放在最後一個回傳值。
func ParseAge(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		// %w 會「包裝」原始錯誤，保留錯誤鏈
		return 0, fmt.Errorf("解析年齡 %q 失敗: %w", s, err)
	}
	if n < 0 {
		return 0, fmt.Errorf("年齡 %d: %w", n, ErrNegative)
	}
	return n, nil
}

// ValidationError 是自訂錯誤型別，只要實作 Error() string 就是 error。
type ValidationError struct {
	Field string
}

func (e *ValidationError) Error() string {
	return "欄位驗證失敗: " + e.Field
}

// ValidateName 在名稱為空時回傳 *ValidationError。
func ValidateName(name string) error {
	if name == "" {
		return &ValidationError{Field: "name"}
	}
	return nil
}
