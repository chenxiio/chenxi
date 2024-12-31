package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/comm/rpc"
	"github.com/urfave/cli"
)

func ApendIOServerCmd(app *cli.App) {
	gflgs := []cli.Flag{
		cli.StringFlag{
			Name:     "password,p",
			Value:    "123456",
			Usage:    "token密码",
			Required: false,
		}, cli.StringFlag{
			Name:     "url,u", // 配置名称
			Value:    "localhost:10600",
			Required: false,
			Usage:    "默认地址", // 配置描述
		},
		cli.StringFlag{
			Name:        "token,t", // 配置名称
			Value:       "",
			Required:    false,
			Usage:       "jwt token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.ae22SYWjZ_RJRRhfWDpVFzWThu_6EQ-iBgAn8vdrR-w", // 配置描述
			Destination: &rpc.Cfg.Token,
		},
	}
	app.Flags = append(app.Flags, gflgs...)
	wcmd := []cli.Command{*ReadIntCmd, *ReadDoubleCmd, *ReadStringCmd, *WriteDoubleCmd, *WriteIntCmd, *WriteStringCmd}
	app.Commands = append(app.Commands, wcmd...)
}

var ReadIntCmd = &cli.Command{
	Name:  "readint",
	Usage: "读取整数值(示例:cli.exe readint -name IO1)",

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name", // 配置名称
			Value:    "",
			Required: true,      // 缺省配置值
			Usage:    "io name", // 配置描述
		},
	},

	Action: func(cctx *cli.Context) error {
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		closer, err := rpc.GetRpcClient(cctx.GlobalString("token"), "ws://"+cctx.GlobalString("url"), "ioserver", &w.Internal, time.Second*10)

		if err != nil {
			return err
		}
		defer closer()

		ret, err := w.ReadInt(context.Background(), cctx.String("name"))
		if err != nil {
			return err
		}
		fmt.Println(ret)

		return nil
	},
}
var ReadStringCmd = &cli.Command{
	Name:  "readstring",
	Usage: "读取字符串值(示例:cli.exe readstring -name IO3)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name",
			Value:    "",
			Required: true,
			Usage:    "io name",
		},
	},
	Action: func(cctx *cli.Context) error {
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		closer, err := rpc.GetRpcClient(cctx.GlobalString("token"), "ws://"+cctx.GlobalString("url"), "ioserver", &w.Internal, time.Second*10)
		defer closer()
		if err != nil {
			return err
		}
		ret, err := w.ReadString(context.Background(), cctx.String("name"))
		if err != nil {
			return err
		}
		fmt.Println(ret)
		return nil
	},
}
var ReadDoubleCmd = &cli.Command{
	Name:  "readdouble",
	Usage: "读取浮点数值(示例:cli.exe readdouble -name IO2)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name",
			Value:    "",
			Required: true,
			Usage:    "io name",
		},
	},
	Action: func(cctx *cli.Context) error {
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		closer, err := rpc.GetRpcClient(cctx.GlobalString("token"), "ws://"+cctx.GlobalString("url"), "ioserver", &w.Internal, time.Second*10)
		defer closer()
		if err != nil {
			return err
		}
		ret, err := w.ReadDouble(context.Background(), cctx.String("name"))
		if err != nil {
			return err
		}
		fmt.Println(ret)
		return nil
	},
}
var WriteIntCmd = &cli.Command{
	Name:  "writeint",
	Usage: "写入整数值(示例:cli.exe writeint -name IO1 -value 6)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name",
			Value:    "",
			Required: true,
			Usage:    "io name",
		},
		cli.IntFlag{
			Name:     "value",
			Value:    0,
			Required: true,
			Usage:    "integer value",
		},
	},
	Action: func(cctx *cli.Context) error {
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		closer, err := rpc.GetRpcClient(cctx.GlobalString("token"), "ws://"+cctx.GlobalString("url"), "ioserver", &w.Internal, time.Second*10)
		defer closer()
		if err != nil {
			return err
		}
		err = w.WriteInt(context.Background(), cctx.String("name"), int32(cctx.Int("value")))
		if err != nil {
			return err
		}
		fmt.Println("WriteInt successful")
		return nil
	},
}
var WriteStringCmd = &cli.Command{
	Name:  "writestring",
	Usage: "写入字符串值(示例:cli.exe writestring -name IO3 -value 我是字符串)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name",
			Value:    "",
			Required: true,
			Usage:    "io name",
		},
		cli.StringFlag{
			Name:     "value",
			Value:    "",
			Required: true,
			Usage:    "string value",
		},
	},
	Action: func(cctx *cli.Context) error {
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		closer, err := rpc.GetRpcClient(cctx.GlobalString("token"), "ws://"+cctx.GlobalString("url"), "ioserver", &w.Internal, time.Second*10)
		defer closer()
		if err != nil {
			return err
		}
		err = w.WriteString(context.Background(), cctx.String("name"), cctx.String("value"))
		if err != nil {
			return err
		}
		fmt.Println("WriteString successful")
		return nil
	},
}
var WriteDoubleCmd = &cli.Command{
	Name:  "writedouble",
	Usage: "写入浮点数值(示例：cli.exe writedouble -name IO2 -value 0.9)",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "name",
			Value:    "",
			Required: true,
			Usage:    "io name",
		},
		cli.Float64Flag{
			Name:     "value",
			Value:    0.0,
			Required: true,
			Usage:    "float value",
		},
	},
	Action: func(cctx *cli.Context) error {
		var w api.IOServerAPIStruct = api.IOServerAPIStruct{}
		closer, err := rpc.GetRpcClient(cctx.GlobalString("token"), "ws://"+cctx.GlobalString("url"), "ioserver", &w.Internal, time.Second*10)
		defer closer()
		if err != nil {
			return err
		}
		err = w.WriteDouble(context.Background(), cctx.String("name"), cctx.Float64("value"))
		if err != nil {
			return err
		}
		fmt.Println("WriteDouble successful")
		return nil
	},
}
