package job

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/api"
	"github.com/chenxiio/chenxi/cfg"
	"github.com/chenxiio/chenxi/comm"
)

// src
//lock prepick  preout
// unlock
// pick pick
//

// dist

// preplace prein

type UnitjobTM struct {
	// i *int
	// s string
	name      string
	statename string //
	tp        string // TM,_MV ,如果是_MV类型比较特殊
	data      any
	curaction *TMAction

	TMActions []*TMAction
	unitcfg   *cfg.Module
	unitapi   api.TMApi
	cond      *sync.Cond

	//iswait    bool

	state     string
	lockPause sync.Mutex
}

var defaulttms map[string]api.TMApi = make(map[string]api.TMApi)

func NewUnitjobTM(name string, utype string) (*UnitjobTM, error) {
	// at := strings.Split(apitype, ".")
	// if len(at) != 2 {
	// 	return nil, fmt.Errorf("NewNewUnitjobTM(%s,%s) apitype 格式错误", name, utype)
	// }

	statename := fmt.Sprintf("%s.state", name)
	n, err := chenxi.CX.IOServer.ReadString(context.TODO(), statename)
	if err != nil {
		//log.Error(err.Error())
		n = ""
	}
	// if n != "IDLE" {
	// 	// 执行
	// 	return nil, err
	// }
	var ucfg cfg.Module
	if utype != "_MV" {
		ucfg, err = chenxi.CX.Cfg.Modules.GetCfgByUnit(name)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
	} else {
		ucfg.Slot_count = 1
	}

	ujtm := &UnitjobTM{name: name, tp: utype, data: n, state: "INIT",
		unitcfg: &ucfg, statename: statename, cond: sync.NewCond(&sync.Mutex{})}

	err = chenxi.CX.IOServer.Sub("subio", ujtm.statename, ujtm)

	go ujtm.StartExecution()
	return ujtm, nil
}
func (t *UnitjobTM) init(parm string) error {

	log.Info("Call UnitjobTM init", "Module", t.name)

	if t.curaction != nil {
		cuaction := t.curaction
		t.curaction = nil

		cuaction.moveActions.Lockfrom.TryLock()
		cuaction.moveActions.Lockfrom.Unlock()
		cuaction.moveActions.Lockto.TryLock()
		cuaction.moveActions.Lockto.Unlock()
	}
	t.lockPause.TryLock()
	t.lockPause.Unlock()

	t.state = "INIT"
	err := chenxi.CX.IOServer.SetState(context.TODO(), t.statename, "IDLE")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	if t.tp == "_MV" {
		if _, ok := defaulttms[t.name]; !ok {
			defaulttms[t.name] = &TMTest{Name: t.name}
		}
		t.unitapi = defaulttms[t.name]
		// if err != nil {
		// 	log.Error(err.Error())
		// 	return err
		// }
		// t.unitapi = uapi.(api.TMApi)

	} else {
		uapi, err := chenxi.CX.GetModule(t.name)
		if err != nil {
			log.Error(err.Error())
			return err
		}
		t.unitapi = uapi.(api.TMApi)
	}

	t.TMActions = []*TMAction{}

	return t.unitapi.Init(context.TODO(), parm)
}
func (t *UnitjobTM) Dispatch(data ...any) {
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
func (t *UnitjobTM) dispatching(old any) error {
	t.lockPause.Lock()

	defer t.lockPause.Unlock()
	defer t.lockPause.TryLock()
	switch t.state {
	case "PAUSED":
		//加入队列
	case "ABORTED":
		return nil
	case "IDLE":
		var err error

		switch old {
		case "Init":
			t.state = "IDLE"
		case "Ready":

		case "PROCESS":
			// process 结束
		case "Abort":
		case "Resume":
		case "End":

		case "PrePick":
			// 等待 src unit 完成

			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}

			t.curaction.moveActions.Lockfrom.Lock()
			if t.tp == "_MV" {
				t.curaction.moveActions.Lockto.Lock()
			}
			// 执行pick
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}

			if t.tp == "_MV" {

				err = jobIns.cjpms[t.curaction.form].move()
				if err != nil {
					break
				}

				err = t.postplace()
				if err != nil {
					break
				}
			} else {
				err = t.pick()
				if err != nil {
					break
				}
			}

		case "PrePlace":
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}
			t.curaction.moveActions.Lockto.Lock()
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}
			err = t.place()
			if err != nil {
				return err
			}
		case "Pick":
			//设置 src idle
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}
			srcpm := t.curaction.moveActions.Slots[0].Name
			distpm := 0
			// delete
			k := 0
			// j := 0 // 间隔

			for i := 0; i < t.unitcfg.Slot_count && k < len(t.curaction.moveActions.Slots); i++ {
				v := t.curaction.moveActions.Slots[k]
				if v.Name == i+srcpm {
					err = t.Swap(t.curaction.form, t.curaction.Arm, v.WaferId, v.Name, distpm+i+1)
					if err != nil {
						return err
					}
					k++
				}
				//j++
			}
			if chenxi.CX.Cfg.Modules.TMcfgs[t.name].Type == cfg.CARRIER {
				// 移动fuop

				err := t.Swapcarrier(t.curaction.form, t.name)
				if err != nil {
					return err
				}

			}
			err = t.postpick()
			if err != nil {
				return err
			}
			err = chenxi.CX.IOServer.SetState(context.TODO(), t.curaction.form+".state", "IDLE")
			if err != nil {
				return err
			}

		case "Place":
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}

			err = t.postplace()
			if err != nil {
				return err
			}

			err = chenxi.CX.IOServer.SetState(context.TODO(), t.curaction.to+".state", "IDLE")
			if err != nil {
				return err
			}

		case "PostPick":
			// 转移位置
			//t.unitcfg.
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}
			if t.curaction.times == 0 {
				jobIns.cjpms[t.curaction.to].lock.Lock()
			}
			t.curaction.times++
			err = t.preplace()
			if err != nil {
				return err
			}
			// err = jobIns.cjpms[t.curaction.to].prein(context.Background(), "")
			// if err != nil {
			// 	//setalarm
			// }
		case "PostPlace":
			// bstm := 0
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}
			t.curaction.moveActions.Lockto.Lock()
			t.curaction.moveActions.Lockto.Unlock()
			t.curaction.moveActions.Lockfrom.Lock()
			t.curaction.moveActions.Lockfrom.Unlock()
			if t.curaction == nil {
				log.Error(fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name).Error())
				return nil
			}

			srcpm := t.curaction.moveActions.Slots[0].Name
			distpm, err := jobIns.cjpms[t.curaction.to].LastSlot()

			if err != nil {
				return err
			}

			// if distpm == 0 {
			// 	distpm = 1
			// }
			uj := jobIns.cjpms[t.curaction.to].unitjob
			if CTCIns.cfg.ReturningMode == 0 && uj != nil {
				if len(uj.pj.mainrecipe.Steps) < uj.StepName {
					scan, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("env.wid.scan.%s", t.curaction.moveActions.Slots[0].WaferId))
					if err != nil {
						return err
					}
					uname, uslot, err := comm.UnmarshalUSlot(scan)
					if err != nil {
						return err
					}
					// 如果是起始cm,获取原始位置
					// 或者直接判断 step
					if uname == t.curaction.to {
						distpm = uslot - 1
					}
				}
				//\&&

			}

			// delete
			k := 0
			//j := 0 // 间隔

			if chenxi.CX.Cfg.Modules.TMcfgs[t.name].Type == cfg.CARRIER {
				// 移动fuop
				err := t.Swapcarrier(t.name, t.curaction.to)
				if err != nil {
					return err
				}
			}
			for i := 0; i < t.unitcfg.Slot_count && k < len(t.curaction.moveActions.Slots); i++ {
				v := t.curaction.moveActions.Slots[k]
				if v.Name == i+srcpm {
					if t.tp == "_MV" {
						err = t.Swap(t.curaction.form, t.curaction.to, v.WaferId, v.Name, distpm+i+1)
					} else {
						err = t.Swap(t.curaction.Arm, t.curaction.to, v.WaferId, i+1, distpm+i+1)
					}

					if err != nil {
						return err
					}
					// 删除元素
					// t.curaction.moveActions.Slots = append(t.curaction.moveActions.Slots[:i], t.curaction.moveActions.Slots[i+1:]...)
					// i--
					k++
				}

			}
			// wids, err := chenxi.CX.IOServer.ReadFromPrefix(context.Background(), "env.wid.")
			// if err != nil {
			// 	break loop
			// }
			// wstr, err := json.MarshalIndent(wids, "", "")
			// if err != nil {
			// 	break loop
			// }
			// fmt.Println(string(wstr))

			t.curaction.moveActions.Slots = t.curaction.moveActions.Slots[k:]

			// 判断是否搬运完成
			if len(t.curaction.moveActions.Slots) > 0 {
				// 继续
				err = t.prepick()
				if err != nil {
					return err
				}
			} else {
				// 判断 是否有next

				// 判断uj是否全部搬运完成,如果没有uj 判断pj是否全部完成

				// to 执行 process

				// 没有uj 的单独动作

				if t.curaction.end != t.curaction.to {
					//
					err = jobIns.Move3(context.TODO(), t.curaction.to, t.curaction.end, t.curaction.moveActions.pj.MATERIAL_LIST, t.curaction.moveActions.pj)
					if err != nil {
						return err
					}
				} else if uj != nil {
					// 所有的
					allmvcmplt, err := jobIns.cjpms[t.curaction.to].IsMoveCmplt()

					if err != nil {
						return err
					}
					if allmvcmplt {
						// 所有完成
						if uj.StepName > len(uj.pj.mainrecipe.Steps) {
							// 最后一步，返回 cm
							// err := jobIns.cjpms[t.curaction.to].complete(context.Background(), "")
							// if err != nil {
							// 	// set alrm
							// }
							//jobIns.AddAction(uj)
							//return paths
						} else if uj.StepName == -1 {
							//// return wafer
							// err := jobIns.cjpms[t.curaction.to].complete(context.Background(), "")
							// if err != nil {
							// 	// set alrm
							// }
							if _, ok := chenxi.CX.Cfg.CTCCfg.Return_paths[uj.Unit]; ok {
								jobIns.AddAction(uj)
							}

						} else {
							if jobIns.cjpms[t.curaction.to].unitcfg.Slot_count >= jobIns.cjpms[t.curaction.form].unitcfg.Slot_count*jobIns.cjpms[t.curaction.to].merge &&
								jobIns.cjpms[t.curaction.to].merge > len(jobIns.cjpms[t.curaction.to].unitjob.ujs) {

							} else {
								err = jobIns.cjpms[t.curaction.to].process()
								if err != nil {
									return err
								}
							}
						}
					}

				} else {
					// 完成判断pj是否全部完成,物料是否回到起始位置
					if t.curaction.moveActions.pj == nil {
						// pj ==nil, 没有pj 的单独move动作

					} else {
						err = t.curaction.moveActions.pj.pjcomplete(context.TODO())
						if err != nil {
							return err
						}
					}

				}

				//jobIns.cjpms[t.curaction.form].unitjob.StepName
				// err := jobIns.cjpms[t.curaction.to].process()
				// if err != nil {
				// 	// setalarm
				// }
				tca := t.curaction
				t.curaction = nil
				jobIns.cjpms[tca.to].lock.Unlock()

				// 查找新的moveaction 然后执行

				for _, tm := range jobIns.cjtms {
					tm.cond.L.Lock()

					tm.cond.Signal()
					tm.cond.L.Unlock()
				}
			}
		}

	case "INIT":
		fallthrough
	default:
		if old == "Init" {
			t.state = "IDLE"
		}
		return nil
	}
	return nil
}

