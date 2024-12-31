package job

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm"
	"github.com/chenxiio/chenxi/models"
	"github.com/syndtr/goleveldb/leveldb"
)

// Process Job Management Standard (E40)
// PIS: Process Job State Model
// Q/P:QOUEUED/POOLED
// SU:SETTING UP
// WFS:WAITING FOR START
// P:PROCESSING
// PG:PAUSING
// PC:PROCESS COMPLETE
// PD:PAUSED
// STP:STOPPING
// A:ABORTING

// 在 300mm 晶圆的半导体制造中,process job 通常是在 SECS(半导体设备通信标准)协议中定义的。
// SECS 中定义了多种消息来描述和控制 process job,主要有:
// - S2F41 - 过程作业管理(Process Job Management)
// - S2F47 - 过程作业创建(Process Job Create)
// - S2F49 - 过程作业修改(Process Job Modify)
// - S2F51 - 过程作业删除(Process Job Delete)
// - S2F53 - 过程作业查询(Process Job Query)
// 在 S2F47 消息中会定义一个 process job,主要参数包括:
// - PROCESS_JOB_ID - 过程作业ID
// - PROCESS_JOB_NAME - 过程作业名称
// - PROCESS_DEFINITION_ID - 对应的工艺定义ID
// - PARAMETER_LIST - 作业参数列表
// - MATERIAL_LIST - 材料列表
// - PRIORITY - 优先级
// - PROCESS_JOB_NOTES - 注释
// 设备收到 S2F47 后将创建该过程作业,并可以用 S2F49 进行修改。S2F41 用于控制作业的执行,如等待、开始、暂停等。
// 通过 SECS 协议中的标准消息,不同设备可以兼容地运行和控制 process job,实现自动化的半导体生产流程。
type ProcessJob struct {
	// PROCESS_JOB_ID        string
	// PROCESS_JOB_NAME      string
	// PROCESS_DEFINITION_ID string
	// PARAMETER_LIST        map[string]string // recipe list,recipe1,recipe2
	// MATERIAL_LIST         []string          // A1,B2,C3  A1~A17
	// PRIORITY              string
	// PROCESS_JOB_NOTES     string
	// State                 string // wait  run end down Pause(不退料) Abort(退料)
	// Createtime            int64
	// Starttime             int64
	// Endtime               int64
	models.ProcessJob
	allocated  int //已分配数量
	mainrecipe *cfg.ProcessRecipe
	unitjobs   []*Unitjob
	issubpj    bool
}

// type PJManagement struct {
// 	pjlist []*ProcessJob
// }

// func NewPJManagement() *PJManagement {

// 	// wf := make([][]bool, c)
// 	// for i := 0; i < c; i++ {
// 	// 	wf[i] = make([]bool, n)
// 	// }
// 	return &PJManagement{}
// }

//	func (c *ProcessJob) Create(ctx context.Context, MATERIAL_LIST []string) (*Unitjob, error) {
//		return nil, nil
//	}
//
//	func (c *ProcessJob) Modify(ctx context.Context, MATERIAL_LIST []string) (*Unitjob, error) {
//		return nil, nil
//	}

// func (p *ProcessJob) UJCreate(pjid string, matlist []string, unit string, step int) {

// 	unitjob := &Unitjob{PROCESS_JOB_ID: pjid,
// 		MATERIAL_LIST: matlist, Unit: unit, State: "SETTINGUP",
// 		StepName: step, rcp: p.mainrecipe, pj: p}

//		err := UnitjobDAOInstance.Insert(p.unitjob)
//		if err != nil {
//			return err
//		}
//	}

