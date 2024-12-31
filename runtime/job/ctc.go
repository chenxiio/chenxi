package job

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/cfg"
)

// 搬运指令 例如 PM1-PM2
// 生成所有可达指令(依次执行 group 调用tm的 CvtIns 转换具体指令 调用 （enable disable） 得到最终结果 )
// 根据busy ，优先级，选择搬运单元执行

// type MoveAction struct {
// 	id   string

// }

type TMAction struct {
	Name string // tmname
	// unitjob 开始执行 ,选择form单元后，在按tm 优先级先后顺序排序， 发送给tm
	Priority int
	// tm收到所有待执行指令按照 Pick排序然后执行
	PickPriority  int
	PlacePriority int
	// Pick          string // pick指令 unit-tm.arm
	// Place         string // place指令 tm.arm-unit.s
	curcmd   string
	form, to string
	Arm      string // 如果 len(Arm)==0 ，子路径 必须设置unitpathmoveActions
	end      string
	times    int
	//unitpathmoveActions *[]MoveActions
	// tm 队列使用 =
	moveActions *MoveActions
}

type MoveActionCmds struct {
	From string // from - to
	To   string
	// unitjob 开始执行 ,按照所有的form 单元 先后顺序排序，排序后依次发送给tm
	Priority int

	Actions []TMAction
}
type MoveActions struct {
	// 按照优先级排序后
	Cmds MoveActionCmds

	Slots    cfg.Slots // wafer 位置
	pj       *ProcessJob
	mvid     string
	Lock     *sync.Mutex // 获得锁才能执行
	Lockfrom *sync.Mutex // 等待from idel 锁
	Lockto   *sync.Mutex // 等待to idel 锁
}

// func (m *TMAction) isFromMoveCmplt() bool {
// 	if len(m.moveActions.Slots) == 0 {
// 		return true
// 	}
// 	for _, v := range m.moveActions.Slots {

// 		_, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.wid.%d", m.form, v.Name))
// 		if err == nil {
// 			return false
// 		}

// 	}

// 	return true
// }

// func (tma *MoveActionCmds) SortByPriority() {
// 	sort.Slice(tma.Actions, func(i, j int) bool {
// 		return tma.Actions[i].Priority > tma.Actions[j].Priority
// 	})
// }

type CTC struct {
	cfg *cfg.CTCCfg
	// from to 所有可用TM路径
	Actions map[string]MoveActionCmds
	lock    sync.Mutex
	counter int64
}

func (c *CTC) getUniqueNumber() int64 {

	return atomic.AddInt64(&c.counter, 1)

}

var CTCIns *CTC
var ctconce sync.Once

func GetCTCInstance(cfg *cfg.CTCCfg) (*CTC, error) {
	ctconce.Do(func() {
		// for _, v := range cfg.Interlocking {
		// 	//v.Pall = make(map[string]cfg.TAction)
		// }

		CTCIns = &CTC{cfg: cfg, Actions: map[string]MoveActionCmds{}}

	})

	return CTCIns, nil
}

func (c *CTC) Init(parm string) error {
	return nil
}

// func (u *CTC) InterlockingLock(form, to string) bool {

// 	for _, v := range u.cfg.Interlocking {
// 		var uf cfg.Module
// 		var ut cfg.Module

// 		for _, v1 := range v.Units {
// 			if v1 == form {
// 				uf, _ = chenxi.CX.Cfg.Modules.GetCfgByUnit(form)
// 				continue
// 			}
// 			if v1 == to {
// 				uf, _ = chenxi.CX.Cfg.Modules.GetCfgByUnit(form)
// 			}
// 		}
// 		if uf.Name == form && ut.Name == to {
// 			isin := uf.PositionNo - ut.PositionNo

// 			if v.Lock.TryLock() {
// 				v.Rlock.Lock()
// 				defer v.Rlock.Unlock()
// 				v.Isin = isin

// 				return true
// 			} else {

// 				v.Rlock.RLock()
// 				defer v.Rlock.RUnlock()
// 				if isin == 0 {
// 					return false
// 				} else if isin > 0 {
// 					//  判断原值

// 					if v.Isin > 0 {
// 						return true
// 					} else {
// 						return false
// 					}
// 				} else {

// 					if v.Isin < 0 {
// 						return true
// 					} else {
// 						return false
// 					}
// 				}
// 			}

// 		}
// 	}

