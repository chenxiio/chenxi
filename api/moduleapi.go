package api

import (
	"context"
)

// m_dictFunc["Ready"] = DoReady
// m_dictFunc["End"] = DoEnd
// m_dictFunc["PreRecv"] = DoPreRecv
// m_dictFunc["PostRecv"] = DoPostRecv
// m_dictFunc["PreSend"] = DoPreSend
// m_dictFunc["PostSend"] = DoPostSend
// m_dictFunc["Map"] = DoMap
// m_dictFunc["Load"] = DoLoad
// m_dictFunc["Unload"] = DoUnload
// m_dictFunc["Home"] = DoHome
// m_dictFunc["Init"] = DoInit

type ModuleApi interface {
	Init(ctx context.Context, parm string) error   //perm:none
	PreIn(ctx context.Context, parm string) error  //perm:none
	In(ctx context.Context, parm string) error     //perm:none
	PreOut(ctx context.Context, parm string) error //perm:none
	Out(ctx context.Context, parm string) error    //perm:none
	Move(ctx context.Context, parm string) error   //perm:none
	Next(ctx context.Context, parm string) error   //perm:none
}
