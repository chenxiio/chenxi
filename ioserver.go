package chenxi

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm/socketio"
	"github.com/chenxiio/chenxi/models"
	eventbus "github.com/chenxiio/go-eventbus"

	socketios "github.com/googollee/go-socket.io"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type IOEvent struct {
	// i *int
	//dataChannel chan map[string]interface{}
	topic    string
	s        socketios.Conn
	ioserver *IOServer
}

func (n *IOEvent) Sub(k string) {
	n.ioserver.bus.On(k, n)
}
func (n *IOEvent) UnSub(k string) {
	n.ioserver.bus.Off(k, n)
}

// func (n *IOEvent) OnDisconnect(k string) {

//		n.ioserver.bus.Off(n.topic, n)
//	}
func (n *IOEvent) Dispatch(data ...any) {

	go n.s.Emit("subio"+n.topic, data[0])
	// for key, value := range data[0] {
	// 	fmt.Printf("%s: %v\n", key, value)
	// }
}

type IOServer struct {
	cfg       *cfg.IOCfg
	db        *leveldb.DB
	bus       *eventbus.Bus[any]
	lock      sync.Mutex
	socketio  *socketio.SocketioServer
	his       *HistoriesDao
	writeChan chan *models.His
}

func (c *IOServer) Close() error {
	c.bus.Clean()
	return c.db.Close()
}

const sqlbulkcount = 2000

func NewIOServer(path string, _cfg *cfg.IOCfg, socketio *socketio.SocketioServer) (*IOServer, error) {

	db, err := leveldb.OpenFile(path+"cache", nil)
	if err != nil {
		return nil, err
	}
	sqldb, err := sql.Open("sqlite3", path+"data/his.db")
	if err != nil {
		slog.Error("Failed to connect to database: ", "err", err.Error())
		return nil, err

		//panic(err)
	}
	hisdb := &HistoriesDao{db: sqldb}

	err = hisdb.createTable()
	if err != nil {
		slog.Error("Failed to connect to database: ", "err", err.Error())
		return nil, err

		//panic(err)
	}
	bus := eventbus.New[any]()
	if _cfg == nil {
		_cfg = &cfg.IOCfg{Items: make(cfg.IODefines)}
	}

	// 创建一个缓冲大小为10的chan通道
	writeChan := make(chan *models.His, sqlbulkcount)

	ios := &IOServer{
		cfg:       _cfg,
		db:        db,
		bus:       bus,
		socketio:  socketio,
		his:       hisdb,
		writeChan: writeChan,
	}

	// 启动一个goroutine来处理写入操作
	go func() {
		for {
			count := len(writeChan)
			if count == 0 {
				count = 1
			}
			if count > sqlbulkcount {
				count = sqlbulkcount
			}
			//fmt.Println(count)
			histories := make([]*models.His, count)
			for i := 0; i < count; i++ {
				//len(writeChan)
				histories[i] = <-writeChan
			}
			// err := hisdb.BulkInsert(histories)
			// if err != nil {
			// 	// 错误处理
			// 	slog.Error(err.Error())
			// }

		}
	}()

	ios.socketio.OnEvent("/", "subio", func(s socketios.Conn, prexx string) {
		slog.Debug("subio", s.ID(), prexx)
		ie := &IOEvent{topic: prexx, s: s, ioserver: ios}

		ios.socketio.Sub(s.ID(), prexx, ie)

		// 返回该前缀所有数据
		//s.Emit("reply", "have "+prexx)
		mp, err := ios.ReadFromPrefix(context.TODO(), prexx)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		s.Emit("subio"+prexx, mp)
	})

	ios.socketio.OnEvent("/", "unsubio", func(s socketios.Conn, prexx string) {
		slog.Debug("unsubio", s.ID(), prexx)

		ios.socketio.UnSub(s.ID(), prexx)
	})
	// ios.socketio.JoinRoom("","")
	// ios.socketio.BroadcastToRoom()

	//	err = ios.Appendrpc()
	return ios, err
}