//		return true
//	}
//
// 自定义实现
func (u *CTC) InterlockingLock(id string, form, to string) bool {

	for _, v := range u.cfg.Interlocking {
		var uf cfg.Module
		var ut cfg.Module

		for _, v1 := range v.Units {
			if v1 == form || v1 == to {

				uf, _ = chenxi.CX.Cfg.Modules.GetCfgByUnit(form)

				ut, _ = chenxi.CX.Cfg.Modules.GetCfgByUnit(to)

				isin := uf.PositionNo - ut.PositionNo

				if v.Lock.TryLock() {
					v.Rlock.Lock()
					defer v.Rlock.Unlock()
					v.Isin = isin
					v.Pall[id] = &cfg.TAction{From: form, To: to, Isdone: false}
					return true
				} else {

					v.Rlock.RLock()
					defer v.Rlock.RUnlock()
					if p, ok := v.Pall[id]; ok {
						v.Isin = isin
						p.From = form
						p.To = to
						p.Isdone = false
						return true
					}

					if isin == 0 {
						return false
					} else if isin > 0 {
						//  判断原值

						if v.Isin > 0 {
							if len(v.Pall) < v.OutParallel {
								v.Pall[id] = &cfg.TAction{From: form, To: to, Isdone: false}
								return true
							}
							return false
						} else {
							return false
						}
					} else {

						if v.Isin < 0 {
							if len(v.Pall) < v.InParallel {
								v.Pall[id] = &cfg.TAction{From: form, To: to, Isdone: false}
								return true
							}
							return false

						} else {
							return false
						}
					}
				}

			}

		}

	}

	return true
}
func (u *CTC) InterlockingUnLockAll() {
	for _, v := range u.cfg.Interlocking {
		v.Lock.TryLock()

		v.Isin = 0
		v.Pall = make(map[string]*cfg.TAction)
		v.Lock.Unlock()
	}
}
func (u *CTC) InterlockingUnLock(id string, form, to string) bool {

	for _, v := range u.cfg.Interlocking {
		var uf cfg.Module
		var ut cfg.Module

		for _, v1 := range v.Units {
			if v1 == form {
				uf, _ = chenxi.CX.Cfg.Modules.GetCfgByUnit(form)
				continue
			}
			if v1 == to {
				ut, _ = chenxi.CX.Cfg.Modules.GetCfgByUnit(to)
			}
		}
		if uf.Name == form || ut.Name == to {
			for _, p := range v.Pall {
				if p.To == form && p.Isdone == true {
					p.From = form
					p.To = to
				}
			}

			if p, ok := v.Pall[id]; ok {
				p.From = form
				p.To = to
				p.Isdone = true
				if ut.Name != to {
					v.Rlock.Lock()
					defer v.Rlock.Unlock()

					//v.Isin = 0
					// for _, p := range v.Pall {
					// 	if p.From == form && p.To == to {
					// 		delete(v.Pall, id)
					// 	}
					// }
					var idsToDelete []string

					for id, p := range v.Pall {
						if p.From == form && p.To == to {
							idsToDelete = append(idsToDelete, id)
						}
					}

					for _, id := range idsToDelete {
						delete(v.Pall, id)
					}
					if len(v.Pall) == 0 {
						v.Isin = 0
						v.Lock.Unlock()
						return true
					}
				}
			}
		}
		// if uf.Name == form && ut.Name != to {
		// 	//出去
		// 	v.Rlock.Lock()
		// 	defer v.Rlock.Unlock()

		// 	v.Isin = 0

		// 	return true
		// }
		//isin := uf.PositionNo - ut.PositionNo

	}

	return true
}

