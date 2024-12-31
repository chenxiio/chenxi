package chenxi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chenxiio/chenxi/comm/socketio"
	"github.com/chenxiio/chenxi/logger"
)

type TEvent struct {
	// i *int
	// s string
	ios any
}

func (n *TEvent) Dispatch(data ...any) {

	for key, value := range data[0].(map[string]any) {

		fmt.Printf("%s: %v\n", key, value)
	}
}

func TestChenxi_ReadInt32(t *testing.T) {
	ctx := context.Background()
	socketioServer := socketio.NewServer()
	slog = logger.GetLog("io", "", "./")
	//mux1.Handle("/socket.io/", socketioServer)
	c, err := NewIOServer("./testdb/", nil, socketioServer)
	if err != nil {
		t.Fatal(err)
	}

	c.Sub("subio", "k", &TEvent{})
	// Write a value to the database
	value1, err := c.ReadInt(ctx, "key1")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(value1)
	err = c.WriteInt(ctx, "key1", 1284)
	if err != nil {
		t.Fatal(err)
	}
	value := int32(0)
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		// Read the value from the database
		value, err = c.ReadInt(ctx, "key1")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	elapsed := time.Since(start)
	fmt.Println("执行时间:", elapsed)
	// Check that the value was read correctly
	if value != 1284 {
		t.Errorf("unexpected value: %d", value)
	}

}

func TestChenxi_ReadString(t *testing.T) {
	ctx := context.Background()

	socketioServer := socketio.NewServer()
	mux1.Handle("/socket.io/", socketioServer)
	c, err := NewIOServer("./testdb", nil, socketioServer)
	if err != nil {
		t.Fatal(err)
	}

	c.bus.On("k", &TEvent{})
	// Write a value to the database
	err = c.WriteString(ctx, "keyeeee", "")
	if err != nil {
		t.Fatal(err)
	}
	err = c.SetState(ctx, "key1", "hello")
	if err != nil {
		fmt.Println(err)
	}
	err = c.SetState(ctx, "key1", "hello")
	if err != nil {
		fmt.Println(err)
	}
	err = c.SetState(ctx, "key1", "IDLE")
	if err != nil {
		fmt.Println(err)
	}
	err = c.SetState(ctx, "key1", "IDLE")
	if err != nil {
		fmt.Println(err)
	}
	err = c.SetState(ctx, "key1", "hello")
	if err != nil {
		fmt.Println(err)
	}

	// Read the value from the database
	value, err := c.ReadString(ctx, "key1")
	if err != nil {
		t.Fatal(err)
	}

	// Check that the value was read correctly
	if value != "hello" {
		t.Errorf("unexpected value: %s", value)
	}
}

func TestChenxi_ReadFloat64(t *testing.T) {
	ctx := context.Background()

	socketioServer := socketio.NewServer()
	mux1.Handle("/socket.io/", socketioServer)
	c, err := NewIOServer("./testdb", nil, socketioServer)
	if err != nil {
		t.Fatal(err)
	}
	c.bus.On("key1", &TEvent{})
	// Write a value to the database
	err = c.WriteDouble(ctx, "key1", 3.146)
	if err != nil {
		t.Fatal(err)
	}

	// Read the value from the database
	value, err := c.ReadDouble(ctx, "key1")
	if err != nil {
		t.Fatal(err)
	}

	// Check that the value was read correctly
	if value != 3.14 {
		t.Errorf("unexpected value: %f", value)
	}
}

// func TestChenxi_ReadBool(t *testing.T) {
// 	ctx := context.Background()

// 	c, err := NewIOServer("testdb", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	c.bus.On("k", &TEvent{})
// 	// Write a value to the database
// 	err = c.WriteBool(ctx, "key1", true)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Read the value from the database
// 	value, err := c.ReadBool(ctx, "key1")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// Check that the value was read correctly
// 	if value != true {
// 		t.Errorf("unexpected value: %t", value)
// 	}
// }