func (c *IOServer) insertWorker(writeChan <-chan *models.His, batchSize int) {
	// 创建一个切片用于存储待插入的历史记录
	histories := make([]*models.His, 0, batchSize)
	for history := range writeChan {
		// 将历史记录添加到切片中
		histories = append(histories, history)
		// 当切片中的历史记录数量达到批量大小时，执行批量插入
		if len(histories) >= batchSize {
			err := c.his.BulkInsert(histories)
			if err != nil {
				// 错误处理
				fmt.Println(err)
			}
			// 清空切片，准备下一批数据
			histories = make([]*models.His, 0, batchSize)
		}
	}
	// 处理剩余的历史记录
	if len(histories) > 0 {
		err := c.his.BulkInsert(histories)
		if err != nil {
			// 错误处理
			fmt.Println(err)
		}
	}
}

// func NewIOClient(url string) (api.IOServerAPI, error) {

// }

// func (c *IOServer) Appendrpc() error {
// 	mutex.Lock()
// 	defer mutex.Unlock()

// 	var out IOServerAPIStruct

// 	cors := &rpc.CorsHandler{Origin: "*", HandlerFunc: rpc.RegisterRpc(c, &out.Internal, &out, "ioserver").ServeHTTP}

// 	mux1.Handle("/ioserver/v0", cors).Name("/ioserver/v0")

// 	slog.Info("Server name : /ioserver/v0")

// 	return nil
// }

// func (w *IOServer) StratrpcServer(url string) error {

// 	mux1 = mux.NewRouter()

// 	//	rpcServer := jsonrpc.NewServer()

// 	var out IOServerAPIStruct
// 	// auth.PermissionedProxy(rpc.AllPermissions, rpc.DefaultPerms, w, &out.Internal)
// 	// rpcServer.Register("wallet", &out)

// 	mux1.Handle("/chenxi/v0", rpc.RegisterRpc(w, &out.Internal, &out, "chenxi"))

// 	ah := &auth.Handler{
// 		Verify: rpc.AuthVerify,
// 		Next:   mux1.ServeHTTP,
// 	}

//		return http.ListenAndServe(url, ah)
//	}
func (c *IOServer) Unsub(group string, key string, e socketio.SocketsubDispatch) error {
	slog.Debug("IOServer Off ", "", key)
	c.bus.Off(key, e)
	return nil
}

