// 練習 11 共用的型別與錯誤（題目已提供，不需修改）。
package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// BcryptCost 是 bcrypt 的計算成本。越高越安全也越慢（每 +1 慢一倍）。
// 測試時會調成 bcrypt.MinCost 加快速度；正式環境請用預設值或更高。
var BcryptCost = bcrypt.DefaultCost

// Claims 是 token 內攜帶的資訊。
type Claims struct {
	UserID    int64  `json:"uid"`
	Role      string `json:"role"`
	ExpiresAt int64  `json:"exp"` // Unix 秒數
}

var (
	ErrWeakPassword = errors.New("密碼至少需要 8 個字元")
	ErrInvalidToken = errors.New("token 無效")
	ErrExpired      = errors.New("token 已過期")
)
