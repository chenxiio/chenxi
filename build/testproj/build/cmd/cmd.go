package main

import (
	"fmt"
	"os"

	"github.com/chenxiio/chenxi/build/testproj"
	"github.com/chenxiio/chenxi/cmd"
)

// go build -o ../../cmd.exe
// 初始化系统
var name string = "cmd"

//var bsdir string = "../../"

var bsdir string = "./"

// func init() {

// 	//logger.Init(bsdir, slog.LevelDebug)
// 	fmt.Println("main init")
// 	var err error
// 	chenxi.CX, err = chenxi.New(name, bsdir)

// 	if err != nil {
// 		slog.Error(err.Error())
// 		return
// 	}

// }
func main() {
	fmt.Println(os.Args)
	// _, err := chenxi.CX.LoadModule(os.Args[1], testproj.TypesMap, testproj.TypesMapApi)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	err := cmd.Excmd(testproj.TypesMapApi, bsdir, os.Args[1:]...)
	if err != nil {
		fmt.Println(err)
	}
}
