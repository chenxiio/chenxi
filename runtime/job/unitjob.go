package job

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm"
	"github.com/chenxiio/chenxi/models"
	"github.com/syndtr/goleveldb/leveldb"
)

type Unitjob struct {
	// Unitjob_id     int64
	// PROCESS_JOB_ID string
	// MATERIAL_LIST  []string
	// PreUnit        string // 上一个单元
	// Unit           string // place 成功之后 设置
	// StepName       int    //
	// //Recipe         string // 通过unit 和processrecipe 获得 unitrecipe ，通过 processrecipe 获取单元路径
	// State      string
	// Createtime int64
	// Starttime  int64
	// Endtime    int64
	models.Unitjob
	pj  *ProcessJob
	rcp *cfg.ProcessRecipe
	ujs []*Unitjob // 上一步uj，合并时有多个uj
	//excaction *MoveActions
	// src 单元+ 可用的tm 路径， tm抢占lock保证只执行一次
	//moveActions []*MoveActions
	//srcCarrier cfg.Carrier //最初位置
	//carriers  map[string]api.CMApi
	subpj []*ProcessJob
}

// func NewUnitjob(pjid string, matlist []string, unit string) (*Unitjob, error) {

// 	return &Unitjob{PROCESS_JOB_ID: pjid, MATERIAL_LIST: matlist, Unit: unit, State: "Run"}, nil
// }

type UnitjobControl struct {
	// i *int
	// s string
	name      string
	statename string //

	api       string // HP ,CP
	data      any
	unitjob   *Unitjob
	unitcfg   *cfg.Module
	unitapi   any
	lock      sync.Mutex  // 执行锁，获得锁，才能调用运送接口执行，out 时 释放锁， tm postplace时 释放
	innerLock *sync.Mutex // 群组锁
	lockPause sync.Mutex
	tmaction  *TMAction
	state     string //INIT IDLE RUN PAUSED ABORTED
	merge     int    // 是否合并 1：合并
	//tmslotcount int    // 最多合并carrier次数计算
	//readycount int
}

func (t *UnitjobControl) LastSlot() (int, error) {
	ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), t.name+".wid.")
	if err != nil {
		//log.Error(err.Error())
		return 0, err
	}
	if len(ws) == 0 {
		return 0, nil
	}
	slot := 0
	//keys := make([]int, 0, len(ws))
	for key := range ws {
		_, slot1, err := comm.UnmarshalUSlot(key)
		if err != nil {
			return 0, err
		}
		if slot1 > slot {
			slot = slot1
		}
	}
	//sort.Strings(keys)
	if t.merge > 1 {
		c2 := t.unitcfg.Slot_count / t.merge
		slot = (slot/c2 + 1) * c2
	}
	return slot, nil
}

func (t *UnitjobControl) IsMoveCmplt() (bool, error) {
	// ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), t.name+".wid.")
	// if err != nil {
	// 	//log.Error(err.Error())
	// 	return false, err
	// }
	isallcmplet := true
	for _, str := range t.unitjob.MATERIAL_LIST {

		unitp, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+str)
		if err != nil {
			log.Error(err.Error())
			return false, err
		}

		uname, _, err := comm.UnmarshalUSlot(unitp)

		if uname != t.name {
			isallcmplet = false
			break
		}
	}
	return isallcmplet, nil

}

func (t *UnitjobControl) IsActRemoveCmplt() bool {
	// ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), t.name+".wid.")
	// if err != nil {
	// 	//log.Error(err.Error())
	// 	return false, err
	// }

	for _, s := range t.tmaction.moveActions.Slots {

		_, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.wid.%d", t.name, s.Name))
		if err == nil {
			//log.Error(err.Error())
			return false
		}

		// uname, _, err := comm.UnmarshalUSlot(unitp)

		// if uname == t.name {
		// 	isallcmplet = false
		// 	break
		// }
	}
	return true

}