func (ut *UnitjobTM) Swapcarrier(src, dist string) error {
	fid, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.carrier.id", src))
	if err != nil {
		return err
	}
	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.carrier.id", dist), fid)
	if err != nil {
		return err
	}
	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.carrier.id", src), "")
	if err != nil {
		return err
	}

	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.carrier.%s", fid), dist)
	if err != nil {
		return err
	}
	punit, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("env.wid.%s", fid))
	if err != nil {
		return nil
	}
	if strings.HasPrefix(punit, src) {
		err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.%s", fid), "")
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *UnitjobTM) Swap(src, dist, wid string, sslot, dslot int) error {

	err := chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.%s", wid), fmt.Sprintf("%s.%d", dist, dslot))
	if err != nil {
		return err
	}
	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.wid.%d", dist, dslot), wid)
	if err != nil {
		return err
	}
	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.wid.%d", src, sslot), "")
	if err != nil {
		return err
	}
	return nil
}

func (t *UnitjobTM) AddAction(act TMAction) {
	t.cond.L.Lock()
	defer t.cond.L.Unlock()
	log.Info("UnitjobTM AddAction ", "tm", t.name, "from", act.form, "to", act.to)
	t.TMActions = append(t.TMActions, &act)
	// src-dst 优先级
	sort.Slice(t.TMActions, func(i, j int) bool {
		return t.TMActions[i].PickPriority > t.TMActions[j].PickPriority
	})

	if len(t.TMActions) == 1 {

		t.cond.Signal()
	}
}

