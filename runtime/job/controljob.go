package job

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/chenxiio/chenxi/models"
)

// Control Job Management Standard (E94)
// CIS: Control Job State Model
// Q：OUEUED
// S: SELECTED
// WFS:WAITING FOR START
// E: EXECUTING
// P:PAUSED
// C:COMPLETED
// QUEUED，EXECUTING，PAUSED,ABORTED,ENDING,COMPLETED
// - CONTROL_JOB_ID:控制作业的唯一ID
// - CARRIER_ID:对应的载具ID
// - PRIORITY:控制作业的优先级
// - PROCESS_JOB_LIST:包含的工艺作业列表
// - EXECUTE_STATUS:预期的执行状态
// - Mode        : // Auto ,Manual

type ControlJob struct {
	// CONTROL_JOB_ID   string
	// CARRIER_ID       []string
	// PRIORITY         string
	// PROCESS_JOB_LIST []string
	// Mode             string // Auto ,Manual
	// State            string
	// Createtime       int64
	// Starttime        int64
	// Endtime          int64
	models.ControlJob
	processJobs []*ProcessJob
	lock        sync.Mutex
}

func (c *ControlJob) AddPj(act *ProcessJob) {
	c.lock.Lock()
	defer c.lock.Unlock()
	log.Debug("call ControlJob AddAction", c.CONTROL_JOB_ID, act.PROCESS_JOB_ID)
	c.processJobs = append(c.processJobs, act)
	// // src-dst 优先级
	sort.Slice(c.processJobs, func(i, j int) bool {
		return c.processJobs[i].PRIORITY > c.processJobs[j].PRIORITY
	})

}

func (c *ControlJob) Start(ctx context.Context, parm string) error {
	log.Debug("call ControlJob Start", c.CONTROL_JOB_ID, parm)
	// 如果发生错误，状态设置Error
	c.lock.Lock()
	defer c.lock.Unlock()

	sort.Slice(c.processJobs, func(i, j int) bool {
		return c.processJobs[i].PRIORITY > c.processJobs[j].PRIORITY
	})
	//

	// 先load 再createcj
	// processjob 判断是否load 再次load，根据业务情况，设置load状态，可控制是否需要二次load map

	// for _, v := range c.CARRIER_ID {

	// 	load, _ := chenxi.CX.IOServer.ReadString(context.TODO(), "env.carrier."+v)
	// 	// if err != nil {
	// 	// 	return fmt.Errorf("Load 不正确，请重新load,%s,%s", v, err.Error())
	// 	// }
	// 	if load == {
	// 		//
	// 		// return fmt.Errorf("Load 不正确，请重新load,%s,%s", v, err.Error())
	// 		log.Debug("called Load", v, "")
	// 		err := jobIns.cjpms[v].loadandmap(ctx, "")
	// 		if err != nil {
	// 			return fmt.Errorf("Load 不正确，请重新load,%s,%s", v, err.Error())
	// 		}

	// 	}
	// }
	jobIns.state = "RUN"
	c.State = "EXECUTING"
	c.Starttime = time.Now().UnixNano()
	err := c.processJobs[0].Start(ctx, "")
	if err != nil {
		//c.processJobs[0].Abort(context.TODO(), "")
		return err
	}
	err = ControlJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	return nil
}

// func (c *ControlJob) SetError(ctx context.Context, parm string) error {
// 	c.State = "ABORT"

//		err := ControlJobDAOInstance.Update(c)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
func (c *ControlJob) End(ctx context.Context, parm string) error {
	log.Debug("call ControlJob End", c.CONTROL_JOB_ID, parm)
	// end
	c.State = "ENDING"

	err := ControlJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	return nil
}
func (c *ControlJob) Pause(ctx context.Context, parm string) error {
	log.Debug("call ControlJob Pause", c.CONTROL_JOB_ID, parm)
	if c.State != "EXECUTING" {
		//
		return fmt.Errorf("ControlJob Pause error C.state:%s", c.State)

	}
	c.State = "PAUSED"
	for _, v := range c.processJobs {
		if v.State == "EXECUTING" {
			err := v.Pause(ctx, "")
			if err != nil {
				log.Error(err.Error())
			}
		}

	}

	err := ControlJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	return nil
}
func (c *ControlJob) Resume(ctx context.Context, parm string) error {
	log.Debug("call ControlJob Resume", c.CONTROL_JOB_ID, parm)
	if c.State != "PAUSED" {
		//
		return fmt.Errorf("ControlJob Pause error C.state:%s", c.State)

	}
	c.State = "EXECUTING"
	for _, v := range c.processJobs {
		if v.State == "PAUSED" {
			err := v.Resume(ctx, "")
			if err != nil {
				log.Error(err.Error())
			}
		}

	}
	err := ControlJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	//wait
	return nil
}
func (c *ControlJob) Abort(ctx context.Context, parm string) error {
	log.Debug("call ControlJob Abort", c.CONTROL_JOB_ID, parm)
	if c.State != "PAUSED" && c.State != "EXECUTING" {
		//
		return fmt.Errorf("ControlJob Abort error C.state:%s", c.State)

	}
	c.State = "ABORTED"
	for _, v := range c.processJobs {
		if v.State == "PAUSED" || v.State == "EXECUTING" {
			err := v.Abort(ctx, "")
			if err != nil {
				log.Error(err.Error())
			}
		}

	}
	err := ControlJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	return nil
}

// func (c *ControlJob) End(ctx context.Context, parm string) error {

//		return nil
//	}
func (c *ControlJob) List(ctx context.Context, parm string) error {

	return nil
}
