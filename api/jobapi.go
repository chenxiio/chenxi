package api

import (
	"context"

	"github.com/chenxiio/chenxi/models"
)

// INIT	””或者单元名	CallCommand("CTC", "INIT", "", ref rspData);
// CallCommand("CTC", "INIT", "CSR", ref rspData);	””:CTC主要io变量里面做数据同步
// 单元名：CTC主要调用单元的Init指令
// SYS_EVENT	字符串:sender,event,data	CallCommand("CTC", "SYS_EVENT", "CSR,Placed,P01", ref rspData);	主要是给CTC发送系统事件，目前没有使用
// ABORT	””或者单元名	CallCommand("CTC", "ABORT", "", ref rspData);
// CallCommand("CTC", "ABORT", "CSR", ref rspData);	””:CTC Abort所有单元和CJ
// 单元名:CTC Abort单元
// CJOB_CREATE	obj:json string或者file:json file	CallCommand("CTC", "CJOB_CREATE", "obj:json string", ref rspData);
// CallCommand("CTC", "CJOB_CREATE", "file:json file name", ref rspData);	obj:json，结构见下页
// file:file name, LOT recipe目录下的jsonwenj，文件内容和obj传入的要一样
// CJOB_START	cjob id	CallCommand("CTC", "CJOB_START", "cjob id", ref rspData);	开始一个JOB，主要对于不是自动开始的CJOB
// CJOB_STOP	cjob id	CallCommand("CTC", "CJOB_STOP", "cjob id", ref rspData);	STOP CJOB, 对于没有从FOUP里面的Wafer，不再取片，已经出来的要做完工艺
// CJOB_PAUSE	cjob id	CallCommand("CTC", "CJOB_PAUSE", "cjob id", ref rspData);	暂停CJOB
// CJOB_RESUME	cjob id	CallCommand("CTC", "CJOB_RESUME", "cjob id", ref rspData);	对于暂停的CJOB，重新开始流片
// CJOB_ABORT	cjob id	CallCommand("CTC", "CJOB_ABORT", "cjob id", ref rspData);	ABORT CJOB, 整个JOB立即停止，对于已经在进行的动作，会等待做完。此时处理需要初始化有问题的单元，或者All Init，然后使用return wafer功能退片。
// GET_CJOB_LIST	””	CallCommand("CTC", "GET_CJOB_LIST", "", ref rspData);	返回所有CJ的数据，包括PJ，主要用于UI显示
// PJOB_CREATE	obj:json string或者file:json file	CallCommand("CTC", "PJOB_CREATE", "obj:json string", ref rspData);
// CallCommand("CTC", "PJOB_CREATE", "file:json file name", ref rspData);	创建PJ，结构见下页
// PJOB_SETUP	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现
// PJOB_START	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现
// PJOB_STOP	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现
// PJOB_CANCEL	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现
// PJOB_PAUSE	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现
// PJOB_RESUME	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现
// PJOB_ABORT	pjob id	CallCommand("CTC", "PJOB_SETUP", "pjob id", ref rspData);	暂未实现

type JobApi interface {
	Init(ctx context.Context, parm string) error                                     //perm:none
	Auto(ctx context.Context) error                                                  //perm:none
	Pause(ctx context.Context) error                                                 //perm:none
	Resume(ctx context.Context) error                                                //perm:none
	Abort(ctx context.Context) error                                                 //perm:none
	End(ctx context.Context) error                                                   //perm:none
	PJCreate(ctx context.Context, pj *models.ProcessJob) (*models.ProcessJob, error) //perm:none
	CJCreate(ctx context.Context, cj *models.ControlJob) (*models.ControlJob, error) //perm:none
	CJStart(ctx context.Context, cjid string) (*models.ControlJob, error)            //perm:none
	CJEnd(ctx context.Context, cjid string) (*models.ControlJob, error)              //perm:none
	CJPause(ctx context.Context, cjid string) (*models.ControlJob, error)            //perm:none
	CJResume(ctx context.Context, cjid string) (*models.ControlJob, error)           //perm:none
	CJAbort(ctx context.Context, cjid string) (*models.ControlJob, error)            //perm:none
	CJList(ctx context.Context) ([]*models.ControlJob, error)                        //perm:none
	ClearWafer(ctx context.Context, unit string) error                               //perm:none
	Continue(ctx context.Context) error                                              //perm:none
}
