package main

import (
	"fmt"

	"github.com/chenxiio/chenxi/comm/rpc"
)

func main() {
	if err := rpc.Generate("../", "api", "api", "../proxy_gen.go"); err != nil {
		fmt.Println("error: ", err)
	}
}
