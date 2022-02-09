package db

import "time"

func DialOptions(opts ...DialOption) dialOptions {
	dopts := dialOptions{
		DB:          "0",
		DialTimeout: time.Second * 10,
	}
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	return dopts
}

type dialOptions struct {
	DB          string
	Addr        string
	Account     string
	Password    string
	DialTimeout time.Duration
}

type DialOption interface {
	apply(*dialOptions)
}

type funcDialOption func(*dialOptions)

func (fdo funcDialOption) apply(do *dialOptions) {
	fdo(do)
}

func WithDB(dbname string) DialOption {
	return funcDialOption(func(do *dialOptions) {
		do.DB = dbname
	})
}

func WithAddr(addr string) DialOption {
	return funcDialOption(func(do *dialOptions) {
		do.Addr = addr
	})
}

func WithDailTimeout(timeout time.Duration) DialOption {
	return funcDialOption(func(do *dialOptions) {
		do.DialTimeout = timeout
	})
}

func WithAccount(account string) DialOption {
	return funcDialOption(func(do *dialOptions) {
		do.Account = account
	})
}

func WithPassword(password string) DialOption {
	return funcDialOption(func(do *dialOptions) {
		do.Password = password
	})
}
