package cfg

import (
	"encoding/json"
	"os"
	"sync"
)

//	type Unit struct {
//		Name string
//		No   int // no相同直接锁定，from -to 大于0 流入 小于0 流出
//	}

type TAction struct {
	From   string
	To     string
	Isdone bool
}

type Interlock struct {
	Units       []string
	InParallel  int
	OutParallel int
	Lock        sync.Mutex          `json:"-"`
	Rlock       sync.RWMutex        `json:"-"`
	Isin        int                 `json:"-"`
	Pall        map[string]*TAction `json:"-"`
}
type CTCCfg struct {
	ProcessJobMode string              // ProcessRecipe:创建processjob时传入 processjob，参数模式
	Group          map[string][]string // move.1 后面跟两个单元，如果遇到这种类型，调用 src.preout  dist.prein move  src.postout  dist.postin
	Enable_paths   []string
	Disable_paths  []string          //主要是禁用手臂
	Return_paths   map[string]string //可以直接查找的跨单元子流程流出流程
	In_paths       map[string]string //可以直接查找的跨单元子流程流入流程 执行单元流程时注意锁定 Interlocking
	Priority       map[string]int    // 搬运指令优先级，数字越大优先级越高
	Interlocking   []*Interlock      // 不能同时搬运,流入锁和流出锁

	ReturningMode int // 做完recipe 工艺后返回模式，0 返回原位，1 调用return_paths, 2，固定单元列表，不分先后, 3 返回fuopid 位置 ,4,结束
	path          string
}

func (u *CTCCfg) SaveFile() error {
	data, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(u.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (u *CTCCfg) ReadFile() error {
	data, err := os.ReadFile(u.path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, u)
	if err != nil {
		return err
	}
	return nil
}
