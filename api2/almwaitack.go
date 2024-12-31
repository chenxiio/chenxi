package api2

import (
	"fmt"
	"sync"

	"github.com/chenxiio/chenxi/models"
)

type AlmWaitAck struct {
	Socketioapi
}

type AlmWaitAckEvent struct {
	aid  int64
	lock sync.Mutex
	alm  models.Alarm
}

func (a *AlmWaitAckEvent) Dispatch(data ...any) {
	for _, v := range data[0].([]models.Alarm) {
		//a.once.
		if v.Sn == a.aid && v.ClearType != 0 {
			a.alm = v
			a.lock.Unlock()
		}
	}
}

func (a *AlmWaitAck) WaitAlmAck(aid int64) (int, error) {
	e := AlmWaitAckEvent{aid: aid}
	e.lock.Lock()
	err := a.Sub("subalm", fmt.Sprintf("alm.%d", aid), &e)
	if err != nil {
		return 0, err
	}
	e.lock.Lock()
	e.lock.Unlock()
	return e.alm.ClearType, nil
}
