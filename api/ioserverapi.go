package api

import (
	"context"
)

// ReadRange(ctx context.Context, prexx string) ([]IO, error)        //perm:read
//
//	type IO struct {
//		Key   string
//		Value any
//	}
type IOServerAPI interface {
	ReadInt(ctx context.Context, key string) (int32, error)                    //perm:none
	ReadString(ctx context.Context, key string) (string, error)                //perm:none
	ReadDouble(ctx context.Context, key string) (float64, error)               //perm:none
	WriteInt(ctx context.Context, key string, value int32) error               //perm:none
	WriteString(ctx context.Context, key string, value string) error           //perm:none
	WriteDouble(ctx context.Context, key string, value float64) error          //perm:none
	ReadFromPrefix(ctx context.Context, prefix string) (map[string]any, error) //perm:none
	SetState(ctx context.Context, key string, value string) error              //perm:none
}
