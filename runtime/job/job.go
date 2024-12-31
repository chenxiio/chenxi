package job

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm"
	"github.com/chenxiio/chenxi/logger"
	"github.com/chenxiio/chenxi/models"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type Job struct {
	ControlJobs   []*ControlJob
	tmpProcessjob []*ProcessJob
	cjpms         map[string]*UnitjobControl
	cjtms         map[string]*UnitjobTM
	// cjcms         map[string]*UnitjobCM
	isauto bool
	ujs    [][]*Unitjob
	state  string //INIT IDLE RUN  DOWN PM PAUSED
	cond   *sync.Cond
	iswait bool
	cjlock sync.Mutex
}

var jobonce sync.Once

var jobIns *Job

// GetInstance 返回Job的单例实例
func CreateJobInstance() *Job {
	jobonce.Do(func() {
		log = logger.GetLog("job", "", chenxi.CX.Cfg.Basedir)
		if chenxi.CX.Name != "ioserver" {
			log.Error("job 必须和ioserver同一个进程 ")
			panic(fmt.Errorf("job 必须和ioserver同一个进程 "))
		}
		jobIns = &Job{
			ControlJobs:   []*ControlJob{},
			tmpProcessjob: []*ProcessJob{},
			cond:          sync.NewCond(&sync.Mutex{}),
			ujs:           make([][]*Unitjob, 0),
			iswait:        false,
			isauto:        true,
			state:         "INIT",
			cjpms:         map[string]*UnitjobControl{},
			cjtms:         map[string]*UnitjobTM{},
			// cjcms:         map[string]*UnitjobCM{},
		}
		// 加载数据库

		err := os.MkdirAll(chenxi.CX.Cfg.Basedir+"data/", os.ModePerm)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}
		db, err := sql.Open("sqlite3", chenxi.CX.Cfg.Basedir+"data/job.db")
		if err != nil {
			log.Error("Failed to connect to database: ", "err", err.Error())
			panic(err)
		}
		InitUnitjobDAOInstance(db)
		InitProcessJobDAOInstance(db)
		InitControlJobDAOInstance(db)
		_, err = GetCTCInstance(&chenxi.CX.Cfg.CTCCfg)
		if err != nil {
			log.Error(err.Error())
			panic(err)
		}

		// 最近七天
		cjlist, err := ControlJobDAOInstance.GetList(0, 0, nil, map[string]any{"State": "COMPLETED"})
		if err != nil {
			log.Error(err.Error())

		}
		jobIns.ControlJobs = cjlist

		for _, v := range cjlist {
			v.State = "ABORTED"
			for _, p := range v.PROCESS_JOB_LIST {
				pj, err := ProcessJobDAOInstance.GetByPJID(p)
				if err != nil {
					log.Error(err.Error())
					panic(err)
				}
				pj.State = "ABORTED"
				ujs, err := UnitjobDAOInstance.GetByPJID(p)
				if err != nil {
					log.Error(err.Error())
					panic(err)
				}

				pj.unitjobs = ujs
				for _, u := range ujs {
					u.pj = pj
					u.State = "ABORTED"
					pj.allocated += len(u.MATERIAL_LIST)
				}
				v.processJobs = append(v.processJobs, pj)
			}
		}
		for _, v := range chenxi.CX.Cfg.Modules.Items {

			switch v.API {

			case "TM":
				tuj, err := NewUnitjobTM(v.Name, v.Type)
				if err != nil {
					log.Error(err.Error())
					panic(err)
				}
				jobIns.cjtms[v.Name] = tuj
			case "PM":
				fallthrough
			case "CM":
				puj, err := NewUnitjobControl(v.Name, v.API)
				if err != nil {
					log.Error(err.Error())
					panic(err)
				}
				jobIns.cjpms[v.Name] = puj
			default:

			}

		}
		for tm := range CTCIns.cfg.Group {
			if strings.HasPrefix(tm, "_MV") {

				tuj, err := NewUnitjobTM(tm, "_MV")
				if err != nil {
					log.Error(err.Error())
					panic(err)
				}
				jobIns.cjtms[tm] = tuj
			}
		}

		go jobIns.executing()
	})
	return jobIns
}
func (j *Job) checkwafer(v cfg.Module) error {
	// 检查 是否有wafer 在单元中
	switch v.API {
	case "PM":
		fallthrough
	case "TM":
		ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), v.Name+".wid.")
		if err != nil {
			//log.Error(err.Error())
			return err
		}
		fmt.Println(ws)
		for _, wid := range ws {
			if len(wid.(string)) > 0 {
				return fmt.Errorf("%s 有wafer 不能初始化", v.Name)
			}
		}
	case "CM":

	default:
	}

	return nil
}
func (j *Job) checkall() error {

	for _, v := range j.cjpms {
		err := v.init("")
		if err != nil {

			return err
		}
		err = j.checkwafer(*v.unitcfg)
		if err != nil {
			return err
		}
	}
	// for _, v := range j.cjcms {
	// 	err := v.init("")
	// 	if err != nil {

	// 		return err
	// 	}
	// 	err = j.checkwafer(*v.unitcfg)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	for _, v := range j.cjtms {
		err := v.init("")
		if err != nil {
			return err
		}
		err = j.checkwafer(*v.unitcfg)
		if err != nil {
			return err
		}
	}
	CTCIns.InterlockingUnLockAll()
	return nil
}
func (j *Job) clearallcj() error {
	j.cjlock.Lock()
	defer j.cjlock.Unlock()

	//  删除所有job

	for _, v := range j.ControlJobs {
		v.ControlJob.State = "COMPLETED"
		err := ControlJobDAOInstance.Update(v)
		if err != nil {
			return err
		}
	}
	j.ControlJobs = make([]*ControlJob, 0)
	return nil
}
func (j *Job) Init(ctx context.Context, parm string) error {
	log.Info("job init called")
	if err := j.checkall(); err != nil {
		j.state = "DOWN"
		return err
	}

	// 设置设备状态 READY
	err := j.clearallcj()
	if err != nil {
		j.state = "DOWN"
		return err
	}
	j.state = "IDLE"
	for i := 0; i < len(j.ujs); i++ {
		j.ujs[i] = make([]*Unitjob, 0)
	}

	// 删除所有 env.wid.

	//j.cjcms = map[string]*UnitjobControl{}
	// err := j.checkwafer(v)
	// if err != nil {
	// 	log.Error(err.Error())
	// 	return err
	// }

	//  结束全部job

	testjob()
	return nil
}
func testjob() {
	// allinit
	// err := jobIns.Init(context.Background(), "")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	//jobIns.PJCreate()
	// 初始化cm
	//jobIns.cjpms["cmtest1"]
	// load and map
	jobIns.cjpms["pmtest2"].merge = 2
	err := jobIns.cjpms["cmtest1"].loadandmap(context.Background(), "")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = jobIns.cjpms["cmtest2"].loadandmap(context.Background(), "")
	if err != nil {
		fmt.Println(err)
		return
	}
	// 创建pj jobIns
	pj1 := models.ProcessJob{PROCESS_JOB_NAME: "testjob", PARAMETER_LIST: map[string]string{"main": "processtest"},
		MATERIAL_LIST: []string{"foup_cmtest1_1", "foup_cmtest1_2", "foup_cmtest2_1", "foup_cmtest2_2"}}
	// pj1.PARAMETER_LIST["main"] = "processtest"
	// pj2 := models.ProcessJob{PROCESS_JOB_NAME: "testjob", PARAMETER_LIST: map[string]string{"main": "processtest"},
	// 	MATERIAL_LIST: []string{"foup_cmtest2_1", "foup_cmtest2_2"}}
	// pj1.PARAMETER_LIST["main"] = "processtest"
	pj, err := jobIns.PJCreate(context.TODO(), &pj1)
	if err != nil {
		fmt.Println(err)
		return
	}
	// pj_2, err := jobIns.PJCreate(context.TODO(), &pj2)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// 创建 cj
	cj1 := models.ControlJob{CARRIER_ID: []string{"foup_cmtest1", "foup_cmtest2"},
		PROCESS_JOB_LIST: []string{pj.PROCESS_JOB_ID}, State: "QUEUED"}
	cj, err := jobIns.CJCreate(context.Background(), &cj1)
	if err != nil {
		fmt.Println(err)
		return
	}
	// start //
	_, err = jobIns.CJStart(context.Background(), cj.CONTROL_JOB_ID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (j *Job) Continue(ctx context.Context) error {

	for _, v := range j.cjtms {
		v.cond.L.Lock()
		v.cond.Signal()
		v.cond.L.Unlock()
	}
	jobIns.cond.L.Lock()
	jobIns.cond.Signal()
	jobIns.cond.L.Unlock()
	return nil
}
func (j *Job) AddAction(act *Unitjob) {
	j.cond.L.Lock()
	log.Debug("Call Job AddAction", act.Unit, act.StepName)
	for i := len(j.ujs); i <= act.StepName; i++ {
		j.ujs = append(j.ujs, make([]*Unitjob, 0))
	}

	j.ujs[act.StepName] = append(j.ujs[act.StepName], act)
	// // src-dst 优先级
	// sort.Slice(t.TMActions, func(i, j int) bool {
	// 	return t.TMActions[i].PickPriority > t.TMActions[j].PickPriority
	// })

	j.cond.Signal()
	// log.Debug("Signal end")
	j.cond.L.Unlock()
	// log.Debug("AddAction Unlock end")
}
func (j *Job) executing() {
	for {
		if !j.isauto {
			break
		}
		j.cond.L.Lock()
		if j.iswait {
			log.Debug("job waiting")
			j.cond.Wait()

		}
		j.iswait = true

		for i := len(j.ujs) - 1; i >= 0; i-- {
			sort.Slice(j.ujs[i], func(i1, j1 int) bool {
				return j.ujs[i][i1].Unitjob_id < j.ujs[i][j1].Unitjob_id
			})

			for j1 := 0; j1 < len(j.ujs[i]); j1++ {
				b := j.createuj(*j.ujs[i][j1])
				// if b {
				// 	j.iswait = false
				// }
				if b {
					//
					j.ujs[i] = append(j.ujs[i][:j1], j.ujs[i][j1+1:]...)
					j1--
				}
			}

		}
		j.cond.L.Unlock()
	}
}
func (j *Job) createuj(uj Unitjob) bool {
	if j.state == "ENDING" {
		return false
	}
	if uj.pj.State == "ENDING" {
		//
		return false
	}
	// ABROTED 、
	//PAUSED
	if uj.State == "ABROTED" || uj.State == "PAUSED" {
		return false
	}
	if uj.StepName == 0 {
		//
		err := uj.pj.UJCreate(context.TODO(), uj)
		if err != nil {
			log.Error(err.Error())
			// 如果程序错误，会出现循环报错的情况
			return false
		}
		if uj.pj.allocated == len(uj.pj.MATERIAL_LIST) {
			return true
		}
		return false
	} else {
		err := uj.pj.UJContinue(context.TODO(), uj)
		if err != nil {
			log.Warn(err.Error())
			return false
		}
		return true
	}

}

func (c *Job) Move2Home(ctx context.Context, from string, mats []string) error {

	return nil
}

// form to
func (c *Job) Move3(ctx context.Context, from, to string, mats []string, pj *ProcessJob) error {
	// prefix := "env.wid."
	// if home == 0 {
	// 	prefix = "env.wid.scan."
	// }
	moveActions := []*MoveActions{}

	for _, str := range mats {
		if len(str) > 0 {

			unitp2, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+str)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			_, slotname, err := comm.UnmarshalUSlot(unitp2)

			if err != nil {
				log.Error(err.Error())
				return err
			}

			// }
			var tmactions *MoveActions
			for _, v := range moveActions {

				if v.Cmds.From == from && v.Cmds.To == to {
					tmactions = v
					break
				}
			}
			if tmactions == nil {

				// actstr := fmt.Sprintf("%s-%s", from, to)

				tas, err := CTCIns.GenActions(from, to)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				tmactions = tas
				tmactions.pj = pj
				moveActions = append(moveActions, tmactions)
			}

			tmactions.Slots = append(tmactions.Slots,
				cfg.Slot{Name: slotname, WaferId: str})

		} else {
			return fmt.Errorf("waferid格式错误 :%s", str)
		}
	}
	//form unit 优先级 排序
	sort.Slice(moveActions, func(i, j int) bool {
		return moveActions[i].Cmds.Priority > moveActions[j].Cmds.Priority
	})
	return c.Move(ctx, moveActions)
}

// home:0 返回 1：起始位置 ，>0 step
func (c *Job) Move2(ctx context.Context, unit string, mats []string, home int, pj *ProcessJob) error {
	prefix := "env.wid."
	if home == 0 {
		prefix = "env.wid.scan."
	}
	moveActions := []*MoveActions{}

	for _, str := range mats {
		if len(str) > 0 {
			unitp, err := chenxi.CX.IOServer.ReadString(context.TODO(), prefix+str)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			uname, slotname, err := comm.UnmarshalUSlot(unitp)

			if err != nil {
				log.Error(err.Error())
				return err
			}
			if home == 0 {
				unitp2, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+str)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				_, slotname2, err := comm.UnmarshalUSlot(unitp2)

				if err != nil {
					log.Error(err.Error())
					return err
				}
				slotname = slotname2
			}
			var tmactions *MoveActions
			for _, v := range moveActions {
				fro := ""
				to := ""
				if home == 0 {
					fro = unit
					to = uname
				} else {
					fro = uname
					to = unit
				}
				if v.Cmds.From == fro && v.Cmds.To == to {
					tmactions = v
					break
				}
			}
			if tmactions == nil {
				// actstr := ""

				if home == 0 {
					// actstr = fmt.Sprintf("%s-%s", unit, uname)
					tmactions, err = CTCIns.GenActions(unit, uname)
				} else {
					// actstr = fmt.Sprintf("%s-%s", uname, unit)
					tmactions, err = CTCIns.GenActions(uname, unit)
				}

				// tas, err := CTCIns.GenActions(actstr)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				// tmactions = tas
				tmactions.pj = pj
				moveActions = append(moveActions, tmactions)
			}

			//cfgslot := chenxi.CX.Cfg.Modules.CMcfgs[uname].Slots[slotname-1]
			tmactions.Slots = append(tmactions.Slots,
				cfg.Slot{Name: slotname, WaferId: str})

			if home == 1 && CTCIns.cfg.ReturningMode == 0 && !pj.issubpj {
				// src 起始位置, pj 结束时 清楚
				err := chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.scan.%s", str), unitp)
				if err != nil {
					return err
				}
			}

		} else {
			return fmt.Errorf("waferid格式错误 :%s", str)
		}
	}
	//form unit 优先级 排序
	sort.Slice(moveActions, func(i, j int) bool {
		return moveActions[i].Cmds.Priority > moveActions[j].Cmds.Priority
	})
	return c.Move(ctx, moveActions)
}
func (c *Job) Move(ctx context.Context, moveActions []*MoveActions) error {

	// 按照顺序发送执行指令给tm，tm根据优先级执行搬运动作

	for _, v := range moveActions {
		mv := MoveActions{Cmds: MoveActionCmds{From: v.Cmds.From, To: v.Cmds.To, Priority: v.Cmds.Priority},
			Slots: v.Slots, Lock: v.Lock, Lockfrom: v.Lockfrom, Lockto: v.Lockto, pj: v.pj}
		tmActions := []TMAction{}
		for _, act := range v.Cmds.Actions {
			act.moveActions = &mv
			tmActions = append(tmActions, act)
		}
		// from  and tm 优先级
		sort.Slice(tmActions, func(i, j int) bool {
			return tmActions[i].Priority > tmActions[j].Priority
		})

		for _, tv := range tmActions {
			jobIns.cjtms[tv.Name].AddAction(tv)
		}
	}
	return nil
}
func (c *Job) Auto(ctx context.Context) error {
	log.Debug("Call Job Auto")
	if c.state != "IDLE" {
		return fmt.Errorf("Job Auto error state ：%s", c.state)
	}
	if len(c.ControlJobs) == 0 {
		return fmt.Errorf("Job Auto error controlsjob count =0")
	}
	c.isauto = true

	c.state = "RUN"

	// if err != nil {
	// 	return  err
	// }
	// return c.CJList(ctx)
	return c.ControlJobs[0].Start(ctx, "")
}

