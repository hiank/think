package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hiank/think/oauth"
)

var (
	Export_getHelperServer = func(sh *ServeHelper) *http.Server {
		return sh.server
	}
	Export_getHelperUpgrader = func(sh *ServeHelper) *websocket.Upgrader {
		return sh.upgrader
	}
	Export_getHelperAuther = func(sh *ServeHelper) oauth.Auther {
		return sh.auther
	}
	Export_getHelperDopts = func(sh *ServeHelper) *options {
		return sh.dopts
	}
	Export_getOptionsConnMaker = func(sh *options) ConnMaker {
		return sh.connMaker
	}
)

// server   *http.Server
// upgrader *websocket.Upgrader //NOTE: use default options
// auther   oauth.Auther
// dopts    *options
// net.ChanAccepter
