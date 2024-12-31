package main

import (
	"fmt"

	urpc "github.com/chenxiio/chenxi/comm/rpc"
)

func main() {

	ge2()

}

func ge1() {
	// latest (v1)
	if err := urpc.Generate2("../api", "api", "api", "../api/proxy_gen2.go"); err != nil {
		fmt.Println("error: ", err)
	}

	if err := urpc.Generate("../api", "api", "api", "../api/proxy_gen.go"); err != nil {
		fmt.Println("error: ", err)
	}

}

func ge2() {
	// latest (v1)
	// if err := urpc.Generate2("../chain", "chain", "chain", "../chain/proxy_gen2.go"); err != nil {
	// 	fmt.Println("error: ", err)
	// }

	if err := urpc.Generate("../", "chenxi", "chenxi", "../proxy_gen.go"); err != nil {
		fmt.Println("error: ", err)
	}

	// v0

}