func (c *Job) Pause(ctx context.Context) error {
	log.Debug("Call Job Pause")
	if c.state != "RUN" {
		return fmt.Errorf("Job Pause error state:%s", c.state)
	}
	c.state = "PAUSED"
	for _, v := range c.cjtms {

		if v.state == "IDLE" {
			err := v.pause(ctx, "")
			if err != nil {
				return err
			}
		}

	}
	for _, v := range c.ControlJobs {
		if v.State == "EXECUTING" {
			err := v.Pause(ctx, "")
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}
	return nil
}
func (c *Job) Resume(ctx context.Context) error {
	log.Debug("Call Job Resume")
	if c.state != "PAUSED" {
		return fmt.Errorf("Job Pause error state:%s", c.state)
	}
	c.state = "RUN"
	for _, v := range c.cjtms {

		if v.state == "PAUSED" {
			err := v.resume(ctx, "")
			if err != nil {
				return err
			}
		}

	}
	for _, v := range c.ControlJobs {
		if v.State == "PAUSED" {
			err := v.Resume(ctx, "")
			if err != nil {
				log.Error(err.Error())
				return err
			}
		}
	}
	return nil
}

func (c *Job) Abort(ctx context.Context) error {

	log.Debug("Call Job Abort")
	// if c.state != "RUN" && c.state != "PAUSED" {
	// 	return fmt.Errorf("Job Abort error state:%s", c.state)
	// }
	c.state = "ABORTED"
	for _, v := range c.cjtms {

		if v.state == "IDLE" || v.state == "PAUSED" {
			err := v.abort(ctx, "")
			if err != nil {
				return err
			}
		}

	}
	for _, v := range c.ControlJobs {
		//
		if v.State == "PAUSED" || v.State == "EXECUTING" {
			//
			err := v.Abort(context.TODO(), "")
			if err != nil {
				return err
			}
		}
		//ret = append(ret, &v.ControlJob)
	}
	// for _, v := range c.cjpms {
	// 	err := v.abort(ctx, "")
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
func (c *Job) End(ctx context.Context) error {
	log.Debug("Call Job End")
	//ret := []*models.ControlJob{}
	for _, v := range c.ControlJobs {
		//
		err := v.End(context.TODO(), "")
		if err != nil {
			return err
		}
		//ret = append(ret, &v.ControlJob)
	}

	return nil
}

func (c *Job) PJCreate(ctx context.Context, pj1 *models.ProcessJob) (*models.ProcessJob, error) {
	log.Debug("Call Job PJCreate")
	// 判断cm 是不是load，map 以及结果是否一致
	// cm 初始化时选择是否重置load map 状态，如果不是load
	//
	if c.state != "IDLE" && c.state != "RUN" {
		return nil, fmt.Errorf("设备状态:%s", c.state)
	}
	cm := map[string]bool{}
	for _, v := range pj1.MATERIAL_LIST {

		unitp, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+v)
		if err != nil {
			return nil, fmt.Errorf("map 不正确，请重新map,%s,%s", unitp, err.Error())
		}
		us := strings.Split(unitp, ".")
		if len(us) < 2 {
			return nil, fmt.Errorf("map 不正确，请重新map,%s,%s", unitp, err.Error())
		}
		if _, ok := cm[us[0]]; ok {
			continue
		}
		load, err := chenxi.CX.IOServer.ReadInt(context.TODO(), us[0]+".load")
		if err != nil || load != 1 {
			// 重新load
			return nil, fmt.Errorf("Load 不正确，请重新Load,%s,%s", unitp, err.Error())
		}

		cm[us[0]] = true
	}
	pj := ProcessJob{ProcessJob: *pj1, unitjobs: []*Unitjob{}}
	pj.State = "QUEUED"
	if len(pj.PROCESS_JOB_ID) == 0 {
		pj.PROCESS_JOB_ID = uuid.New().String()
	}
	pj.Createtime = time.Now().UnixNano()
	// err := ProcessJobDAOInstance.Insert(&pjm)
	// if err != nil {
	// 	return nil, err
	// }
	log.Debug("Call Job PJCreate", pj.PROCESS_JOB_NAME, pj.PROCESS_JOB_ID)
	c.tmpProcessjob = append(c.tmpProcessjob, &pj)
	return &pj.ProcessJob, nil
}

func (c *Job) CJCreate(ctx context.Context, cj1 *models.ControlJob) (*models.ControlJob, error) {
	log.Debug("Call Job CJCreate")
	if c.state != "IDLE" && c.state != "RUN" {
		return nil, fmt.Errorf("设备状态:%s", c.state)
	}
	cj := &ControlJob{ControlJob: *cj1, processJobs: []*ProcessJob{}}

	// 删除 tmppj
	for k, v := range cj1.PROCESS_JOB_LIST {
		for i := 0; i < len(c.tmpProcessjob); i++ {
			if c.tmpProcessjob[i].PROCESS_JOB_ID == v {
				err := ProcessJobDAOInstance.Insert(c.tmpProcessjob[i])
				if err != nil {
					return nil, err
				}
				cj.processJobs = append(cj.processJobs, c.tmpProcessjob[i])
				c.tmpProcessjob = append(c.tmpProcessjob[:i], c.tmpProcessjob[i+1:]...)
				break
			}
		}
		if len(cj.processJobs) == k {
			return nil, fmt.Errorf("pj 不存在：%s", v)
		}
	}

	// sort.Slice(cj.processJobs, func(i, j int) bool {
	// 	return cj.processJobs[i].PRIORITY > cj.processJobs[j].PRIORITY
	// })
	c.cjlock.Lock()
	defer c.cjlock.Unlock()
	c.ControlJobs = append(c.ControlJobs, cj)

	sort.Slice(c.ControlJobs, func(i, j int) bool {
		return c.ControlJobs[i].PRIORITY > c.ControlJobs[j].PRIORITY
	})
	if cj.State == "EXECUTING" {
		//
		c.state = "RUN"
		err := cj.Start(context.TODO(), "")

		if err != nil {
			log.Error(err.Error())

			return cj1, err
		}
	} else {
		cj.State = "QUEUED"
	}
	err := ControlJobDAOInstance.Insert(cj)
	if err != nil {
		return nil, err
	}

	log.Debug("Call Job CJCreate", cj.CONTROL_JOB_ID, cj.CARRIER_ID)

	return &cj.ControlJob, nil
}

func (c *Job) CJStart(ctx context.Context, cjid string) (*models.ControlJob, error) {
	log.Debug("Call Job CJStart", cjid)
	c.state = "RUN"
	// if c.state == "RUN" {
	// 	//
	// }
	for _, v := range c.ControlJobs {
		if v.CONTROL_JOB_ID == cjid {
			err := v.Start(ctx, "")
			if err != nil {
				log.Error(err.Error())
				//err = v.Abort(context.TODO(), "")
				return &v.ControlJob, err
			}
			return &v.ControlJob, nil
		}
	}

	return nil, fmt.Errorf("Job CJStart not fond :%s", cjid)
}

func (c *Job) CJEnd(ctx context.Context, cjid string) (*models.ControlJob, error) {
	log.Debug("Call Job CJEnd", cjid)
	for _, v := range c.ControlJobs {
		if v.CONTROL_JOB_ID == cjid {
			err := v.End(ctx, "")
			if err != nil {
				log.Error(err.Error())
				//err = v.Abort(context.TODO(), "")
				return &v.ControlJob, err
			}
		}
	}

	return nil, fmt.Errorf("Job CJStop not fond :%s", cjid)
}

func (c *Job) CJPause(ctx context.Context, cjid string) (*models.ControlJob, error) {
	log.Debug("Call Job CJPause", cjid)
	for _, v := range c.ControlJobs {
		if v.CONTROL_JOB_ID == cjid {

			err := v.Pause(ctx, "")
			if err != nil {
				log.Error(err.Error())
				//err = v.Abort(context.TODO(), "")
				return &v.ControlJob, err
			}

		}
	}

	return nil, fmt.Errorf("Job CJPause not fond :%s", cjid)
}

func (c *Job) CJResume(ctx context.Context, cjid string) (*models.ControlJob, error) {
	log.Debug("Call Job CJResume", cjid)
	for _, v := range c.ControlJobs {
		if v.CONTROL_JOB_ID == cjid {
			err := v.Resume(ctx, "")
			if err != nil {
				log.Error(err.Error())
				//err = v.Abort(context.TODO(), "")
				return &v.ControlJob, err
			}
		}
	}

	return nil, fmt.Errorf("Job CJResume not fond :%s", cjid)
}

func (c *Job) CJAbort(ctx context.Context, cjid string) (*models.ControlJob, error) {
	log.Debug("Call Job CJAbort", cjid)
	for _, v := range c.ControlJobs {
		if v.CONTROL_JOB_ID == cjid {
			err := v.Abort(ctx, "")
			if err != nil {
				log.Error(err.Error())
				//err = v.Abort(context.TODO(), "")
				return &v.ControlJob, err
			}
		}
	}

	return nil, fmt.Errorf("Job CJAbort not fond :%s", cjid)
}

func (c *Job) cjcomplete(ctx context.Context) error {
	// c.cjlock.Lock()
	// defer c.cjlock.Unlock()
	isall := true
	for i := 0; i < len(c.ControlJobs); i++ {
		v := c.ControlJobs[i]
		isallcmplet := true
		for _, p := range v.processJobs {
			if p.State != "COMPLETED" {
				isallcmplet = false
				break
			}
			// if p.PROCESS_JOB_ID == pj.PROCESS_JOB_ID {
			// 	iscj =true
			// 	break
			// }
		}
		if isallcmplet {
			v.State = "COMPLETED"
			v.Endtime = time.Now().UnixNano()
			err := ControlJobDAOInstance.Update(v)
			if err != nil {
				return err
			}
			// 自动unload ，需要判断所有 物料完成
		} else {
			isall = false
		}
	}
	if isall {
		c.state = "IDLE"
	}

	return nil
}
func (c *Job) CJList(ctx context.Context) ([]*models.ControlJob, error) {

	ret := []*models.ControlJob{}
	for _, v := range c.ControlJobs {

		ret = append(ret, &v.ControlJob)
	}
	return ret, nil
}

func (c *Job) clearAllWafer(ctx context.Context) error {
	for _, v := range c.cjpms {
		err := c.ClearWafer(ctx, v.name)
		if err != nil {
			//log.Error(err.Error())
			return err
		}
	}
	for _, v := range c.cjtms {
		err := c.ClearWafer(ctx, v.name)
		if err != nil {
			//log.Error(err.Error())
			return err
		}
	}
	// 清除全部
	ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), "env.wid.")
	if err != nil {
		//log.Error(err.Error())
		return err
	}
	unit := map[string]bool{}
	for key, uid := range ws {

		err = chenxi.CX.IOServer.WriteString(ctx, key, "")
		if err != nil {
			return err
		}
		uss := strings.Split(uid.(string), ".")
		if len(uss) < 2 {
			return err
		}
		ukey := uss[0]
		for i := 1; i < len(uss)-1; i++ {
			ukey += "."
			ukey += uss[i]
		}
		ukey += ".wid"
		ukey += uss[len(uss)-1]

		err = chenxi.CX.IOServer.WriteString(ctx, ukey, "")
		if err != nil {
			return err
		}
		unit[uss[0]] = true
	}
	for k := range unit {
		err = c.ClearWafer(ctx, k)
		if err != nil {
			return err
		}
	}

	c.ControlJobs = make([]*ControlJob, 0)
	return nil
}