// key =subio
func (c *IOServer) Sub(key string, parm string, f socketio.SocketsubDispatch) error {
	slog.Debug("IOServer On ", key, parm)
	c.bus.On(parm, f)

	mp, err := c.ReadFromPrefix(context.TODO(), parm)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	f.(eventbus.Event[any]).Dispatch(mp)
	return nil
}
func (c *IOServer) Loop() {
	var groupitems map[int]cfg.IODefines = make(map[int]cfg.IODefines)
	for k, v := range c.cfg.Items {
		if _, ok := groupitems[v.Rs]; !ok {
			groupitems[v.Rs] = make(cfg.IODefines)
		}
		groupitems[v.Rs][k] = v
	}

	for k, iodf := range groupitems {
		go func(rs int, ioitems cfg.IODefines) {
			for {
				time.Sleep(time.Millisecond * time.Duration(rs))
				//fmt.Println(time.Now())
				for _, v := range ioitems {
					if v.Cat == "IO" {
						if drv, ok := CX.Modules.Items[v.Drv].(api.Drvapi); ok {

							switch v.DT {
							case "int":
								ret, err := drv.ReadInt(context.TODO(), v.Pr)
								if err != nil {
									slog.Error(err.Error(), v.Name, v.Pr)
									continue
								}
								if ok, _ := c.ischange(context.TODO(), v.Name, ret); !ok {
									continue
								}
								err = c.put(context.TODO(), v.Name, ret, models.IO_TYPE_INT)
								if err != nil {
									slog.Error(err.Error(), v.Name, v.Pr)
									continue
								}
							case "string":
								ret, err := drv.ReadString(context.TODO(), v.Pr)
								if err != nil {
									slog.Error(err.Error(), v.Name, v.Pr)
									continue
								}
								if ok, _ := c.ischange(context.TODO(), v.Name, ret); !ok {
									continue
								}
								err = c.put(context.TODO(), v.Name, ret, models.IO_TYPE_STRING)
								if err != nil {
									slog.Error(err.Error(), v.Name, v.Pr)
									continue
								}
							case "double":
								ret, err := drv.ReadDouble(context.TODO(), v.Pr)
								if err != nil {
									slog.Error(err.Error(), v.Name, v.Pr)
									continue
								}
								if ok, _ := c.ischange(context.TODO(), v.Name, ret); !ok {
									continue
								}
								err = c.put(context.TODO(), v.Name, ret, models.IO_TYPE_DOUBLE)
								if err != nil {
									slog.Error(err.Error(), v.Name, v.Pr)
									continue
								}
							default:
								slog.Error("dt 配置不对", v.DT, v.Name, v.Pr)
							}

						} else {
							slog.Error("driver is not driver.Drvapi %d", v.Drv, v.Name, v.Pr)
						}

					}
				}
			}
		}(k, iodf)
	}

}
func (c *IOServer) ReadInt(ctx context.Context, key string) (int32, error) {

	//slog.Debug("IOServer readint", key)
	bkey := []byte(key)
	existingValue, err := c.get(bkey)
	if err != nil {

		return 0, err
	}
	if value, ok := existingValue.(int32); ok {
		return value, nil
	}
	return 0, fmt.Errorf("value with key %s is not an int", key)
}

func (c *IOServer) ReadString(ctx context.Context, key string) (string, error) {
	bkey := []byte(key)
	existingValue, err := c.get(bkey)
	if err != nil {
		return "", err
	}
	if value, ok := existingValue.(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("value with key %s is not a string", key)
}

func (c *IOServer) ReadDouble(ctx context.Context, key string) (float64, error) {
	bkey := []byte(key)
	existingValue, err := c.get(bkey)
	if err != nil {
		return 0, err
	}
	if value, ok := existingValue.(float64); ok {
		return value, nil
	}
	return 0, fmt.Errorf("value with key %s is not a float64", key)
}

//	func (c *IOServer) ReadBool(ctx context.Context, key string) (bool, error) {
//		bkey := []byte(key)
//		existingValue, err := c.get(bkey)
//		if err != nil {
//			return false, err
//		}
//		if value, ok := existingValue.(bool); ok {
//			return value, nil
//		}
//		return false, fmt.Errorf("value with key %s is not a bool", key)
//	}
func (c *IOServer) WriteInt(ctx context.Context, key string, value int32) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	slog.Debug("IOServer WriteInt", key, value)
	// 判断是否是io变量或者内存变量，先写入io在写入缓存
	if ok, _ := c.ischange(context.TODO(), key, value); !ok {

		return nil
	}

	if v, ok := c.cfg.Items[key]; ok {
		if v.DT != "int" {
			return fmt.Errorf("io 配置数据类型不匹配 dt=%s %s:%v", v.DT, key, value)
		}
		// 判断大小值是否合法

		min, err := v.GetMin()
		if err != nil {
			return err
		}

		if min != nil && value < min.(int32) {
			return fmt.Errorf("value %d is less than minimum allowed value %d", value, min.(int32))
		}

		max, err := v.GetMax()
		if err != nil {
			return err
		}
		if max != nil && value > max.(int32) {
			return fmt.Errorf("value %d is greater than maximum allowed value %d", value, max.(int32))
		}

		// 写io
		if v.Cat == "IO" {
			if drv, ok := CX.Modules.Items[v.Drv].(api.Drvapi); ok {
				err = drv.WriteInt(ctx, v.Pw, value)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("driver is not driver.Drvapi %d", v.Drv)
			}
		}
	}

	err := c.put(ctx, key, value, models.IO_TYPE_INT)
	if err != nil {
		return err
	}
	return nil
}