func (c *ProcessJob) UJCreate(ctx context.Context, uj Unitjob) error {
	log.Debug("call ProcessJob UJCreate", c.PROCESS_JOB_ID, uj.Unitjob_id)

	v := c.mainrecipe.Steps[0]
	for _, u := range v.Unit {
		// 判断u是否空闲，
		us := strings.Split(u, "|")
		if len(us) != 2 {
			return fmt.Errorf("配方格式错误：%s", c.PARAMETER_LIST["main"])
		}
		uname := us[0]
		rcpname := us[1]

		err := chenxi.CX.IOServer.SetState(ctx, uname+".rcp", rcpname)
		if err != nil {
			//log.Info("unit is busy")
			continue
		}

		pmcfg, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(uname)
		if err != nil {
			return err
		}
		// 根据 cm tm.slotcount不同创建多个uj，
		mat := c.MATERIAL_LIST[0]
		us2, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+mat)
		if err != nil {
			return err
		}
		u2, _, err := comm.UnmarshalUSlot(us2)
		if err != nil {
			return err
		}
		mv, err := CTCIns.GenActions(u2, uname)
		if err != nil {
			return err
		}
		upcfg, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(mv.Cmds.Actions[0].Name)
		if err != nil {
			return err
		}
		endwafer := 1
		if upcfg.Slot_count < pmcfg.Slot_count {
			endwafer := c.allocated + pmcfg.Slot_count

			if endwafer > len(c.MATERIAL_LIST) {
				endwafer = len(c.MATERIAL_LIST)
			}

		} else {
			for i := 1; i < len(c.MATERIAL_LIST); i++ {
				mat1 := c.MATERIAL_LIST[i]
				us3, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+mat1)
				if err != nil {
					return err
				}
				u3, _, err := comm.UnmarshalUSlot(us3)
				if err != nil {
					return err
				}
				if u3 == u2 {
					endwafer++
				} else {
					break
				}
			}
			endwafer += c.allocated
		}
		matlist := append([]string(nil), c.MATERIAL_LIST[c.allocated:endwafer]...)
		c.allocated = endwafer
		uj.MATERIAL_LIST = matlist
		// matlist := c.MATERIAL_LIST[c.allocated:endwafer]
		// c.allocated = endwafer
		// uj.MATERIAL_LIST = matlist
		// uj := &Unitjob{Unitjob: models.Unitjob{PROCESS_JOB_ID: c.PROCESS_JOB_ID,
		// 	MATERIAL_LIST: matlist, Unit: uname, State: "SETTINGUP", StepName: 1},
		// 	rcp: c.mainrecipe, pj: c}

		// 生成unitjob，并执行
		err = jobIns.cjpms[uname].Startjob(context.TODO(), &uj)
		if err != nil {
			return err
		}
		if endwafer == len(c.MATERIAL_LIST) {
			// 执行下一个processjob
			for k := 0; k < len(jobIns.ControlJobs); k++ {
				cj := jobIns.ControlJobs[k]

				for i := 0; i < len(cj.processJobs); i++ {
					if cj.processJobs[i].PROCESS_JOB_ID == c.PROCESS_JOB_ID {

						for _, v := range cj.processJobs {
							if v.State == "QUEUED" {
								return v.Start(context.TODO(), "")
							}
						}

						for _, v := range jobIns.ControlJobs {
							if v.State == "QUEUED" {
								return v.Start(context.TODO(), "")
							}
						}

						break
					}
				}
			}
			break
		} else {
			//return fmt.Errorf("")
		}
	}

	// 根据空闲unit数量生成unitjob数量
	// 如果没做完继续生成下一个unitjob
	// 如果做完了 执行下一步step

	return nil
}

