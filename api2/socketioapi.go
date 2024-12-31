package api2

import (
	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/comm/socketio"
)

type Socketioapi interface {
	Sub(group string, key string, f socketio.SocketsubDispatch) error
	Unsub(group string, key string, f socketio.SocketsubDispatch) error
}

type IOServerAPI interface {
	api.IOServerAPI
	Socketioapi
}

type IOServerClient struct {
	api.IOServerAPI
	Socketioapi
}
type ALMWaitackApi interface {
	WaitAlmAck(aid int64) (int, error)
}
type ALMApi interface {
	ALMWaitackApi
	api.ALMApi
	Socketioapi
}
type ALMClient struct {
	ALMWaitackApi
	api.ALMApi
	Socketioapi
}
