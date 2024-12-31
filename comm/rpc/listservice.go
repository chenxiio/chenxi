package rpc

import "time"

type RpcService struct {
	Url           string
	Token         []byte
	V             interface{}
	Closer        func()
	Reaction_time time.Duration
}
