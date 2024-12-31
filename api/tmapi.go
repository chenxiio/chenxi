package api

import (
	"context"
)

// m_dictFunc["End"] = DoEnd;
// m_dictFunc["Ready"] = DoReady;
// m_dictFunc["PreRecv"] = DoPreRecv;
// m_dictFunc["PostRecv"] = DoPostRecv;
// m_dictFunc["PreSend"] = DoPreSend;
// m_dictFunc["PostSend"] = DoPostSend;
// m_dictFunc["Pick"] = DoPick;
// m_dictFunc["Place"] = DoPlace;
// m_dictFunc["MoveToSend"] = DoMoveToSend;
// m_dictFunc["MoveToRecv"] = DoMoveToRecv;
// m_dictFunc["Map"] = DoMap;
// m_dictFunc["Transfer"] = DoTransfer;
// m_dictFunc["Exit"] = DoExit;
// m_dictFunc["srv_exit"] = OnServerExit;
// m_dictFunc["srv_active"] = OnServerActive;
//GenAction(ctx context.Context, parm string) (pick []string, place []string, err error) //perm:none

type TMApi interface {
	Init(ctx context.Context, parm string) error      //perm:none
	PrePick(ctx context.Context, parm string) error   //perm:none
	PrePlace(ctx context.Context, parm string) error  //perm:none
	Pick(ctx context.Context, parm string) error      //perm:none
	Place(ctx context.Context, parm string) error     //perm:none
	PostPick(ctx context.Context, parm string) error  //perm:none
	PostPlace(ctx context.Context, parm string) error //perm:none
}
