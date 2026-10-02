// 練習 05 共用的錯誤（題目已提供，不需修改）。
package concurrency

import "errors"

// ErrTimeout 由 WithTimeout 在逾時時回傳。
var ErrTimeout = errors.New("逾時")
