package net

import (
	"github.com/hiank/think/token"
	"context"
)

var netCtx context.Context

//Init 初始化net 包
func Init(ctx context.Context) {

	token.InitBuilder(ctx)

	netCtx = ctx
}