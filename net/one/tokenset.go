package one

import (
	"context"
	"sync"

	"github.com/hiank/think/net/box"
	"github.com/hiank/think/run"
)

var (
	tokenset     *tokenSet
	tokensetOnce sync.Once
)

const (
	ErrTokenExisted = run.Err("one: cannot build root Token when it already existed")
	ErrNonRootoken  = run.Err("one: non root Token")
)

type tokenSet struct {
	ctx context.Context //base to TODO
	m   sync.Map
}

//Build build root Token.
//will return error when Token existed
func (ts *tokenSet) Build(uid string) (box.Token, error) {
	if _, ok := ts.m.Load(uid); !ok {
		rt := box.NewToken(context.WithValue(ts.ctx, box.ContextkeyTokenUid, uid))
		if _, loaded := ts.m.LoadOrStore(uid, rt); !loaded {
			return rt, nil
		}
	}
	return nil, ErrTokenExisted
}

func (ts *tokenSet) Derive(uid string) (tk box.Token) {
	v, ok := ts.m.Load(uid)
	if !ok {
		rt := box.NewToken(context.WithValue(ts.ctx, box.ContextkeyTokenUid, uid))
		v, _ = ts.m.LoadOrStore(uid, rt)
	}
	return v.(box.Token).Fork()
	// return
}

// //Get root Token for uid
// func (ts *tokenSet) Derive(uid string) (tk box.Token, err error) {
// 	if v, ok := ts.m.Load(uid); ok {
// 		tk = v.(box.Token).Fork()
// 	} else {
// 		err = ErrNonRootoken
// 	}
// 	return
// }

func (ts *tokenSet) Kill(uid string) (err error) {
	if v, loaded := ts.m.LoadAndDelete(uid); loaded {
		v.(box.Token).Close()
	} else {
		err = ErrNonRootoken
	}
	return
}

func TokenSet() Tokenset {
	tokensetOnce.Do(func() {
		tokenset = &tokenSet{
			ctx: TODO(),
		}
	})
	return tokenset
}