func (t *UnitjobControl) IsRemoveCmplt() bool {
	// ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), t.name+".wid.")
	// if err != nil {
	// 	//log.Error(err.Error())
	// 	return false, err
	// }

	if t.unitjob == nil {
		for _, s := range t.tmaction.moveActions.Slots {

			_, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.wid.%d", t.name, s.Name))
			if err == nil {
				//log.Error(err.Error())
				return false
			}

			// uname, _, err := comm.UnmarshalUSlot(unitp)

			// if uname == t.name {
			// 	isallcmplet = false
			// 	break
			// }
		}
		return true
	}
	isallcmplet := true

	ups, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), t.name+".wid")
	if err != nil {
		log.Error(err.Error())
		return false
	}
	if len(ups) > 0 {
		for _, str := range t.unitjob.MATERIAL_LIST {

			for _, v := range ups {
				if str == v {
					return false
				}

			}
		}
	}
	return isallcmplet

}

// func (t *UnitjobControl) GetCount() (int, error) {
// 	ws, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), t.name+".wid.")
// 	if err != nil {
// 		//log.Error(err.Error())
// 		return 0, err
// 	}

// 	// for _, wid := range ws {
// 	// 	// if len(wid.(string)) > 0 {
// 	// 	// 	return 0, fmt.Errorf("%s 有wafer 不能初始化", v.Name)
// 	// 	// }

