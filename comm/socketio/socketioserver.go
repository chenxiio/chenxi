package socketio

import (
	"sync"

	socketios "github.com/googollee/go-socket.io"
	"golang.org/x/exp/slog"
)

type SocketsubDispatch interface {
	Dispatch(data ...any)
}
type SocketsubEvent interface {
	Sub(k string)
	UnSub(k string)
}
type SocketEvent interface {
	SocketsubEvent
	SocketsubDispatch
}

type SocketioServer struct {
	*socketios.Server
	Events map[string]map[string]SocketsubEvent
	lock   sync.Mutex
}

func NewServer() *SocketioServer {
	server := SocketioServer{Server: socketios.NewServer(nil), Events: map[string]map[string]SocketsubEvent{}}

	server.OnConnect("/", func(s socketios.Conn) error {
		//s.SetContext("")
		server.lock.Lock()
		defer server.lock.Unlock()
		if _, ok := server.Events[s.ID()]; !ok {
			server.Events[s.ID()] = make(map[string]SocketsubEvent)
		}
		slog.Info("connected:", s.ID())
		return nil
	})

	server.OnError("/", func(s socketios.Conn, e error) {
		slog.Error("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketios.Conn, reason string) {
		slog.Error("closed", reason)
		server.lock.Lock()
		defer server.lock.Unlock()
		for k, v := range server.Events[s.ID()] {
			v.UnSub(k)
		}
		server.Events[s.ID()] = nil

	})

	go func() {
		if err := server.Serve(); err != nil {
			slog.Error("socketio listen error: " + err.Error())
		}
		//server.ForEach("","",)
		// for i := 0; i < 100; i++ {
		// 	time.Sleep(time.Second * 10)

		// }
	}()
	return &server
}

func (s *SocketioServer) Sub(id, key string, e SocketsubEvent) {
	s.lock.Lock()
	defer s.lock.Unlock()
	e.Sub(key)
	if _, ok := s.Events[id]; !ok {
		s.Events[id] = map[string]SocketsubEvent{}
	}
	s.Events[id][key] = e
}

func (s *SocketioServer) UnSub(id, key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.Events[id]; !ok {
		s.Events[id] = map[string]SocketsubEvent{}
		return
	}
	s.Events[id][key].UnSub(key)
	delete(s.Events[id], key)
}
