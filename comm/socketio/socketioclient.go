package socketio

import (
	"fmt"
	"log"
	"sync"

	gosio "github.com/gnabgib/go-sio"
	"github.com/gnabgib/go-sio/transport"
)

type Socketioclient struct {
	*gosio.Client
	subs map[string][]string
	lock sync.Mutex
}

func NewClient(ip string, port int) *Socketioclient {
	parms := make(map[string]string)
	tr := transport.GetDefaultWebsocketTransport()
	sk := &Socketioclient{gosio.New(gosio.GetURL(ip, port, false, &parms), tr), map[string][]string{}, sync.Mutex{}}

	sk.OnDisconnect(func(c *gosio.Channel) {
		log.Println("Disconnected to server1")

		err := sk.Dial()
		if err != nil {

			fmt.Printf("dial err %s", err.Error())
		}

	})
	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(1) // 增加等待组计数器
	sk.OnConnect(func(c *gosio.Channel) {
		sk.lock.Lock()
		defer sk.lock.Unlock()

		log.Println("Connected to server1")
		for k, v := range sk.subs {
			for _, v1 := range v {
				log.Println(k, v1)

				err := c.Emit(k, v1)
				if err != nil {
					log.Println(err.Error())
				}
			}

		}
		if wg != nil {
			wg.Done()
		}

	})

	err := sk.Dial()
	if err != nil {
		wg.Done()
		log.Println(err.Error())
	}
	wg.Wait() // 等待所有等待组计数器归零
	wg = nil
	return sk
}

func (s *Socketioclient) Sub(group string, key string, e SocketsubDispatch) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	fmt.Println("Socketioclient sub ")
	if _, ok := s.subs[group]; !ok {
		s.subs[group] = make([]string, 0)
	}
	s.subs[group] = append(s.subs[group], key)

	s.On(group+key, func(c *gosio.Channel, msg any) {
		e.Dispatch(msg)
	})
	err := s.Emit(group, key)
	if err != nil {
		return err
	}
	return nil
}

func (s *Socketioclient) Unsub(group string, key string, e SocketsubDispatch) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.subs[group]; ok {
		for k, v := range s.subs[group] {
			if v == key {
				s.subs[key] = append(s.subs[group][:k], s.subs[group][k+1:]...)
				break
			}
		}
	}
	err := s.Emit("un"+group, key)
	if err != nil {
		return err
	}
	// if subs, ok := s.subs[key]; ok {
	// 	for i := 0; i < len(subs); i++ {
	// 		if subs[i] == p {
	// 			copy(subs[i:], subs[i+1:])
	// 			s.subs[key] = subs[:len(subs)-1]
	// 			break
	// 		}
	// 	}
	// }
	return nil
}