//	func (c *CTC) GenActions(form ,to  string) (*MoveActions, error) {
//		return
//	}
func (c *CTC) GenActions(from, to string) (*MoveActions, error) {
	var parm = fmt.Sprintf("%s-%s", from, to)
	c.lock.Lock()
	defer c.lock.Unlock()
	// ma := &MoveAction{id: parm}
	if _, ok := c.Actions[parm]; !ok {

		// ft := strings.Split(parm, "-")
		// if len(ft) != 2 {
		// 	return nil, fmt.Errorf("格式错误 parm ：%s", parm)
		// }
		// from := ft[0]
		// to := ft[1]

		u2, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(from)
		if err != nil {
			log.Error(err.Error())
			u2 = cfg.Module{}
		}
		uto, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(to)
		if err != nil {
			log.Error(err.Error())
			uto = cfg.Module{}
		}
		tmActions := MoveActionCmds{From: from, To: to, Priority: u2.Priority}

		// group
		// 如何通过 group 查找相差几个单元的路径
		// 注意不同 TM 不同src 相同 dist 问题，或者同一个src 不同tm 不同dist ，都会出现问题（要以src 或者 dist 避免相同TM的访问
		for k, units := range c.cfg.Group {

			if containsBoth(units, from, to) {
				// tm, err := chenxi.CX.GetModule(k)
				// if err != nil {
				// 	log.Error("ctc 配置错误，不存在的tm", "tm", k)
				// 	continue
				// }
				// pikc, place, err := tm.(api.TMApi).GenAction(context.TODO(), parm)
				// enable disable
				u1, err := chenxi.CX.Cfg.Modules.GetCfgByUnit(k)
				if err != nil {
					//log.Error(err.Error())
					u1 = cfg.Module{}
				}
				// 同一个tm 两个arm 只能同时动一个
				if _, ok := chenxi.CX.Cfg.Modules.TMcfgs[k]; !ok {
					chenxi.CX.Cfg.Modules.TMcfgs[k] = cfg.TMcfg{Name: k, Arms: []string{"k"}}
				}

				for _, v := range chenxi.CX.Cfg.Modules.TMcfgs[k].Arms {

					instruct1 := fmt.Sprintf("%s-%s", from, v)
					instruct2 := fmt.Sprintf("%s-%s", v, to)

					p1 := 0
					p1 = c.cfg.Priority[instruct1]
					p2 := 0
					p2 = c.cfg.Priority[instruct2]
					if len(c.cfg.Disable_paths) > 0 {
						if containsAny(c.cfg.Disable_paths, instruct1, instruct2) {
							continue
						}

						tmActions.Actions = append(tmActions.Actions,
							TMAction{Name: k, Priority: u1.Priority, Arm: v, form: from, to: to, end: to,
								PickPriority: p1, PlacePriority: p2})

					} else {
						if containsBoth(c.cfg.Enable_paths, instruct1, instruct2) {
							tmActions.Actions = append(tmActions.Actions,
								TMAction{Name: k, Priority: u1.Priority, Arm: v, form: from, to: to, end: to,
									PickPriority: p1, PlacePriority: p2})

						}

					}
				}

			}
		}

		if len(tmActions.Actions) == 0 {
			// inpaths
			paths := c.cfg.In_paths
			if u2.PositionNo > uto.PositionNo {
				paths = c.cfg.Return_paths
			}

			pps := deeppath(paths, []string{from}, from, to)
			if pps != nil {
				// from to 路径
				for _, ps := range pps {
					mvas, err := c.GenActions(ps[0], ps[1])

					if err != nil {
						return nil, err
					}
					for _, ac := range mvas.Cmds.Actions {
						ac.end = to
					}
					tmActions.Actions = append(tmActions.Actions, mvas.Cmds.Actions...)
					// tmpa := mvas
					// for i := 1; i < len(ps)-1; i++ {
					// 	mvas1, err := c.GenActions(fmt.Sprintf("%s_%s", ps[i], ps[1+1]))
					// 	// if err != nil {
					// 	// 	return nil, err
					// 	// }
					// 	// for _, t := range tmpa.Cmds.Actions {
					// 	// 	t.Next = tmpa
					// 	// }
					// 	tmpa = mvas1
					// }

				}
			}
		}

		c.Actions[parm] = tmActions
		// 优先级排序
	}

	return &MoveActions{Cmds: c.Actions[parm], Lock: &sync.Mutex{}, Lockfrom: &sync.Mutex{}, Lockto: &sync.Mutex{}}, nil
}

func deeppath(mappath map[string]string, rpath []string, from, to string) [][]string {
	paths := [][]string{}
	if ps, ok := mappath[from]; ok {
		if strings.Contains(ps, to) {
			return append(paths, append(rpath, to))
		} else {
			p := strings.Split(ps, ",")
			for _, v := range p {
				subpaths := deeppath(mappath, append(rpath, v), v, to)
				if len(subpaths) >= 0 {
					paths = append(paths, subpaths...)
				}
			}
		}
	}
	return paths
}
func containsBoth(arr []string, from string, to string) bool {
	fromExists := false
	toExists := false
	for _, v := range arr {
		if v == from {
			fromExists = true
		}
		if v == to {
			toExists = true
		}
		if fromExists && toExists {
			return true
		}
	}
	return false
}
func containsAny(arr []string, elements ...string) bool {
	for _, v := range arr {
		for _, element := range elements {
			if v == element {
				return true
			}
		}
	}
	return false
}
