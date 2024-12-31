package modules

import (
	"context"
	"fmt"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/logger"
)

type CMTest struct {
	Name      string
	statename string
	//carid     string
	log *logger.Logger
}

func Map(dist, cid, wid string, dslot int) error {

	err := chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.%s", wid), fmt.Sprintf("%s.%d", dist, dslot))
	if err != nil {
		return err
	}
	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.wid.carrier.%s", wid), cid)
	if err != nil {
		return err
	}
	err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.wid.%d", dist, dslot), wid)
	if err != nil {
		return err
	}
	return nil
}
func (p *CMTest) Init(ctx context.Context, parm string) error {
	p.log = logger.GetLog(p.Name, "CM", chenxi.CX.Cfg.Basedir)

	p.log.Debug("Init called with parm:", parm)
	p.statename = p.Name + ".state"
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Init")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}

	// 检查无异常设置 设备状态 IDLE
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")

}
func (p *CMTest) PreIn(ctx context.Context, parm string) error {
	p.log.Debug("PreIn called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PreIn")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) In(ctx context.Context, parm string) error {
	p.log.Debug("In called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "In")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) PreOut(ctx context.Context, parm string) error {
	p.log.Debug("PreOut called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PreOut")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) Out(ctx context.Context, parm string) error {
	p.log.Debug("Out called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Out")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) Move(ctx context.Context, parm string) error {
	p.log.Debug("Move called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Move")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) Map(ctx context.Context, parm string) error {
	p.log.Debug("Map called with parm:", parm)

	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Map")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	cid, err := chenxi.CX.IOServer.ReadString(context.TODO(), fmt.Sprintf("%s.carrier.id", p.Name))
	if err != nil {
		return fmt.Errorf("%s.carrier.id , %s", p.Name, err.Error())
	}
	if chenxi.Simulation {
		//cid := "foup1"
		for i := 0; i < 2; i++ {
			Map(p.Name, cid, fmt.Sprintf("%s_%d", cid, i+1), i+1)
		}
	}

	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}

func (p *CMTest) Next(ctx context.Context, parm string) error {
	p.log.Debug("Next called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Next")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}

	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) Load(ctx context.Context, parm string) error {
	p.log.Debug("Load called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Load")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	if chenxi.Simulation {
		err := chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.carrier.id", p.Name), "foup_"+p.Name)
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.carrier.%s", "foup_"+p.Name), p.Name)
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteInt(context.TODO(), p.Name+".load", 1)
		if err != nil {
			return err
		}

		err = chenxi.CX.IOServer.WriteInt(context.TODO(), p.Name+".unload", 0)
		if err != nil {
			return err
		}
	}

	// p.carid = "carid"
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) Unload(ctx context.Context, parm string) error {
	p.log.Debug("Unload called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Unload")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	// 删除unit 上所有 env.wid.
	if chenxi.Simulation {

		err = chenxi.CX.IOServer.WriteInt(context.TODO(), p.Name+".load", 0)
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteInt(context.TODO(), p.Name+".unload", 1)
		if err != nil {
			return err
		}
		err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("env.carrier.%s", "foup_"+p.Name), "")
		if err != nil {
			return err
		}

		err = chenxi.CX.IOServer.WriteString(context.TODO(), fmt.Sprintf("%s.carrier.id", p.Name), "")
		if err != nil {
			return err
		}

	}

	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *CMTest) Home(ctx context.Context, parm string) error {
	p.log.Debug("Home called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Home")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