//		// }
//		return len(ws), nil
//	}
func NewUnitjobControl(name string, utype string) (*UnitjobControl, error) {
	// at := strings.Split(apitype, ".")
	// if len(at) != 2 {
	// 	return nil, fmt.Errorf("NewUnitjobControl(%s,%s) apitype 格式错误", name, apitype)
	// }

	statename := fmt.Sprintf("%s.state", name)
	n, err := chenxi.CX.IOServer.ReadString(context.TODO(), statename)
	if err != nil {
		//log.Error(err.Error())
		n = ""
	}

	ucfg, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(name)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	pm := &UnitjobControl{name: name, api: utype, data: n,
		unitcfg: &ucfg, statename: statename, state: "INIT", merge: 1,
	}

	err = chenxi.CX.IOServer.Sub("subio", pm.statename, pm)

	return pm, err
}
func (t *UnitjobControl) init(parm string) error {

	log.Info("Call UnitJob init", "Module", t.name)
	t.state = "INIT"
	err := chenxi.CX.IOServer.SetState(context.TODO(), t.statename, "IDLE")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	err = chenxi.CX.IOServer.SetState(context.TODO(), t.name+".rcp", "IDLE")
	if err != nil {
		return err
	}
	uapi, err := chenxi.CX.GetModule(t.name)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	t.unitapi = uapi

	if t.unitjob != nil {

		t.unitjob = nil

	}
	if t.lock.TryLock() {

	}
	t.lock.Unlock()

	t.lockPause.TryLock()
	t.lockPause.Unlock()
	t.tmaction = nil

	return t.unitapi.(api.ModuleApi).Init(context.TODO(), parm)
}
func (t *UnitjobControl) Dispatch(data ...any) {
	for key, value := range data[0].(map[string]any) {
		fmt.Printf("%s: %v\n", key, value)

		if value == "IDLE" {
			go func(d any) {
				err := t.dispatching(d)
				if err != nil {
					log.Error(err.Error())
					err = jobIns.Abort(context.TODO())
					if err != nil {
						log.Error(err.Error())
					}
				}
			}(t.data)

		}

		t.data = value
	}

}
func (t *UnitjobControl) dispatching(old any) error {
	t.lockPause.Lock()
	defer t.lockPause.Unlock()
	switch t.state {
	case "PAUSED":
		//加入队列 阻塞队列

	case "ABORTED":
		return nil
	case "IDLE":

		var err error

		switch old {
		case "Init":
			log.Info("2 set uj state IDLE", "uj", t.name)
			t.state = "IDLE"
		case "PreIn":
			if t.tmaction == nil {
				log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
				return nil
			}
			t.tmaction.moveActions.Lockto.TryLock()
			t.tmaction.moveActions.Lockto.Unlock()

			if jobIns.cjtms[t.tmaction.Name].tp != "_MV" {
			}
			log.Info("Call UnitJob Place", "Module", t.name)

			err = chenxi.CX.IOServer.SetState(context.TODO(), t.statename, "Place")

			if err != nil {
				break
			}

		case "Place":

			err = t.in()
			if err != nil {
				return err
			}
			// if err != nil {
			// 	// alarm
			// }
		case "In":
			//
			if t.tmaction == nil {
				log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
				return nil
			}
			t.tmaction.moveActions.Lockto.TryLock()
			t.tmaction.moveActions.Lockto.Unlock()

			CTCIns.InterlockingUnLock(t.tmaction.moveActions.mvid, t.tmaction.form, t.tmaction.to)

		case "PreOut":
			if t.tmaction == nil {
				log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
				return nil
			}
			// TM继续执行
			t.tmaction.moveActions.Lockfrom.TryLock()
			t.tmaction.moveActions.Lockfrom.Unlock()

			if jobIns.cjtms[t.tmaction.Name].tp != "_MV" {
				log.Info("Call UnitJob Pick", "Module", t.name)

				err := chenxi.CX.IOServer.SetState(context.TODO(), t.statename, "Pick")

				if err != nil {
					return err
				}
			}

		case "Pick":
			// if t.unitjob == nil {
			// 	log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
			// }
			err := t.out()

			if err != nil {
				return err
			}

		case "Out":
			// if t.unitjob == nil {
			// 	log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
			// }
			if t.IsActRemoveCmplt() {
				if t.IsRemoveCmplt() {

					if t.unitjob != nil {

						err := t.complete(context.Background(), "")
						if err != nil {
							return err
						}
						t.unitjob = nil
						// if t.tp=="PM" {

						// }
						// 如果 当前 unit.carrier.id  != ""  不清除 rcp

					}
					carrid, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.carrier.id", t.name))
					if err != nil {
						if !strings.HasPrefix(err.Error(), leveldb.ErrNotFound.Error()) {
							return err
						}
					}
					if len(carrid) == 0 {
						err = chenxi.CX.IOServer.SetState(context.TODO(), t.name+".rcp", "IDLE")
						if err != nil {
							return err
						}
					}

					if t.api == "PM" {
						jobIns.cond.L.Lock()
						jobIns.cond.Signal()
						jobIns.cond.L.Unlock()
					}
				}

				t.lock.Unlock()

				for _, tm := range jobIns.cjtms {
					tm.cond.L.Lock()
					tm.cond.Signal()
					tm.cond.L.Unlock()
				}
			}

			t.tmaction.moveActions.Lockfrom.Unlock()

		case "Move":
			err := t.out()
			if err != nil {
				return err
			}
			err = chenxi.CX.IOServer.SetState(context.TODO(), t.tmaction.to+".state", "IDLE")
			if err != nil {
				return err
			}

		case "Ready":
			if t.unitjob == nil {
				log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
			}
			// rcp 准备完成之后，调用 前一个单元 preout
			err := jobIns.Move2(context.Background(), t.name, t.unitjob.MATERIAL_LIST, t.unitjob.StepName, t.unitjob.pj)
			if err != nil {
				return err
			}
		case "Process":
			if t.unitjob == nil {
				log.Error("UnitjobControl   unitjob == nil ", "unit", t.name)
			}
			// process 结束
			// t.unitjob.Endtime = time.Now().UnixNano()
			// t.unitjob.State = "COMPLETED"
			// err := UnitjobDAOInstance.Update(t.unitjob)
			// if err != nil {
			// 	// set alrm
			// }
			// 回到原位置才算结束

			// err := t.complete(context.Background(), "")
			// if err != nil {
			// 	// set alrm
			// }
			// 判断下一步

			jobIns.AddAction(t.unitjob)
			// // src 执行
			// if pm, ok := jobIns.cjpms[t.name]; ok {

			// 	// if err != nil {
			// 	// 	//setalrm
			// 	// }
			// 	jobIns.AddAction(pm.unitjob)
			// }

		case "Abort":
		case "Resume":
		case "End":
		case "Load":
		case "Unload":
			// for _, v := range v.CARRIER_ID {
			// 	cm, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("env.carrier.%s", v))
			// 	if err != nil {
			// 		return err
			// 	}
			// 	err = c.cjpms[cm].unload(ctx, "")
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			err := t.unload(context.Background(), "")
			if err != nil {
				return err
			}
			// c.ControlJobs = append(c.ControlJobs[i:], c.ControlJobs[i+1:]...)
			// i--
		default:

		}

	case "INIT":
		fallthrough
	default:
		if old == "Init" {
			log.Info("1 set uj state IDLE", "uj", t.name)
			t.state = "IDLE"
		}
		return nil
	}
	return nil
}