//	func (t *UnitjobTM) RemoveAction(act *TMAction) {
//		t.cond.L.Lock()
//		log.Info("UnitjobTM RemoveAction ", act.form, act.to)
//		// 删除一个元素
//		// 删除一个元素
//		for i, a := range t.TMActions {
//			if a == act {
//				t.TMActions = append(t.TMActions[:i], t.TMActions[i+1:]...)
//				break
//			}
//		}
//		t.cond.L.Unlock()
//		//t.StartExecution()
//	}
var lock sync.Mutex

func (ta *UnitjobTM) StartExecution() {
	tai := 0

	for {
		tai++
		//fmt.Println(ta.name, tai, "循环前")
		ta.cond.L.Lock()
		//fmt.Println(ta.name, tai, "循环后")
		//fmt.Println("UnitjobTM start")

		// for len(ta.TMActions) == 0 || ta.curaction != nil {
		// 	log.Info("UnitjobTM wait ", "tm", ta.name)
		// 	ta.cond.Wait()
		// }

		log.Info("UnitjobTM wait ", "tm", ta.name)
		ta.cond.Wait()

		if ta.curaction != nil {
			ta.cond.L.Unlock()
			continue
		}
		// 执行tmactions
		var names []string
		for _, a := range ta.TMActions {
			names = append(names, a.Arm)
		}
		log.Info("UnitjobTM 查找   ", "list", names)
		//fmt.Println("UnitjobTM A开始执行tmactions", strings.Join(names, ","))
		// 查找并执行
		lock.Lock()

		for i := 0; i < len(ta.TMActions); i++ {
			v := ta.TMActions[i]

			if v.moveActions.Lock.TryLock() {

				if jobIns.cjpms[v.form].lock.TryLock() {
					//  锁定 Interlocking 后才算完整

					// if !v.moveActions.pj.issubpj {

					// }
					v.moveActions.mvid = v.moveActions.Slots[0].WaferId
					if CTCIns.InterlockingLock(v.moveActions.mvid, v.form, v.to) || v.moveActions.pj.issubpj {
						ta.curaction = v
						ta.curaction.times = 0
						if ta.tp == "_MV" {
							if !jobIns.cjpms[v.to].lock.TryLock() {
								//
								ta.curaction = nil
								CTCIns.InterlockingUnLock(v.moveActions.mvid, v.form, v.to)
								jobIns.cjpms[v.form].lock.Unlock()
								v.moveActions.Lock.Unlock()
								continue
							}

							//jobIns.cjpms[v.form].preout(context.TODO(),ta.curaction)
						}

						err := ta.prepick()
						if err != nil {
							log.Error(err.Error())
							ta.curaction = nil
							jobIns.Abort(context.Background())
						}

						break

					}
					jobIns.cjpms[v.form].lock.Unlock()
				}

				v.moveActions.Lock.Unlock()

			} else {
				// 删除
				log.Info("UnitjobTM delete", "tm", ta.name, "from", v.form, "to", v.to)
				ta.TMActions = append(ta.TMActions[:i], ta.TMActions[i+1:]...)
				i--

			}

		}
		lock.Unlock()
		ta.cond.L.Unlock()
	}

	// setalarm
}

