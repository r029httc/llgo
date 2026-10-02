// 練習 04 共用的錯誤與型別（題目已提供，不需修改）。
package errs

import (
	"errors"
	"fmt"
)

// 哨兵錯誤：呼叫端用 errors.Is(err, ErrXxx) 判斷。
var (
	ErrNegative     = errors.New("數值不可為負")
	ErrTooOld       = errors.New("年齡超過 150")
	ErrEmptyName    = errors.New("名稱不可為空")
	ErrInvalidEmail = errors.New("email 格式錯誤")
	ErrInvalidAge   = errors.New("年齡不合理")
	ErrPermanent    = errors.New("永久性錯誤，不應重試")
)

// ParseError 表示設定檔某一行格式錯誤。用 errors.As 取出。
type ParseError struct {
	Line int    // 從 1 開始的行號
	Text string // 該行原始內容
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("第 %d 行格式錯誤: %q", e.Line, e.Text)
}

// User 用於 4.3。
type User struct {
	Name  string
	Email string
	Age   int
}