// wafer数量不能超过设备单元能力
// unitjob做完才能进行下一个 job设置
func (p *UnitjobControl) Startjob(ctx context.Context, uj *Unitjob) error {

	// //actions :=
	// if p.merge > 1 {
	// 	// 查找
	// 	for k, v := range pj.unitjobs {
	// 		if v.State != "COMPLETED" && v.Unit == unit && v.StepName == step && p.merge < len(v.ujs) {
	// 			// 是否能够合并
	// 			// 按照总容量计算

	// 		}
	// 	}
	// }

	if p.unitjob != nil && len(p.unitjob.ujs) < p.merge {
		p.unitjob.MATERIAL_LIST = append(p.unitjob.MATERIAL_LIST, uj.MATERIAL_LIST...)
		p.unitjob.ujs = append(p.unitjob.ujs, uj)
		err := UnitjobDAOInstance.Update(p.unitjob)
		if err != nil {
			return err
		}
	} else {
		p.unitjob = &Unitjob{Unitjob: models.Unitjob{PROCESS_JOB_ID: uj.PROCESS_JOB_ID,
			MATERIAL_LIST: uj.MATERIAL_LIST, Unit: p.name, State: "SETTINGUP", StepName: uj.StepName + 1},
			rcp: uj.pj.mainrecipe, pj: uj.pj}

		err := UnitjobDAOInstance.Insert(p.unitjob)
		if err != nil {
			return err
		}
		p.unitjob.ujs = append(p.unitjob.ujs, uj)
		uj.pj.unitjobs = append(uj.pj.unitjobs, p.unitjob)
	}

	if p.unitjob.StepName > len(p.unitjob.pj.mainrecipe.Steps) || p.unitjob.StepName == -1 || len(p.unitjob.ujs) > 1 {

		return jobIns.Move2(context.Background(), p.name, uj.MATERIAL_LIST, p.unitjob.StepName, uj.pj)
	}
	return p.ready(ctx)
}

func (p *UnitjobControl) ready(ctx context.Context) error {
	log.Info("Call UnitJob Ready", "Module", p.name)
	// err := chenxi.CX.IOServer.SetState(ctx, p.statename, "Ready")
	// if err != nil {
	// 	return err
	// }
	// ready 下载unitrcp
	//chenxi.CX.Cfg.Modules.g
	urcp, err := p.unitjob.rcp.GetUnitRcp(p.unitjob.StepName, p.unitjob.Unit, p.unitcfg.Type)
	if err != nil {
		return err
	}
	urcpstr, err := json.Marshal(urcp)
	if err != nil {
		return err
	}
	err = p.unitapi.(api.PMApi).Ready(context.TODO(), string(urcpstr))
	if err != nil {

		return err
	}
	// moveActions 每个move 对应一个src 单元到本单元所以可以用的tm
	// 每个src单元单独执行设置一个 lock 分发给所有的tm，依次分发 防止 同一个目标多个tm同时操作
	return nil
}

func (p *UnitjobControl) prein(ctx context.Context, act *TMAction) error {
	// if p.tmaction != act {
	// 	return fmt.Errorf(" UnitjobControl prein tmaction 与之前不一致 ")
	// }
	p.lockPause.Lock()
	defer p.lockPause.Unlock()
	p.tmaction = act
	log.Info("Call UnitJob prein", "Module", p.name)
	return p.unitapi.(api.ModuleApi).PreIn(ctx, act.curcmd)

}

func (p *UnitjobControl) in() error {
	log.Info("Call UnitJob in", "Module", p.name)
	return p.unitapi.(api.ModuleApi).In(context.TODO(), p.tmaction.curcmd)

}
func (p *UnitjobControl) out() error {
	log.Info("Call UnitJob out", "Module", p.name)
	err := p.unitapi.(api.ModuleApi).Out(context.TODO(), p.tmaction.curcmd)
	if err != nil {
		return err
	}
	return nil
	//
}
func (p *UnitjobControl) move() error {
	log.Info("Call UnitJob move", "Module", p.name)
	err := p.unitapi.(api.ModuleApi).Move(context.TODO(), p.tmaction.curcmd)
	if err != nil {
		return err
	}
	return nil
	//
}
func (p *UnitjobControl) preout(ctx context.Context, act *TMAction) error {
	p.lockPause.Lock()
	defer p.lockPause.Unlock()
	log.Info("Call UnitJob preout", "Module", p.name)
	p.tmaction = act
	return p.unitapi.(api.ModuleApi).PreOut(ctx, act.curcmd)
}
func (p *UnitjobControl) process() error {
	log.Info("Call UnitJob process", "Module", p.name)
	p.unitjob.State = "PROCESS"
	p.unitjob.Starttime = time.Now().UnixNano()
	err := UnitjobDAOInstance.Update(p.unitjob)
	if err != nil {
		return err
	}
	return p.unitapi.(api.PMApi).Process(context.TODO(), "")
}

