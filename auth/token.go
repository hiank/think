package auth

import (
	"context"
)

//Token 令牌，用于管理持有者的活动状态
type Token struct {
	context.Context
	Invalidate context.CancelFunc
	value      string //NOTE: token 的string值
}

//newToken 构建一个token
func newToken(ctx context.Context, tkStr string) *Token {

	tk := &Token{
		value: tkStr,
	}
	tk.Context, tk.Invalidate = context.WithCancel(ctx)
	return tk
}

//ToString token string 字串
func (tk *Token) ToString() string {

	return tk.value
}

//Derive 派生的Token
func (tk *Token) Derive() *Token {

	return newToken(tk.Context, tk.value)
}
