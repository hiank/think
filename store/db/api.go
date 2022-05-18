package db

import (
	"context"

	"github.com/hiank/think/store"
)

//Dialer dial to database
type Dialer interface {
	Dial(ctx context.Context, opts ...DialOption) (store.EasyDictionary, error)
}

//FuncDialer convert func to Dialer
type FuncDialer[T ~string] func(ctx context.Context, opts ...DialOption) (store.Dictionary[T], error)

func (fd FuncDialer[T]) Dial(ctx context.Context, opts ...DialOption) (ed store.EasyDictionary, err error) {
	d, err := fd(ctx, opts...)
	if err == nil {
		ed = store.ConvertoEasy(d)
	}
	return
}