func (c *ProcessJob) UJContinue(ctx context.Context, uj Unitjob) error {
	log.Debug("call ProcessJob UJContinue", c.PROCESS_JOB_ID, uj.Unitjob_id)

	if len(c.mainrecipe.Steps) < uj.StepName+1 {
		// 如果是子流程 不回起点
		if uj.pj.issubpj {
			// 如果是回
			return nil
		}
		if chenxi.CX.Cfg.CTCCfg.ReturningMode == 0 {
			return jobIns.Move2(context.Background(), uj.Unit, uj.MATERIAL_LIST, 0, c)
		}
		return nil
		// 回起点

	} else if uj.StepName == -1 {
		// return wafer
		if p, ok := chenxi.CX.Cfg.CTCCfg.Return_paths[uj.Unit]; ok {
			err := jobIns.cjpms[p].Startjob(context.TODO(), &uj)
			if err != nil {
				return err
			}
		}
		return nil
	}
	v := &c.mainrecipe.Steps[uj.StepName]
	v1 := &c.mainrecipe.Steps[uj.StepName-1]
	for _, u := range v.Unit {
		// 判断u是否有wafer 或者foup

		// 判断u是否空闲，

		us := strings.Split(u, "|")
		if len(us) != 2 {
			return fmt.Errorf("配方格式错误：%s", c.PARAMETER_LIST["main"])
		}
		uname := us[0]
		rcpname := us[1]
		if uj.Unit == uname {
			uj.StepName += 1
			v.Curmat = uj.MATERIAL_LIST[len(uj.MATERIAL_LIST)-1]
			return c.UJContinue(ctx, uj)
		}
		// 判断 pj 顺序

		// 判断uj顺序
		m1 := -2
		m2 := -1
		if len(v.Curmat) == 0 {
			m1 = -1
			// xin pj 判断 其他pj有没有执行完成，没有的话，返回错误，继续等待
			// 修改队列顺序时注意，只能在合并uj 创建新job完成后才能修改，不然会卡住
			for cjk, cj := range jobIns.ControlJobs {
				for pjk, pj := range cj.processJobs {
					if pj.PROCESS_JOB_ID == c.PROCESS_JOB_ID {
						for i := pjk - 1; i >= 0; i-- {
							pj1 := cj.processJobs[i]
							if pj1.State != "QUEUED" {
								v1 := pj1.mainrecipe.Steps[uj.StepName]
								if v1.Curmat != pj1.MATERIAL_LIST[len(pj1.MATERIAL_LIST)-1] {
									return fmt.Errorf("顺序不对 pj ")
								}
							}
						}
						for i := cjk - 1; i > 0; i-- {
							if jobIns.ControlJobs[i].State != "QUEUED" {
								v1 := jobIns.ControlJobs[i].processJobs[0].mainrecipe.Steps[uj.StepName]
								if v1.Curmat != jobIns.ControlJobs[i].processJobs[0].MATERIAL_LIST[len(jobIns.ControlJobs[i].processJobs[0].MATERIAL_LIST)-1] {
									return fmt.Errorf("顺序不对 pj ")
								}
							}
						}

					}
				}

			}

		}
		for k, m := range c.MATERIAL_LIST {
			if m == v.Curmat {
				m1 = k
				if m2 >= 0 {
					break
				}
			}
			if m == uj.MATERIAL_LIST[0] {
				m2 = k
				if m1 >= -1 {
					break
				}
			}
		}
		m1 += 1
		if m2 != m1 {
			return fmt.Errorf("等待前面物料完成")
		}
		if !CTCIns.InterlockingLock(uj.MATERIAL_LIST[0], uj.Unit, uname) {
			continue
		}
		if jobIns.cjpms[uname].unitcfg.Slot_count >= jobIns.cjpms[uj.Unit].unitcfg.Slot_count*jobIns.cjpms[uname].merge &&
			jobIns.cjpms[uname].unitjob != nil && len(jobIns.cjpms[uname].unitjob.ujs) < jobIns.cjpms[uname].merge {
			fmt.Println(uj)
		} else if len(v1.SubProcess) > 0 && v1.IsSubStart {
			fmt.Println(uj)
		} else {
			err := chenxi.CX.IOServer.SetState(ctx, uname+".rcp", rcpname)
			if err != nil {
				// log.Info("unit is busy")
				continue
			}
		}

		// if jobIns.cjpms[uname].merge > 1 {
		// 	for k, v := range c.unitjobs {
		// 		if v.Unit == uname && v.StepName == uj.StepName+1 && v.State != "COMPLETED" {
		// 			//

		// 		}
		// 	}
		// }

		pmcfg, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(uname)
		if err != nil {
			return err
		}
		if jobIns.cjpms[uj.Unit].merge > 1 &&
			jobIns.cjpms[uj.Unit].unitcfg.Slot_count >= jobIns.cjpms[uname].unitcfg.Slot_count*jobIns.cjpms[uj.Unit].merge {

			count := jobIns.cjpms[uj.Unit].unitcfg.Slot_count / jobIns.cjpms[uj.Unit].merge
			s0 := 1
			uj1 := uj
			uj1.MATERIAL_LIST = []string{}
			ujt := uj
			for i := 0; i < len(uj.MATERIAL_LIST); i++ {
				unitp2, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+uj.MATERIAL_LIST[i])
				if err != nil {
					log.Error(err.Error())
					return err
				}

				_, slotname2, err := comm.UnmarshalUSlot(unitp2)

				if err != nil {
					log.Error(err.Error())
					return err
				}
				if i == 0 {
					s0 = slotname2
				}
				//mats = append(mats, uj.MATERIAL_LIST[i])
				if slotname2-s0+1 >= count {
					ujt.MATERIAL_LIST = append([]string{}, uj.MATERIAL_LIST[:i]...)
					uj1.MATERIAL_LIST = append([]string{}, uj.MATERIAL_LIST[i:]...)
					uj = ujt
					break
				}

			}

			//old := uj
			if len(uj1.MATERIAL_LIST) > 0 {
				go func() {

					jobIns.AddAction(&uj1)
				}()
			}

			// for _, v := range uj.ujs {

			// }
		}
		if pmcfg.Slot_count < len(uj.MATERIAL_LIST) {
			// 拆分
			return fmt.Errorf("pmcfg.Slot_count :%d,实际数量 %d ", pmcfg.Slot_count, len(c.MATERIAL_LIST))
		}

		//matlist := c.MATERIAL_LIST[:]

		if len(v1.SubProcess) > 0 {
			if v1.IsSubStart {

				carrid, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.carrier.id", uname))
				if err != nil {
					return err
				}
				if len(carrid) > 0 {
					v1.IsSubStart = false
					v.Curmat = uj.MATERIAL_LIST[len(uj.MATERIAL_LIST)-1]

					err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.wid.1", uname), "")
					if err != nil {
						return err
					}

					return jobIns.cjpms[uname].Startjob(context.TODO(), &uj)
				} else {
					return fmt.Errorf("载具未到达")
				}

			}
			// 创建 子流程 ，uj.pj
			submats := []string{}

			carrid, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.carrier.id", uj.Unit))
			if err != nil {
				if !strings.HasPrefix(err.Error(), leveldb.ErrNotFound.Error()) {
					return err
				} else {
					// 当前没有 carrierid ，查找并创建sub
					// 出站时，先 start subprocess
					carrid, err = chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("env.wid.carrier.%s", uj.MATERIAL_LIST[0]))
					if err != nil {
						return err
					}
					submats = append(submats, carrid)
					v1.IsSubStart = true
					err = chenxi.CX.IOServer.SetState(ctx, uname+".rcp", "IDLE")
					if err != nil {
						v1.IsSubStart = false
						return err
					}
					err = c.subpj(ctx, submats, v1.SubProcess)
					if err != nil {
						v1.IsSubStart = false
						return err
					}
					go func() {
						jobIns.AddAction(&uj)
					}()
					return nil

				}
			} else {
				// 生成unitjob，并执行
				// 进入时，先startjob
				v.Curmat = uj.MATERIAL_LIST[len(uj.MATERIAL_LIST)-1]
				err = jobIns.cjpms[uname].Startjob(context.TODO(), &uj)
				if err != nil {
					return err
				}
				// 当前有carrierid
				err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.%s", carrid), fmt.Sprintf("%s.%d", uj.Unit, 1))
				if err != nil {
					return err
				}
				submats = append(submats, carrid)
				// slot, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("env.wid.%s", ma))
				// if err != nil {
				// 	return err
				// }
				// _, dslot, err := comm.UnmarshalUSlot(slot)
				// if err != nil {
				// 	return err
				// }
				return c.subpj(ctx, submats, v1.SubProcess)
			}
			// err = chenxi.CX.IOServer.SetState(ctx, uname+".rcp2", v.SubProcess)
			// if err != nil {
			// 	return err
			// }

		} else {
			v.Curmat = uj.MATERIAL_LIST[len(uj.MATERIAL_LIST)-1]
			// 生成unitjob，并执行
			err = jobIns.cjpms[uname].Startjob(context.TODO(), &uj)
			if err != nil {
				return err
			}
		}

		return nil

	}

	// 根据空闲unit数量生成unitjob数量
	// 如果没做完继续生成下一个unitjob
	// 如果做完了 执行下一步step

	return fmt.Errorf("all unit is busy  processid:%s step:%d", c.PROCESS_JOB_ID, uj.StepName+1)
}

