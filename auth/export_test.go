package auth

import "sync"

var (
	Export_newToken        = newToken
	Export_contextkeyToken = contextkeyToken
)

func Export_tokenSetM(ts Tokenset) sync.Map {
	////
	return ts.(*tokenSet).m
}
