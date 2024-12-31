package main

import (
	"fmt"

	"github.com/chenxiio/chenxi/comm/socketio"
)

type IOEvent struct {
	// i *int
	//dataChannel chan map[string]interface{}
	topic string
}

func (n *IOEvent) Dispatch(data ...any) {

	for key, value := range data[0].(map[string]any) {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func main() {
	//connect to server, you can use your own transport settings
	// parms := make(map[string]string)
	// tr := transport.GetDefaultWebsocketTransport()
	// ws := gosio.New(gosio.GetURL("localhost", 10600, false, &parms), tr)

	ws := socketio.NewClient("localhost", 10600)

	ws.Sub("subio", "IO", &IOEvent{topic: "IO"})
	//	ws.On("subio")
	select {}
}