func (c *ProcessJob) subpj(ctx context.Context, submats []string, rcp string) error {

	subpj := &ProcessJob{ProcessJob: models.ProcessJob{
		PROCESS_JOB_ID:        c.PROCESS_JOB_ID,
		PROCESS_JOB_NAME:      "sub pj",
		PROCESS_DEFINITION_ID: "sub pj",
		PARAMETER_LIST:        map[string]string{},
		MATERIAL_LIST:         submats,
		State:                 "QUEUED"},
		issubpj: true,
	}
	subpj.PARAMETER_LIST["main"] = rcp
	//uj.subpj = append(uj.subpj, subpj)
	return subpj.Start(ctx, "")

}
func (c *ProcessJob) Start(ctx context.Context, parm string) error {
	log.Debug("call ProcessJob Start", c.PROCESS_JOB_ID, parm)

	// cm := map[string]bool{}
	// for _, v := range c.MATERIAL_LIST {

	// 	unitp, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+v)
	// 	if err != nil {
	// 		return fmt.Errorf("map 不正确，请重新map,%s,%s", unitp, err.Error())
	// 	}
	// 	us := strings.Split(unitp, ".")
	// 	if len(us) < 2 {
	// 		return fmt.Errorf("map 不正确，请重新map,%s,%s", unitp, err.Error())
	// 	}
	// 	if _, ok := cm[us[0]]; ok {
	// 		continue
	// 	}
	// 	load, err := chenxi.CX.IOServer.ReadInt(context.TODO(), us[0]+".load")
	// 	if err != nil || load != 1 {
	// 		// 重新load

	// 		return fmt.Errorf("Load 不正确，请重新Load,%s,%s", unitp, err.Error())
	// 	}

	// 	cm[us[0]] = true
	// }

	rcp, err := chenxi.CX.Recipe.ReadProcessRecipe(context.TODO(), c.PARAMETER_LIST["main"])
	if err != nil {
		return err
	}
	if !c.issubpj {
		drcp, _ := chenxi.CX.Recipe.ReadProcessRecipe(context.TODO(), "in")

		drcp.Steps = append(drcp.Steps, rcp.Steps...)

		outrcp, _ := chenxi.CX.Recipe.ReadProcessRecipe(context.TODO(), "out")

		drcp.Steps = append(drcp.Steps, outrcp.Steps...)
		c.mainrecipe = &drcp
	} else {
		c.mainrecipe = &rcp
	}

	c.State = "EXECUTING"
	c.Starttime = time.Now().UnixNano()
	if !c.issubpj {
		err = ProcessJobDAOInstance.Update(c)
		if err != nil {
			return err
		}
	}

	go jobIns.AddAction(&Unitjob{Unitjob: models.Unitjob{PROCESS_JOB_ID: c.PROCESS_JOB_ID, StepName: 0}, pj: c})

	return nil
}

