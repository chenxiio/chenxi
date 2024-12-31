package main

// import (
// 	"fmt"

// 	socketio "github.com/googollee/go-socket.io"
// 	"github.com/googollee/go-socket.io/engineio"
// )

// func main() {
// 	c, err := socketio.NewServer()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	c.On("message", func(msg string) {
// 		fmt.Println("received:", msg)
// 	})

// 	c.On(socketio.OnConnection, func(s socketio.Conn) error {
// 		s.Emit("message", "hello")
// 		return nil
// 	})

// 	c.Wait()
// }