func (c *Job) ClearWafer(ctx context.Context, unit string) error {
	if c.state == "RUN" {
		return fmt.Errorf("ClearWafer 请先Abort")
	}
	if unit == "ALL" {
		return c.clearAllWafer(ctx)
	}
	ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), unit+".wid.")
	if err != nil {
		//log.Error(err.Error())
		return err
	}
	for key, wid := range ws {

		err = chenxi.CX.IOServer.WriteString(ctx, key, "")
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteString(ctx, "env.wid."+wid.(string), "")
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Job) UnitPause(ctx context.Context, id string) error {
	log.Debug("Call Job UnitPause", "id", id)

	if v, ok := c.cjpms[id]; ok {
		return v.pause(ctx, "")

	}

	if v, ok := c.cjtms[id]; ok {
		return v.pause(ctx, "")
	}

	return fmt.Errorf("Job UnitPause not fond :%s", id)
}
func (c *Job) UnitResume(ctx context.Context, id string) error {
	log.Debug("Call Job UnitResume", "id", id)

	if v, ok := c.cjpms[id]; ok {
		return v.resume(ctx, "")

	}

	if v, ok := c.cjtms[id]; ok {
		return v.resume(ctx, "")
	}

	return fmt.Errorf("Job CJResume not fond :%s", id)
}

func (c *Job) UnitAbort(ctx context.Context, id string) error {
	log.Debug("Call Job UnitAbort", "id", id)

	if v, ok := c.cjpms[id]; ok {
		return v.abort(ctx, "")

	}

	if v, ok := c.cjtms[id]; ok {
		return v.abort(ctx, "")
	}

	return fmt.Errorf("Job UnitAbort not fond :%s", id)
}

// CJCreate(ctx context.Context, parm string) error //perm:none
// CJStart(ctx context.Context, parm string) error  //perm:none
// CJStop(ctx context.Context, parm string) error   //perm:none
// CJPause(ctx context.Context, parm string) error  //perm:none
// CJResume(ctx context.Context, parm string) error //perm:none
// CJAbort(ctx context.Context, parm string) error  //perm:none
// CJList(ctx context.Context, parm string) error   //perm:none
var _ api.JobApi = (*Job)(nil)