// 停止分配，将已经分配的做完
func (c *ProcessJob) Stop(ctx context.Context, parm string) error {
	log.Debug("call ProcessJob Stop", c.PROCESS_JOB_ID, parm)
	c.State = "STOPED"

	err := ProcessJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	// end
	return nil
}

// 暂停相关单元
func (c *ProcessJob) Pause(ctx context.Context, parm string) error {
	log.Debug("call ProcessJob Pause", c.PROCESS_JOB_ID, parm)

	if c.State != "EXECUTING" {
		//
		return fmt.Errorf("ProcessJob Pause error C.state:%s", c.State)

	}
	c.State = "PAUSED"
	for _, v := range c.unitjobs {

		if pm, ok := jobIns.cjpms[v.Unit]; ok {
			if pm.state == "IDLE" {
				err := pm.pause(ctx, parm)
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}

	err := ProcessJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	// end
	return nil
}
func (c *ProcessJob) Resume(ctx context.Context, parm string) error {
	//wait
	log.Debug("call ProcessJob Resume", c.PROCESS_JOB_ID, parm)
	if c.State != "PAUSED" {
		//
		return fmt.Errorf("ProcessJob Pause error C.state:%s", c.State)
	}
	c.State = "EXECUTING"
	for _, v := range c.unitjobs {
		if v.State == "PAUSED" {

			if pm, ok := jobIns.cjpms[v.Unit]; ok {
				err := pm.resume(ctx, parm)
				if err != nil {
					log.Error(err.Error())
				}
			}
		}

	}

	err := ProcessJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	return nil
}
func (c *ProcessJob) Abort(ctx context.Context, parm string) error {
	log.Debug("call ProcessJob Abort", c.PROCESS_JOB_ID, parm)

	if c.State != "EXECUTING" && c.State != "PAUSED" {
		//
		return fmt.Errorf("ProcessJob Abort error C.state:%s", c.State)
	}
	c.State = "ABORTED"
	for _, v := range c.unitjobs {

		if pm, ok := jobIns.cjpms[v.Unit]; ok {
			if pm.state == "IDLE" || pm.state == "PAUSED" {
				err := pm.abort(ctx, parm)
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}

	err := ProcessJobDAOInstance.Update(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ProcessJob) pjcomplete(ctx context.Context) error {
	isallcmplet := true
	cmpcount := 0
	for _, v := range c.unitjobs {
		if v.State != "COMPLETED" {
			isallcmplet = false
			return nil
		}
		cmpcount += len(v.MATERIAL_LIST)
	}
	//
	switch CTCIns.cfg.ReturningMode {
	case 0:
		for _, str := range c.MATERIAL_LIST {

			unitp, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid.scan."+str)
			if err != nil {
				isallcmplet = false
				break
			}
			unitp2, err := chenxi.CX.IOServer.ReadString(context.TODO(), "env.wid."+str)
			if err != nil {
				log.Error(err.Error())
				return err
			}

			if unitp != unitp2 {
				isallcmplet = false
				break
			}

		}

	case 1:
		// ruturn path 最终节点
	case 2:
		//
	}

	if isallcmplet {
		c.State = "COMPLETED"
		c.Endtime = time.Now().UnixNano()
		err := ProcessJobDAOInstance.Update(c)
		if err != nil {
			return err
		}
		for _, str := range c.MATERIAL_LIST {

			err := chenxi.CX.IOServer.WriteString(context.TODO(), "env.wid.scan."+str, "")
			if err != nil {
				log.Error(err.Error())
				continue
			}

		}

		return jobIns.cjcomplete(ctx)
	}
	return nil
}

// func (c *ProcessJob) UJEnd(ctx context.Context, uj *Unitjob) error {
// 	log.Debug("call ProcessJob UJEnd", c.PROCESS_JOB_ID, uj.Unitjob_id)
// 	//

// 	return nil

// }
