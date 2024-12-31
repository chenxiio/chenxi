package modules

import (
	"context"
	"time"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/logger"
)

type PMTest struct {
	Name      string
	statename string
	log       *logger.Logger
}

func (p *PMTest) Init(ctx context.Context, parm string) error {
	p.log = logger.GetLog(p.Name, "PM", chenxi.CX.Cfg.Basedir)
	p.log.Debug("Init called with parm:", parm)
	p.statename = p.Name + ".state"
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Init")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}

	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *PMTest) PreIn(ctx context.Context, parm string) error {
	p.log.Debug("PreIn called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PreIn")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		// 注意 abort 立刻退出
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) In(ctx context.Context, parm string) error {
	p.log.Debug("In called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "In")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) PreOut(ctx context.Context, parm string) error {
	p.log.Debug("PreOut called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PreOut")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}

func (p *PMTest) Out(ctx context.Context, parm string) error {
	p.log.Debug("Out called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Out")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}

func (p *PMTest) Move(ctx context.Context, parm string) error {
	p.log.Debug("Move called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Move")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}

func (p *PMTest) Next(ctx context.Context, parm string) error {
	p.log.Debug("Next called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Next")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}
func (p *PMTest) Ready(ctx context.Context, parm string) error {
	p.log.Debug("Ready called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Ready")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		// 注意 abort 立刻退出
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) Process(ctx context.Context, parm string) error {
	p.log.Debug("Process called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Process")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		// 注意 abort 立刻退出
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) Abort(ctx context.Context, parm string) error {
	p.log.Debug("Abort called with parm:", parm)
	// //err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Abort")
	// if err != nil {
	// 	// 如果调用者判断调用失败，选择重试时先设置IDLE
	// 	return err
	// }
	go func() {
		time.Sleep(time.Millisecond * 30)
		//chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) Pause(ctx context.Context, parm string) error {
	p.log.Debug("Pause called with parm:", parm)
	// //err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Pause")
	// if err != nil {
	// 	// 如果调用者判断调用失败，选择重试时先设置IDLE
	// 	return err
	// }
	go func() {
		time.Sleep(time.Millisecond * 30)
		//chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) Resume(ctx context.Context, parm string) error {
	p.log.Debug("Resume called with parm:", parm)
	// err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Resume")
	// if err != nil {
	// 	// 如果调用者判断调用失败，选择重试时先设置IDLE
	// 	return err
	// }
	go func() {
		time.Sleep(time.Millisecond * 30)
		//chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}
func (p *PMTest) End(ctx context.Context, parm string) error {
	p.log.Debug("End called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "End")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}
	go func() {
		time.Sleep(time.Millisecond * 30)
		chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
	}()
	return nil
}

// Init(ctx context.Context, parm string) error    //perm:none
// 	PreIn(ctx context.Context, parm string) error   //perm:none
// 	In(ctx context.Context, parm string) error      //perm:none
// 	PreOut(ctx context.Context, parm string) error  //perm:none
// 	Out(ctx context.Context, parm string) error     //perm:none
// 	Ready(ctx context.Context, parm string) error   //perm:none
// 	Process(ctx context.Context, parm string) error //perm:none
// 	Abort(ctx context.Context, parm string) error   //perm:none
// 	Pause(ctx context.Context, parm string) error   //perm:none
// 	Resume(ctx context.Context, parm string) error  //perm:none
// 	End(ctx context.Context, parm string) error     //perm:none
