package chenxi

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/api2"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm/rpc"
	"github.com/chenxiio/chenxi/comm/socketio"
	"github.com/chenxiio/chenxi/logger"
	"github.com/fanhai/mux"
	"github.com/filecoin-project/go-jsonrpc/auth"
	_ "github.com/mattn/go-sqlite3"
)

var slog *logger.Logger
var CX *Chenxi
var Simulation bool = false

type Chenxi struct {
	Name     string
	IOServer api2.IOServerAPI
	Recipe   api.RecipeApi
	Alm      api2.ALMApi
	Job      api.JobApi
	Cfg      *cfg.Cfg
	socketio *socketio.SocketioServer
	Modules  *Modules
}

func (c *Chenxi) GetSocketio() *socketio.SocketioServer {
	return c.socketio
}
func (c *Chenxi) Close() error {
	//	err := c.IOServer.Close()
	logger.Close()
	// c.DB.Close()
	return nil
}
func New(name, path string) (*Chenxi, error) {
	slog = logger.GetLog(name, "", path)
	_cfg := cfg.NewCfg(path)
	err := _cfg.LoadAll()
	if err != nil {
		return nil, err
	}
	mux1 = mux.NewRouter()
	if p, ok := _cfg.Modules.Processes[name]; ok {
		if p.ProcessName == name {
			go stratFullapirpc(_cfg.Modules.Processes[name].Url)
		}
	}

	// dir := path + "DB"
	// err = os.MkdirAll(dir, os.ModePerm)
	// if err != nil {
	// 	slog.Error("Failed to connect to database: ", "err", err.Error())
	// 	return nil, err
	// }
	// db, err := sql.Open("sqlite3", path+"DB/chenxi.db")
	// if err != nil {
	// 	slog.Error("Failed to connect to database: ", "err", err.Error())
	// 	return nil, err
	// }
	socketioServer := socketio.NewServer()
	mux1.Handle("/socket.io/", socketioServer)
	return &Chenxi{
		//IOServer: ioserver,
		Name:     name,
		Cfg:      _cfg,
		socketio: socketioServer,
		// DB:   db,
	}, nil
}

// func (c *Chenxi) LoadModule(name string, gen GenMap, genapi GenMapapi) (any, error) {
// 	//c.Modules = &Modules{Items: make([]interface{}, 0)}
// 	v, err := c.Cfg.Modules.GetCfgByUnit(name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	internal, out := genapi(v.API)

// 	_, err = rpc.GetRpcClient("", "ws://"+c.Cfg.Modules.Processes[v.ProcessName].Url, v.Name, internal, time.Second*10)

//		if err != nil {
//			slog.Error(err.Error())
//		}
//		slog.Debug(fmt.Sprintf("ws://%s/%s/v0", c.Cfg.Modules.Processes[v.ProcessName].Url, v.Name))
//		//	defer closer()
//		//	c.Modules.Items[v.ID] = out
//		return out, nil
//	}
func (c *Chenxi) Load(gen GenMap, genapi GenMapapi) error {
	m := NewModules(c.Cfg.Modules, gen, genapi)
	c.Modules = m
	if c.Name == "ioserver" {
		c.IOServer.(*IOServer).Loop()

		// if c.Job!=nil {

		// }
		//c.Job.Init(context.TODO(), "")
	}

	// for _, v := range c.Cfg.Modules.Items {
	// 	if v.Name == "ioserver" {
	// 		if v1, ok := m.Items[v.ID].(*IOServer); ok {
	// 			c.IOServer = v1

	// 		}
	// 		break
	// 	}
	// }

	// rcp, err := c.GetModule("recipe")
	// if err == nil {
	// 	c.Recipe = rcp.(api.RecipeApi)
	// }
	// if c.Name == "ioserver" {
	// 	c.IOServer.(*IOServer).Loop()
	// 	initsocket()
	// }

	return nil
}
func (c *Chenxi) GetModule(name string) (any, error) {

	for _, v := range c.Cfg.Modules.Items {
		if v.Name == name {
			return c.Modules.Items[v.ID], nil
		}
	}

	return nil, fmt.Errorf("not fond %s", name)
}

var mutex sync.Mutex
var mux1 *mux.Router

func stratFullapirpc(url string) error {
	slog.Info("stratFullapirpc:" + url)

	//err := c.AppendFullapirpc(url)
	// if err != nil {
	// 	return err
	// }
	ah := &auth.Handler{
		Verify: rpc.AuthVerify,
		Next:   mux1.ServeHTTP,
	}

	err := http.ListenAndServe(url, ah)
	if err != nil {
		panic(err)
	}
	return err
}

func appendrpc(name string, obj, Internal, out any) error {
	mutex.Lock()
	defer mutex.Unlock()

	cors := &rpc.CorsHandler{Origin: "*", HandlerFunc: rpc.RegisterRpc(obj, Internal, out, name).ServeHTTP}
	prx := fmt.Sprintf("/%s/v0", name)
	mux1.Handle(prx, cors).Name(prx)

	slog.Info("Server name : " + prx)

	return nil
}
