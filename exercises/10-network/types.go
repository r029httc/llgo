// 練習 10 共用的型別與錯誤（題目已提供，不需修改）。
package chat

import "errors"

// Message 是發布到 Hub 的訊息。
type Message struct {
	Topic string
	Text  string
}

// MaxFrameSize 是單一封包允許的最大長度（1 MiB）。
const MaxFrameSize = 1 << 20

// ErrFrameTooLarge 表示封包超過 MaxFrameSize。
var ErrFrameTooLarge = errors.New("封包太大")
