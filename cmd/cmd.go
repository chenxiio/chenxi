package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm/rpc"
	socketio "github.com/googollee/go-socket.io"
	"github.com/urfave/cli"
	"golang.org/x/exp/slog"
)

func ApendCMDServerCmd(app *cli.App) {
	// gflgs := []cli.Flag{
	// 	cli.StringFlag{
	// 		Name:     "password,p",
	// 		Value:    "123456",
	// 		Usage:    "token密码",
	// 		Required: false,
	// 	}, cli.StringFlag{
	// 		Name:     "url,u", // 配置名称
	// 		Value:    "localhost:10600",
	// 		Required: false,
	// 		Usage:    "默认地址", // 配置描述
	// 	},
	// 	cli.StringFlag{
	// 		Name:        "token,t", // 配置名称
	// 		Value:       "",
	// 		Required:    false,
	// 		Usage:       "jwt token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.ae22SYWjZ_RJRRhfWDpVFzWThu_6EQ-iBgAn8vdrR-w", // 配置描述
	// 		Destination: &rpc.Cfg.Token,
	// 	},
	// }
	// app.Flags = append(app.Flags, gflgs...)

	app.Commands = append(app.Commands, *Cmdserver)
}

var Cmdserver = &cli.Command{
	Name:  "cmd",
	Usage: "通过模块名，方法名 ，参数，调用接口(示例：cmd ioserver ReadInt IO1)",

	Flags: []cli.Flag{
		// cli.StringFlag{
		// 	Name:     "module,m", // 配置名称
		// 	Value:    "",
		// 	Required: true,  // 缺省配置值
		// 	Usage:    "模块名", // 配置描述
		// },
		// cli.StringFlag{
		// 	Name:     "method,md", // 配置名称
		// 	Value:    "",
		// 	Required: true,  // 缺省配置值
		// 	Usage:    "方法名", // 配置描述
		// },
		// cli.StringSliceFlag{
		// 	Name:     "parm,p", // 配置名称
		// 	Required: false,    // 缺省配置值
		// 	Usage:    "参数",     // 配置描述
		// },
	},
	Action: func(cctx *cli.Context) error {
		if len(cctx.Args()) < 2 {
			return fmt.Errorf("参数数量不能少于2个")
		}
		module := cctx.Args().First()

		apiobj, err := chenxi.CX.GetModule(module)
		if err != nil {
			return err
		}
		methodValue := reflect.ValueOf(apiobj).MethodByName(cctx.Args()[1])
		// 检查方法是否存在
		if !methodValue.IsValid() {

			return errors.New("方法不存在")
		}
		// 构造参数列表

		args := []reflect.Value{reflect.ValueOf(context.Background())}

		for _, v := range cctx.Args()[2:] {
			args = append(args, reflect.ValueOf(v))
		}
		// 调用方法并获取返回值
		result := methodValue.Call(args)
		// ret, err := w.ReadInt(context.Background(), cctx.String("name"))
		// if err != nil {
		// 	return err
		// }
		for _, v := range result {
			str, err := json.MarshalIndent(v.Interface(), "", "")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(str))
		}

		return nil
	},
}

func Excmd(genapi chenxi.GenMapapi, bsdir string, args ...string) error {
	cfgmd := cfg.NewCfgModules(bsdir)
	err := cfgmd.ReadConfigFile()
	if err != nil {
		return err
	}

	v, err := cfgmd.GetCfgByUnit(args[0])
	if err != nil {
		return err
	}

	internal, out := genapi(v.API)

	_, err = rpc.GetRpcClient("", "ws://"+cfgmd.Processes[v.ProcessName].Url, v.Name, internal, time.Second*10)

	if err != nil {
		slog.Error(err.Error())
	}
	slog.Debug(fmt.Sprintf("ws://%s/%s/v0", cfgmd.Processes[v.ProcessName].Url, v.Name))
	//	defer closer()
	//	c.Modules.Items[v.ID] = out

	if len(args) < 2 {
		return fmt.Errorf("参数数量不能少于2个")
	}

	methodValue := reflect.ValueOf(out).MethodByName(args[1])
	// 检查方法是否存在
	if !methodValue.IsValid() {
		return errors.New("方法不存在")
	}
	// 构造参数列表

	args1 := []reflect.Value{reflect.ValueOf(context.Background())}

	for _, v := range args[2:] {
		args1 = append(args1, reflect.ValueOf(v))
	}
	// 调用方法并获取返回值
	result := methodValue.Call(args1)
	// ret, err := w.ReadInt(context.Background(), cctx.String("name"))
	// if err != nil {
	// 	return err
	// }
	for _, v := range result {
		str, err := json.MarshalIndent(v.Interface(), "", "")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(str))
	}
	// ie := &IOEvent{topic: "k"}
	// out.(api2.IOServerAPI).su(context.TODO(), "k", ie)

	return nil
}

type IOEvent struct {
	// i *int
	//dataChannel chan map[string]interface{}
	topic string
	s     socketio.Conn
}

func (n *IOEvent) Dispatch(data ...map[string]any) {

	go n.s.Emit("device-values", data[0])
	for key, value := range data[0] {
		fmt.Printf("%s: %v\n", key, value)
	}
}
