package mod

var (
	TaskHub = &taskHub{}
)

type taskHub struct {
	// ctx          context.Context
	// cancel       context.CancelFunc
	// hub          map[string]*token.Master //NOTE: map[uuid]*Token
	// subscription *nats.Subscription
	// mux          sync.RWMutex
}

// func (th *tokenHub) Depend() []think.Module {
// 	return []think.Module{NatsCLI}
// }

// func (th *tokenHub) OnCreate(ctx context.Context) error {
// 	th.ctx, th.cancel = context.WithCancel(context.Background())
// 	th.hub = make(map[string]*token.Master)
// 	return nil
// }

// func (th *tokenHub) OnStart(ctx context.Context) (err error) {
// 	th.subscription, err = NatsCLI.Subscribe(modex.NewsTokenInvalid, func(msg *nats.Msg) {})
// 	return
// }

// func (th *tokenHub) OnStop() {
// 	th.subscription.Unsubscribe()
// 	th.subscription = nil
// }

// func (th *tokenHub) OnDestroy() {
// 	th.cancel()
// }

// func (th *tokenHub) AutoToken(uuid string) *token.Master {
// 	if tk, existed := th.cachedToken(uuid); existed {
// 		return tk
// 	}
// 	th.mux.Lock()
// 	defer th.mux.Unlock()
// 	tk, existed := th.hub[uuid]
// 	if !existed {
// 		tk = token.NewMaster(th.ctx, uuid)
// 		th.hub[uuid] = tk
// 	}
// 	return tk
// }

// func (th *tokenHub) cachedToken(uuid string) (tk *token.Master, existed bool) {
// 	th.mux.RLock()
// 	defer th.mux.RUnlock()
// 	tk, existed = th.hub[uuid]
// 	return
// }
