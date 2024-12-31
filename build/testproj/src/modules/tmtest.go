package modules

import (
	"context"
	"time"

	"github.com/chenxiio/chenxi"
	"github.com/chenxiio/chenxi/logger"
)

type TMTest struct {
	Name      string
	statename string
	log       *logger.Logger
}

func (p *TMTest) Init(ctx context.Context, parm string) error {
	p.log = logger.GetLog(p.Name, "TM", chenxi.CX.Cfg.Basedir)
	p.log.Debug("Init called with parm:", parm)
	p.statename = p.Name + ".state"
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Init")
	if err != nil {
		// 如果调用者判断调用失败，选择重试时先设置IDLE
		return err
	}

	return chenxi.CX.IOServer.SetState(context.Background(), p.statename, "IDLE")
}

func (p *TMTest) PrePick(ctx context.Context, parm string) error {
	p.log.Debug("PrePick called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PrePick")
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
func (p *TMTest) PrePlace(ctx context.Context, parm string) error {
	p.log.Debug("PrePlace called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PrePlace")
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
func (p *TMTest) Pick(ctx context.Context, parm string) error {
	p.log.Debug("Pick called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Pick")
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
func (p *TMTest) Place(ctx context.Context, parm string) error {
	p.log.Debug("Place called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "Place")
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
func (p *TMTest) PostPick(ctx context.Context, parm string) error {
	p.log.Debug("PostPick called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PostPick")
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
func (p *TMTest) PostPlace(ctx context.Context, parm string) error {
	p.log.Debug("PostPlace called with parm:", parm)
	err := chenxi.CX.IOServer.SetState(context.Background(), p.statename, "PostPlace")
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
