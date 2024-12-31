package api

import (
	"context"
)

// 初始化 判断是否有片，如果有片，设置down ，return wafer
type PMApi interface {
	Init(ctx context.Context, parm string) error    //perm:none
	PreIn(ctx context.Context, parm string) error   //perm:none
	In(ctx context.Context, parm string) error      //perm:none
	PreOut(ctx context.Context, parm string) error  //perm:none
	Out(ctx context.Context, parm string) error     //perm:none
	Move(ctx context.Context, parm string) error    //perm:none
	Next(ctx context.Context, parm string) error    //perm:none
	Ready(ctx context.Context, parm string) error   //perm:none
	Process(ctx context.Context, parm string) error //perm:none
	Abort(ctx context.Context, parm string) error   //perm:none
	Pause(ctx context.Context, parm string) error   //perm:none
	Resume(ctx context.Context, parm string) error  //perm:none
	End(ctx context.Context, parm string) error     //perm:none
}
