package alarm

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/api2"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm/socketio"
	"github.com/chenxiio/chenxi/logger"
	"github.com/chenxiio/chenxi/models"
	"github.com/chenxiio/go-eventbus"
	socketios "github.com/googollee/go-socket.io"
	"golang.org/x/exp/slog"
)

type Alarms struct {
	AlarmsDao
	bus      *eventbus.Bus[any]
	lock     sync.Mutex
	socketio *socketio.SocketioServer
	cfg      *cfg.Alarms
	api2.AlmWaitAck
}

type IOEvent struct {
	// i *int
	//dataChannel chan map[string]interface{}
	topic string
	s     socketios.Conn
	alms  *Alarms
}

func (n *IOEvent) Sub(k string) {
	n.alms.bus.On(k, n)
}
func (n *IOEvent) UnSub(k string) {
	n.alms.bus.Off(k, n)
}
func (n *IOEvent) Dispatch(data ...any) {

	go n.s.Emit("subalm"+n.topic, data[0])
	// for key, value := range data[0] {
	// 	fmt.Printf("%s: %v\n", key, value)
	// }
}
func NewAlarms(path string, _cfg *cfg.Alarms, socketio *socketio.SocketioServer, _log *logger.Logger) *Alarms {
	//	path string, _cfg *cfg.IOCfg, socketio *socketio.SocketioServer
	log = _log
	bus := eventbus.New[any]()
	if _cfg == nil {
		_cfg = &cfg.Alarms{Alarms: make(cfg.AlmItems)}
	}
	db, err := sql.Open("sqlite3", chenxi.CX.Cfg.Basedir+"data/alm.db")
	if err != nil {
		log.Error("Failed to connect to database: ", "err", err.Error())
		panic(err)
	}
	ios := &Alarms{
		AlarmsDao: *InitAlarmsDaoInstance(db),
		cfg:       _cfg,
		bus:       bus,
		socketio:  socketio,
	}
	ios.AlmWaitAck = api2.AlmWaitAck{Socketioapi: ios}
	ios.socketio.OnEvent("/", "subalm", func(s socketios.Conn, prexx string) {
		slog.Debug("subalm", s.ID(), prexx)
		ie := &IOEvent{topic: prexx, s: s, alms: ios}

		ios.socketio.Sub(s.ID(), prexx, ie)

		// 返回该前缀所有数据
		//s.Emit("reply", "have "+prexx)

		// mp, err := ios.ReadFromPrefix(context.TODO(), prexx)
		// if err != nil {
		// 	slog.Error(err.Error())
		// 	return
		// }
		alms, err := ios.GetAlarms(context.TODO())
		if err != nil {
			log.Error(err.Error())
		}

		s.Emit("subalm"+prexx, alms)
	})

	ios.socketio.OnEvent("/", "unsubalm", func(s socketios.Conn, prexx string) {
		slog.Debug("unsubalm", s.ID(), prexx)

		ios.socketio.UnSub(s.ID(), prexx)
	})
	return ios
}

func (c *Alarms) Unsub(group string, key string, e socketio.SocketsubDispatch) error {
	slog.Debug("Alarms Unsub ", "", key)
	c.bus.Off(key, e.(eventbus.Event[any]))
	return nil
}

// key =subio
func (c *Alarms) Sub(key string, parm string, f socketio.SocketsubDispatch) error {
	slog.Debug("Alarms Sub ", key, parm)
	c.bus.On(parm, f.(eventbus.Event[any]))

	alms, err := c.GetAlarms(context.TODO())
	if err != nil {
		log.Error(err.Error())
		return err
	}
	f.(eventbus.Event[any]).Dispatch(alms)
	return nil
}

func (c *Alarms) Insert(ctx context.Context, alarm *models.Alarm) error {

	err := c.AlarmsDao.Insert(alarm)
	if err != nil {
		return err
	}
	c.bus.Trigger(fmt.Sprintf("alm.%d", alarm.Sn), alarm)

	return nil
}

func (c *Alarms) ClearAlarms(ctx context.Context, aid int64, cleartype int) error {

	err := c.AlarmsDao.ClearAlarms(aid, cleartype)
	if err != nil {
		return err
	}
	c.bus.Trigger(fmt.Sprintf("alm.%d", aid), models.Alarm{Sn: aid, ClearType: cleartype})

	return nil
}

// GetAlarms(ctx context.Context) ([]models.Alarm, error)                                //perm:none
// GetAlarmsHistory(ctx context.Context, start int64, end int64) ([]models.Alarm, error) //perm:none

var _ api2.ALMApi = (*Alarms)(nil)
