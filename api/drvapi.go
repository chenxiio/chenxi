package api

import (
	"context"
)

// Cmd(ctx context.Context, parm ...any) (any, error)
// ChangeEvnet(ctx context.Context,)
// Cmd(ctx context.Context, cmd string, parm ...any) ([]any, error)
// ChangeEvnet(ctx context.Context, fun func()) error
type Drvapi interface {
	Start(ctx context.Context, parm string) error
	Stop(ctx context.Context, parm string) error
	ReadInt(ctx context.Context, parm string) (int32, error)
	ReadString(ctx context.Context, parm string) (string, error)
	ReadDouble(ctx context.Context, parm string) (float64, error)
	WriteInt(ctx context.Context, parm string, value int32) error
	WriteString(ctx context.Context, parm string, value string) error
	WriteDouble(ctx context.Context, parm string, value float64) error
}