func (c *IOServer) SetState(ctx context.Context, key string, value string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	slog.Debug("IOServer WriteState", key, value)
	ok, old := c.ischange(context.TODO(), key, value)

	// 原先等于 IDLE 时才可以修改

	// if old == "" {
	// 	old = "IDLE"
	// }

	if !ok {
		if value != "IDLE" {
			return fmt.Errorf("%s state is %s", key, old)
		}
		return nil
	} else {
		if old != "IDLE" && value != "IDLE" {
			return fmt.Errorf("%s state is %s", key, old)
		}
	}
	if v, ok := c.cfg.Items[key]; ok {
		if v.DT != "string" {
			return fmt.Errorf("io 配置数据类型不匹配 dt=%s %s:%v", v.DT, key, value)
		}
		// 判断大小值是否合法

		// 写io
		if v.Cat == "IO" {
			if drv, ok := CX.Modules.Items[v.Drv].(api.Drvapi); ok {
				err := drv.WriteString(ctx, v.Pw, value)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("driver is not driver.Drvapi %d", v.Drv)
			}
		}
	}

	err := c.put(ctx, key, value, models.IO_TYPE_STRING)
	if err != nil {
		return err
	}
	return nil
}
func (c *IOServer) WriteString(ctx context.Context, key string, value string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	slog.Debug("IOServer WriteString", key, value)
	// 判断是否是io变量或者内存变量，先写入io在写入缓存
	ok, _ := c.ischange(context.TODO(), key, value)
	if !ok {
		return nil
	}

	if v, ok := c.cfg.Items[key]; ok {
		if v.DT != "string" {
			return fmt.Errorf("io 配置数据类型不匹配 dt=%s %s:%v", v.DT, key, value)
		}
		// 判断大小值是否合法

		// 写io
		if v.Cat == "IO" {
			if drv, ok := CX.Modules.Items[v.Drv].(api.Drvapi); ok {
				err := drv.WriteString(ctx, v.Pw, value)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("driver is not driver.Drvapi %d", v.Drv)
			}
		}
	}

	err := c.put(ctx, key, value, models.IO_TYPE_STRING)
	if err != nil {
		return err
	}
	return nil
}

func (c *IOServer) WriteDouble(ctx context.Context, key string, value float64) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	slog.Debug("IOServer WriteDouble", key, value)
	// 判断是否是io变量或者内存变量，先写入io在写入缓存
	if ok, _ := c.ischange(context.TODO(), key, value); !ok {

		return nil
	}

	if v, ok := c.cfg.Items[key]; ok {
		if v.DT != "double" {
			return fmt.Errorf("io 配置数据类型不匹配 dt=%s %s:%v", v.DT, key, value)
		}
		// 判断大小值是否合法
		min, err := v.GetMin()
		if err != nil {
			return err
		}
		if min != nil && value < min.(float64) {
			return fmt.Errorf("value %f is less than minimum allowed value %f", value, min.(float64))
		}
		max, err := v.GetMax()
		if err != nil {
			return err
		}
		if max != nil && value > max.(float64) {
			return fmt.Errorf("value %f is greater than maximum allowed value %f", value, max.(float64))
		}

		// 写io
		if v.Cat == "IO" {
			if drv, ok := CX.Modules.Items[v.Drv].(api.Drvapi); ok {
				err = drv.WriteDouble(ctx, v.Pw, value)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("driver is not driver.Drvapi %d", v.Drv)
			}
		}
	}
	err := c.put(ctx, key, value, models.IO_TYPE_DOUBLE)
	if err != nil {
		return err
	}
	return nil
}

