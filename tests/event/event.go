package main

import (
	"fmt"

	"github.com/gookit/event"
)

func main() {
	// Register event listener
	event.On("evt1.*", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handle event: %s\n", e.Name())
		return nil
	}), event.Normal)

	// Register multiple listeners
	event.On("evt1.*", event.ListenerFunc(func(e event.Event) error {
		fmt.Printf("handle event: %s\n", e.Data())
		return nil
	}), event.High)

	// ... ...

	// Trigger event
	// Note: The second listener has a higher priority, so it will be executed first.
	event.MustFire("evt1.xxx", event.M{"arg0": "val0", "arg1": "val1"})
}
