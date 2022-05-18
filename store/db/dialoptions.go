package db

import (
	"time"

	"github.com/hiank/think/doc"
	"github.com/hiank/think/run"
)

func DialOptions(opts ...DialOption) dialOptions {
	dopts := dialOptions{
		DB:          "0",
		DialTimeout: time.Second * 10,
		Coder:       doc.NewCoder[doc.PBCoder](), //default is protobuf coder
	}
	for _, opt := range opts {
		opt.Apply(&dopts)
	}
	return dopts
}

type dialOptions struct {
	DB          string
	Addr        string
	Account     string
	Password    string
	Coder       doc.Coder //default use protobuf coder
	DialTimeout time.Duration
}

type DialOption run.Option[*dialOptions]

func WithDB(dbname string) DialOption {
	return run.FuncOption[*dialOptions](func(do *dialOptions) {
		do.DB = dbname
	})
}

func WithAddr(addr string) DialOption {
	return run.FuncOption[*dialOptions](func(do *dialOptions) {
		do.Addr = addr
	})
}

func WithDailTimeout(timeout time.Duration) DialOption {
	return run.FuncOption[*dialOptions](func(do *dialOptions) {
		do.DialTimeout = timeout
	})
}

func WithAccount(account string) DialOption {
	return run.FuncOption[*dialOptions](func(do *dialOptions) {
		do.Account = account
	})
}

func WithPassword(password string) DialOption {
	return run.FuncOption[*dialOptions](func(do *dialOptions) {
		do.Password = password
	})
}

func WithCoder(coder doc.Coder) DialOption {
	return run.FuncOption[*dialOptions](func(do *dialOptions) {
		do.Coder = coder
	})
}