// func (c *IOServer) WriteBool(ctx context.Context, key string, value bool) error {
// 	slog.Debug("IOServer WriteBool", key, value)
// 	// 判断是否是io变量或者内存变量，先写入io在写入缓存
// 	if !c.ischange(context.TODO(), key, value) {
// 		return nil
// 	}

// 	if v, ok := c.cfg.Items[key]; ok {
// 		if v.DT != "int" {
// 			return fmt.Errorf("io 配置数据类型不匹配 dt=%s %s:%v", v.DT, key, value)
// 		}

// 		// 写io
// 		if v.Cat == "IO" {
// 			if drv, ok := CX.Modules.Items[v.Drv].(driver.Drvapi); ok {

// 				i := 0
// 				if value {
// 					i = 1
// 				}
// 				err := drv.WriteInt(ctx, v.Pw, i)
// 				if err != nil {
// 					return err
// 				}
// 			} else {
// 				return fmt.Errorf("driver is not driver.Drvapi %d", v.Drv)
// 			}
// 		}
// 	}

//		err := c.put(ctx, key, value, IO_TYPE_BOOL)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
func (c *IOServer) DeleteFromPrefix(ctx context.Context, prefix string) error {
	if c.db == nil {
		return errors.New("database is not initialized")
	}
	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()
	for iter.Next() {
		key := iter.Key()
		// 删除符合条件的数据
		if err := c.db.Delete(key, nil); err != nil {
			return err
		}
	}
	if err := iter.Error(); err != nil {
		return err
	}
	return nil
}
func (c *IOServer) ReadFromPrefix(ctx context.Context, prefix string) (map[string]any, error) {
	if c.db == nil {
		return nil, errors.New("database is not initialized")
	}

	result := make(map[string]interface{})

	iter := c.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	defer iter.Release()

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		// 将符合条件的数据添加到结果中
		result[string(key)] = models.ConvertValue(value)
	}

	if err := iter.Error(); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *IOServer) ischange(ctx context.Context, key string, value any) (bool, any) {

	bkey := []byte(key)

	// 检查键是否已经存在
	existingValue, err := c.get(bkey)

	if err == nil && existingValue == value {
		// 键已经存在且值相同，不需要写入
		return false, existingValue
	}

	return true, existingValue
}
func (c *IOServer) put(ctx context.Context, key string, value any, iotype models.IO_TYPE) error {
	bkey := []byte(key)
	// 将键值对写入数据库

	// 使用binary.Write 分别将value转化成byte数组，iotype占用第一位
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, iotype)
	switch iotype {
	case models.IO_TYPE_STRING:
		if len(value.(string)) == 0 {
			err := c.db.Delete(bkey, nil)
			if err != nil {
				return err
			}
			c.bus.Trigger(key, map[string]any{key: value})

			return nil
		}
		buf.WriteString(value.(string))
	default:
		binary.Write(&buf, binary.BigEndian, value)
		// if err != nil {
		// 	slog.Error(err.Error())
		// }
	}

	err := c.db.Put(bkey, buf.Bytes(), nil)
	if err != nil {
		return err
	}
	if c.cfg.Items[key].Rcd == 1 {
		// s
		c.writeChan <- &models.His{Parm: key, Value: value, CreateTime: time.Now().UnixNano()}
		// err := c.his.Insert(&models.His{Parm: key, Value: value})
		// if err != nil {
		// 	slog.Error(err.Error())
		// }
	}
	// events
	c.bus.Trigger(key, map[string]any{key: value})

	return nil
}

func (c *IOServer) get(key []byte) (any, error) {
	existingValue, err := c.db.Get(key, nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			if io, ok := c.cfg.Items[string(key)]; ok {
				return io.GetDfval()
			}
		}
		return nil, fmt.Errorf("%s,%s ", err.Error(), string(key))
	}
	v := models.ConvertValue(existingValue)
	return v, nil
}