func (t *UnitjobTM) postplace() error {

	log.Info("Call UnitjobTM postplace", "Module", t.name)
	return t.unitapi.PostPlace(context.TODO(), t.curaction.curcmd)

}
func (t *UnitjobTM) place() error {

	log.Info("Call UnitjobTM place", "Module", t.name)
	// curcmd 改成 dist
	return t.unitapi.Place(context.TODO(), t.curaction.curcmd)

}
func (t *UnitjobTM) preplace() error {

	log.Info("Call UnitjobTM preplace", "Module", t.name)
	if t.curaction == nil {
		log.Error(fmt.Errorf("UnitjobTM preplace  curaction == nil %s", t.name).Error())
		return nil
	}
	t.curaction.moveActions.Lockto.Lock()
	// curcmd 改成 dist

	l, err := jobIns.cjpms[t.curaction.to].LastSlot()
	if err != nil {
		return err
	}
	t.curaction.curcmd = fmt.Sprintf("%s,%s.%d", t.curaction.Arm,
		t.curaction.to, l+1,
	)
	err = t.unitapi.PrePlace(context.TODO(), t.curaction.curcmd)
	if err != nil {
		return err
	}

	return jobIns.cjpms[t.curaction.to].prein(context.Background(), t.curaction)

}
func (t *UnitjobTM) postpick() error {

	log.Info("Call UnitjobTM postpick", "Module", t.name)
	return t.unitapi.PostPick(context.TODO(), t.curaction.curcmd)

}
func (t *UnitjobTM) pick() error {

	log.Info("Call UnitjobTM Pick", "Module", t.name)
	if t.curaction == nil {
		return fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name)
	}
	return t.unitapi.Pick(context.TODO(), t.curaction.curcmd)

}
func (t *UnitjobTM) prepick() error {
	log.Info("Call UnitjobTM PrePick", "Module", t.name)

	if t.curaction != nil {
		// 开始搬运
		//from preout， tm prepick
		t.curaction.moveActions.Lockfrom.Lock()
		if t.curaction == nil {
			return fmt.Errorf("UnitjobTM pick  curaction == nil %s", t.name)
		}
		// 构建搬运指令 ,可分多次搬运
		// 获取unit slots 数量，tm 一次搬运数量
		if len(t.curaction.moveActions.Slots) == 0 {
			return errors.New("UnitjobTM prepick无slot可以搬运")
		}
		t.curaction.curcmd = fmt.Sprintf("%s.%d,%s",
			t.curaction.form, t.curaction.moveActions.Slots[0].Name,
			t.curaction.Arm)
		if t.tp == "_MV" {
			t.curaction.moveActions.Lockto.Lock()
			t.curaction.curcmd = fmt.Sprintf("%s.%d,%s",
				t.curaction.form, t.curaction.moveActions.Slots[0].Name,
				t.curaction.to)

			jobIns.cjpms[t.curaction.to].prein(context.Background(), t.curaction)
		}
		err := t.unitapi.PrePick(context.TODO(), t.curaction.curcmd)
		if err != nil {
			return err
		}

		// 判断from 类型，调用 preout
		return jobIns.cjpms[t.curaction.form].preout(context.Background(), t.curaction)

	} else {
		// 继续等待
		//t.StartExecution()
	}

	// err := chenxi.CX.IOServer.SetState(context.TODO(), t.statename, "PrePick")
	// if err != nil {
	// 	return err
	// }
	// 根据 slot 构建指令
	// t.unitapi.PrePick(context.TODO(), act.)
	return nil

}

func (p *UnitjobTM) pause(ctx context.Context, parm string) error {

	log.Info("Call UnitjobTM pause", "Module", p.name)
	if p.state != "IDLE" {
		//
		return fmt.Errorf("UnitjobControl Pause error p.state:%s", p.state)
	}
	p.lockPause.Lock()
	p.state = "PAUSED"

	return nil
}

func (p *UnitjobTM) abort(ctx context.Context, parm string) error {
	log.Info("Call UnitjobTM abort", "Module", p.name)
	if p.state != "IDLE" && p.state != "PAUSED" {
		//
		return fmt.Errorf("UnitjobControl Pause error p.state:%s", p.state)
	}
	p.state = "ABORTED"

	return nil
}
func (p *UnitjobTM) resume(ctx context.Context, parm string) error {
	log.Info("Call UnitjobTM resume", "Module", p.name)
	if p.state != "PAUSED" {
		//
		return fmt.Errorf("UnitjobControl Pause error p.state:%s", p.state)
	}
	p.state = "IDLE"
	p.lockPause.TryLock()
	p.lockPause.Unlock()
	return nil
}
