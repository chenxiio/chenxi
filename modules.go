package chenxi

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/api2"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm/rpc"
	"github.com/chenxiio/chenxi/comm/socketio"
)

// type Module struct {
// }
type GenMap func(string, string, string, string) any
type GenMapapi func(string) (internal, out any)
type Modules struct {
	//cfg   *cfg.Modules
	Items []interface{}
}

// loadtype: dll ,plugin ,class,exe
// Type : driver,pm,cm,ctc ,custom
func NewModules(c cfg.Modules, gen GenMap, genapi GenMapapi) *Modules {
	ret := Modules{}

	sort.Slice(c.Items, func(i, j int) bool {
		return c.Items[i].ID < c.Items[j].ID
	})
	ret.Items = make([]interface{}, c.Items[len(c.Items)-1].ID+1)
	ioclients := map[string]api2.Socketioapi{}
	for _, v := range c.Processes {
		if v.ProcessName != CX.Name {
			args := strings.Split(c.Processes[v.ProcessName].Url, ":")
			port, err := strconv.Atoi(args[1])
			if err != nil {
				slog.Error(err.Error())
			}
			ioclient := socketio.NewClient(args[0], port)
			ioclients[v.ProcessName] = ioclient
			//c := api2.IOServerClient{ioclient}
		}
	}
	for _, v := range c.Items {

		if v.ProcessName == CX.Name {

			obj := gen(v.Name, v.Type, v.Path, v.Parm)
			ret.Items[v.ID] = obj
			switch v.API {
			case "DRIVER":

			default:
				internal, out := genapi(v.API)
				err := appendrpc(v.Name, obj, internal, out)
				if err != nil {
					slog.Error(err.Error())
				}
			}

			switch v.Name {
			case "ioserver":
				CX.IOServer = ret.Items[v.ID].(api2.IOServerAPI)
				// CX.IOServer = api2.IOServerClient{
				// 	IOServerAPI: ret.Items[v.ID].(api.IOServerAPI),
				// 	Socketioapi: ioclients[v.ProcessName],
				// }
			case "recipe":
				CX.Recipe = ret.Items[v.ID].(api.RecipeApi)
			case "job":

				CX.Job = ret.Items[v.ID].(api.JobApi)
			case "alarm":
				CX.Alm = ret.Items[v.ID].(api2.ALMApi)

			}
		} else {
			switch v.API {
			case "DRIVER":
			case "ioserver":

			default:
				internal, out := genapi(v.API)

				_, err := rpc.GetRpcClient("", "ws://"+c.Processes[v.ProcessName].Url, v.Name, internal, time.Second*10)

				if err != nil {
					slog.Error(err.Error())
				}
				slog.Debug(fmt.Sprintf("ws://%s/%s/v0", c.Processes[v.ProcessName].Url, v.Name))
				//	defer closer()
				ret.Items[v.ID] = out
			}

			switch v.Name {
			case "ioserver":
				//	CX.IOServer = api2.IOServerClient{ret.Items[v.ID].(api.IOServerAPI), ioclients[v.ProcessName]}
				CX.IOServer = api2.IOServerClient{
					IOServerAPI: ret.Items[v.ID].(api.IOServerAPI),
					Socketioapi: ioclients[v.ProcessName],
				}
			case "recipe":
				CX.Recipe = ret.Items[v.ID].(api.RecipeApi)
			case "job":

				CX.Job = ret.Items[v.ID].(api.JobApi)
			case "alarm":
				CX.Alm = api2.ALMClient{
					ALMApi:        ret.Items[v.ID].(api.ALMApi),
					Socketioapi:   ioclients[v.ProcessName],
					ALMWaitackApi: &api2.AlmWaitAck{Socketioapi: ioclients[v.ProcessName]},
				}

			}
		}

	}

	return &ret
}

// func GenericFunc[T any](input T) T {
// 	return input
// }
