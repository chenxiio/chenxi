package main

import (
	"fmt"
	"os"

	"github.com/chenxiio/chenxi/cmd"
	"github.com/chenxiio/chenxi/comm/rpc"
	"github.com/urfave/cli"
)

func main() {
	fmt.Println(os.Args)
	// ctx, cancel := context.WithCancel(context.Background())

	// defer cancel()

	//实例化cli
	app := cli.NewApp()
	//Name可以设定应用的名字
	app.Name = "chenxi cli"
	// Version可以设定应用的版本号
	app.Version = "1.0.0"
	cmd.ApendIOServerCmd(app)
	// Commands用于创建命令
	app.Commands = append(app.Commands, *rpc.AuthNewCmd)
	// 接受os.Args启动程序
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
	// for {
	// 	time.Sleep(time.Second * 1)
	// }
}
