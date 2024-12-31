package main

import (
	"fmt"
	"os"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/build/testproj"
	"github.com/chenxiio/chenxi/logger"
)

// go build -o ../../ioserver.exe
// 初始化系统
var name string = "ioserver"

var bsdir string = "../../"

//var bsdir string = "./"

func main() {
	// for i := 0; i < 100; i++ {
	// 	fmt.Println(os.Args)
	// }
	fmt.Println(os.Args)
	if len(os.Args) > 1 {
		if os.Args[1] == "Simulation" {
			chenxi.Simulation = true
		}
	}

	fmt.Println("main init")
	var err error
	chenxi.CX, err = chenxi.New(name, bsdir)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	chenxi.CX.Load(testproj.TypesMap, testproj.TypesMapApi)
	select {}
}
