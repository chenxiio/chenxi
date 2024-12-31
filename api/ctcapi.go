package api

import "context"

type CTCApi interface {
	Init(ctx context.Context, parm string) error     //perm:none
	Move(ctx context.Context, from, to string) error //perm:none
}
