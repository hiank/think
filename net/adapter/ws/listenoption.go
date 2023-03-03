package ws

import (
	"github.com/hiank/think/auth"
	"github.com/hiank/think/auth/oauth"
	"github.com/hiank/think/run"
)

var (
	ErrUnimplementedAuther = run.Err("ws: unimplemented auther")
)

type ListenOption struct {
	Addr     string
	Auther   oauth.Auther
	Tokenset auth.Tokenset
}

type unimplementedAuther struct{}

func (*unimplementedAuther) Auth(string) (uint64, error) {
	return 0, ErrUnimplementedAuther
}

func withDefaultListenOption(opt ListenOption) ListenOption {
	if opt.Auther == nil {
		opt.Auther = &unimplementedAuther{}
	}
	return opt
}