func (p *UnitjobControl) pause(ctx context.Context, parm string) error {

	log.Info("Call UnitJob pause", "Module", p.name)
	if p.state != "IDLE" {
		//
		return fmt.Errorf("UnitjobControl Pause error p.state:%s", p.state)
	}
	p.lockPause.Lock()
	p.state = "PAUSED"

	if p.unitjob != nil {
		p.unitjob.State = "PAUSED"
		err := UnitjobDAOInstance.Update(p.unitjob)
		if err != nil {
			return err
		}
	}

	return p.unitapi.(api.PMApi).Pause(ctx, parm)
}

func (p *UnitjobControl) abort(ctx context.Context, parm string) error {
	log.Info("Call UnitJob abort", "Module", p.name)
	if p.state != "IDLE" && p.state != "PAUSED" {
		//
		return fmt.Errorf("UnitjobControl Pause error p.state:%s", p.state)
	}
	p.state = "ABORTED"
	if p.unitjob != nil {
		p.unitjob.State = "ABORTED"
		err := UnitjobDAOInstance.Update(p.unitjob)
		if err != nil {
			return err
		}
	}

	return p.unitapi.(api.PMApi).Abort(ctx, parm)
}
func (p *UnitjobControl) resume(ctx context.Context, parm string) error {
	log.Info("Call UnitJob resume", "Module", p.name)
	if p.state != "PAUSED" {
		//
		return fmt.Errorf("UnitjobControl Pause error p.state:%s", p.state)
	}
	p.state = "IDLE"
	p.lockPause.TryLock()
	p.lockPause.Unlock()
	if p.unitjob != nil {
		p.unitjob.State = "RESUMED"
		err := UnitjobDAOInstance.Update(p.unitjob)
		if err != nil {
			return err
		}
	}

	return p.unitapi.(api.PMApi).Resume(ctx, parm)
}

func (p *UnitjobControl) complete(ctx context.Context, parm string) error {
	log.Info("Call UnitJob complete", "Module", p.name)

	p.unitjob.State = "COMPLETED"
	p.unitjob.Endtime = time.Now().UnixNano()
	err := UnitjobDAOInstance.Update(p.unitjob)
	if err != nil {
		return err
	}
	return nil
}

func (p *UnitjobControl) loadandmap(ctx context.Context, parm string) error {
	// if p.tmaction != act {
	// 	return fmt.Errorf(" UnitjobControl prein tmaction 与之前不一致 ")
	// }

	log.Info("Call UnitJob loadandmap", "Module", p.name)
	err := p.unitapi.(api.CMApi).Load(ctx, parm)
	if err != nil {
		return err
	}

	return p.unitapi.(api.CMApi).Map(ctx, parm)
}

func (p *UnitjobControl) unload(ctx context.Context, parm string) error {
	// if p.tmaction != act {
	// 	return fmt.Errorf(" UnitjobControl prein tmaction 与之前不一致 ")
	// }

	log.Info("Call UnitJob unload", "Module", p.name)
	err := p.unitapi.(api.CMApi).Unload(ctx, parm)
	if err != nil {
		return err
	}
	// 设置 wid =""
	wids, err := chenxi.CX.IOServer.ReadFromPrefix(ctx, p.name+".wid.")
	if err != nil {
		return err
	}
	for k, v := range wids {
		//Map(p.Name, fmt.Sprintf("%s_%s_%d", p.Name, cid, i+1), i+1)
		err := chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.%s", v), "")
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.scan.%s", v), "")
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteString(context.TODO(), k, "")
		if err != nil {
			return err
		}
	}

	return err
}
